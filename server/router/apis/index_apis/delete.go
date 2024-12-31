// -------------------------------------------------
// Package index_apis
// Author: hanzhi
// Date: 2024/12/11
// -------------------------------------------------

package index_apis

import (
	"gcnote/server/ability/search_engine"
	"gcnote/server/cache"
	"gcnote/server/config"
	"gcnote/server/dto"
	"gcnote/server/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"net/http"
	"os"
	"path/filepath"
)

// DeleteIndex
// @Summary      删除知识库及其关联文件
// @Description  检查存在+删除知识库(MYSQL + 文件系统已完成 + 更新redis)，同时删除关联的KBFile记录及文件
// @ID           delete-index
// @Tags         index
// @Accept       json
// @Produce      json
// @Param        request  body      dto.IndexRequest true "索引请求体"
// @Success      200      {object}  dto.BaseResponse "成功响应，返回success"
// @Failure      400      {object}  dto.BaseResponse "参数错误(code:40000)"
// @Failure      401      {object}  dto.BaseResponse "Token错误(code:40101)"
// @Failure      404      {object}  dto.BaseResponse "该知识库不存在(code:40201)"
// @Failure      500      {object}  dto.BaseResponse "服务器内部错误(code:50000)"
// @Router       /index/delete_index [post]

func DeleteIndex(ctx *gin.Context) {
	var req dto.IndexRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		zap.S().Debugf("Invalid parameters: %+v", req)
		ctx.JSON(http.StatusBadRequest, dto.Fail(dto.ParamsErrCode))
		return
	}

	// 验证用户身份
	claims, exists := ctx.Get("claims")
	if !exists {
		zap.S().Infof("Unable to get the claims")
		ctx.JSON(http.StatusUnauthorized, dto.Fail(dto.UserTokenErrCode))
		return
	}
	currentUser := claims.(jwt.MapClaims)
	currentUserId := currentUser["sub"].(string)
	if currentUserId == "" {
		zap.S().Debugf("currentUserId is empty")
		ctx.JSON(http.StatusBadRequest, dto.Fail(dto.ParamsErrCode))
		return
	}

	// 开始事务
	tx := config.DB.Begin()

	// 检查知识库是否存在
	var indexInfo model.Index
	if err := tx.Where("index_id = ?", req.IndexId).First(&indexInfo).Error; err != nil {
		zap.S().Errorf("Knowledge base not found: %v", req.IndexId)
		tx.Rollback()
		ctx.JSON(http.StatusNotFound, dto.Fail(dto.IndexNotExistErrCode))
		return
	}

	// 删除关联的KBFile记录
	var kbFiles []model.KBFile
	if err := tx.Where("index_id = ?", req.IndexId).Find(&kbFiles).Error; err != nil {
		zap.S().Errorf("Failed to get KBFile records for index_id %v: %v", req.IndexId, err)
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	// 删除关联的KBFile记录
	if err := tx.Where("index_id = ?", req.IndexId).Delete(&model.KBFile{}).Error; err != nil {
		zap.S().Errorf("Failed to delete KBFile records for index_id %v: %v", req.IndexId, err)
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	// 删除知识库记录
	if err := tx.Delete(&indexInfo).Error; err != nil {
		zap.S().Errorf("Failed to delete index %v: %v", req.IndexId, err)
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	// 删除文件系统中的知识库及关联文件
	kbPath := config.PathCfg.KnowledgeBasePath
	xyPath := filepath.Join(kbPath, req.IndexId)
	if _, err := os.Stat(xyPath); os.IsNotExist(err) {
		zap.S().Infof("Path %v does not exist, skipping deletion", xyPath)
	} else {
		if err := os.RemoveAll(xyPath); err != nil {
			zap.S().Errorf("Failed to delete path %v: %v", xyPath, err)
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
			return
		}
	}

	// -----------------------------------------------
	// 查看是否index在es中是否存在
	err, code := search_engine.IndexExist(config.ElasticClient, "gcnote-"+req.IndexId)
	if err != nil {
		zap.S().Errorf("Failed to check exist of index, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}
	if code != 200 {
		zap.S().Errorf("the index is not exist, err: %v")
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}
	// 删除index
	err = search_engine.IndexDelete(config.ElasticClient, "gcnote-"+req.IndexId)
	if err != nil {
		zap.S().Errorf("Failed to delete the index, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		zap.S().Errorf("Failed to commit transaction: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	// 清理缓存
	// 1. 删除index信息缓存
	err = cache.DelIndexInfo(ctx, req.IndexId)
	if err != nil {
		zap.S().Errorf("Failed to delete index cache: %v", err)
	}
	// 2. 删除index的kb文件列表缓存
	err = cache.DelIndexKBList(ctx, req.IndexId)
	if err != nil {
		zap.S().Errorf("Failed to delete index kb list cache: %v", err)
	}
	// 3. 删除每个kb文件的缓存
	for _, kbFile := range kbFiles {
		err = cache.DelKBInfo(ctx, kbFile.KBFileId)
		if err != nil {
			zap.S().Errorf("Failed to delete kb file cache: %v", err)
		}
	}
	// 4. 刷新用户的index列表缓存
	_, err = cache.RefreshUserIndexList(ctx, currentUserId)
	if err != nil {
		zap.S().Errorf("Failed to refresh user index list cache: %v", err)
	}
	// 5. 刷新用户的最近访问kb列表缓存
	_, err = cache.RefreshRecentKBList(ctx, currentUserId)
	if err != nil {
		zap.S().Errorf("Failed to refresh user recent kb list cache: %v", err)
	}

	zap.S().Infof("Deleted index %v and its associated files for user %v", req.IndexId, currentUserId)
	ctx.JSON(http.StatusOK, dto.Success())
}
