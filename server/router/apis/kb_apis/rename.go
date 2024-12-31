// -------------------------------------------------
// Package kb_apis
// Author: hanzhi
// Date: 2024/12/16
// -------------------------------------------------

package kb_apis

import (
	"errors"
	"gcnote/server/cache"
	"gcnote/server/config"
	"gcnote/server/dto"
	"gcnote/server/model"
	"gcnote/server/router/wrench"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"os"
	"path/filepath"
)

// RenameKBFile
// [rewrite]需要重写
// @Summary 重命名文件
// @Description 重命名文件(
// @ID			rename-kb-file
// @Tags		index
// @Accept		json
// @Produce		json
// @Param		request		body		dto.KBFileRenameRequest true "登录请求体"
// @Success 	200			{object} 	dto.BaseResponse	"成功响应，或文件不存在"
// @Failure 	400			{object} 	dto.BaseResponse	"参数错误"
// @Failure 	401			{object} 	dto.BaseResponse	"用户未登录"
// @Failure 	500			{object} 	dto.BaseResponse	"服务器错误"
// @Router 		/index/rename_file [post]
func RenameKBFile(ctx *gin.Context) {
	var req dto.KBFileRenameRequest
	zap.S().Debugf("Rename request params : %v", req)

	// 参数绑定
	if err := ctx.ShouldBindJSON(&req); err != nil {
		zap.S().Debugf("Invalid parameters: %+v", req)
		ctx.JSON(http.StatusBadRequest, dto.Fail(dto.ParamsErrCode))
		return
	}

	// 验证新文件名是否合法
	if !wrench.ValidateKBName(req.DestKBFileName) {
		zap.S().Debugf("Invalid KB file name: %v", req.DestKBFileName)
		ctx.JSON(http.StatusBadRequest, dto.FailWithMessage(dto.KBFileNameErrCode, "dest KB file name not allowed."))
		return
	}

	// 获取用户信息
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
		ctx.JSON(http.StatusUnauthorized, dto.Fail(dto.ParamsErrCode))
		return
	}

	// 先从缓存验证index是否存在
	index, err := cache.GetIndexInfo(ctx, req.IndexId)
	if errors.Is(err, redis.Nil) {
		// 缓存未命中，从数据库查询
		index, err = cache.RefreshIndexInfo(ctx, req.IndexId)
		if err != nil {
			zap.S().Errorf("Failed to get index info: %v", err)
			ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
			return
		}
	} else if err != nil {
		zap.S().Errorf("Failed to get index from cache: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	if index == nil || index.IndexId == "" || index.UserId != currentUserId {
		ctx.JSON(http.StatusNotFound, dto.Fail(dto.IndexNotExistErrCode))
		return
	}

	// 从缓存验证kb文件是否存在
	kbFile, err := cache.GetKBInfo(ctx, req.KBFileId)
	if errors.Is(err, redis.Nil) {
		// 缓存未命中，从数据库查询
		kbFile, err = cache.RefreshKBInfo(ctx, req.KBFileId)
		if err != nil {
			zap.S().Errorf("Failed to get kb file info: %v", err)
			ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
			return
		}
	} else if err != nil {
		zap.S().Errorf("Failed to get kb file from cache: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	if kbFile == nil || kbFile.KBFileId == "" || kbFile.IndexId != req.IndexId {
		ctx.JSON(http.StatusNotFound, dto.Fail(dto.KBFileNotExistErrCode))
		return
	}

	// 检查原来的系统中是否存在这个名字 不允许重命名
	var existKBFile model.KBFile
	tx := config.DB.Model(&existKBFile).Where("kb_file_name = ? AND index_id = ? ",
		req.DestKBFileName, req.IndexId).First(&existKBFile)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		zap.S().Debugf("kb_file_name: %s is exist", req.DestKBFileName)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}
	if tx.RowsAffected != 0 { // 行数不为0 ，说明已经存在了
		ctx.JSON(http.StatusConflict, dto.Fail(dto.KBFileExistErrCode))
		return
	}

	// 文件名修改
	// os.Rename(oldPath, newPath)
	documentPath := filepath.Join(config.PathCfg.KnowledgeBasePath, req.IndexId, req.KBFileId, req.KBFileName+".md")
	renameDocumentPath := filepath.Join(config.PathCfg.KnowledgeBasePath, req.IndexId, req.KBFileId, req.DestKBFileName+".md")
	zap.S().Debugf("修改文件，原始路径 %v, 新路径 %v", documentPath, renameDocumentPath)
	err = os.Rename(documentPath, renameDocumentPath)
	if err != nil {
		zap.S().Errorf("Failed to rename file, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.FailWithMessage(dto.InternalErrCode, "文件系统重命名文件失败"))
		return
	}

	// 修改SQL
	err = config.DB.Model(&model.KBFile{}).Where("kb_file_id = ?", req.KBFileId).Update("kb_file_name", req.DestKBFileName).Error
	if err != nil {
		zap.S().Errorf("Failed to update kb file, err: %v", err)
		ctx.JSON(http.StatusOK, dto.FailWithMessage(dto.InternalErrCode, "更新文件名失败"))
		return
	}

	// 更新缓存
	// 1. 刷新kb文件信息缓存
	_, err = cache.RefreshKBInfo(ctx, req.KBFileId)
	if err != nil {
		zap.S().Errorf("Failed to refresh kb file cache: %v", err)
	}
	// 2. 刷新index的kb文件列表缓存
	_, err = cache.RefreshIndexKBList(ctx, req.IndexId)
	if err != nil {
		zap.S().Errorf("Failed to refresh index kb list cache: %v", err)
	}
	// 3. 刷新用户的最近访问kb列表缓存
	_, err = cache.RefreshRecentKBList(ctx, currentUserId)
	if err != nil {
		zap.S().Errorf("Failed to refresh user recent kb list cache: %v", err)
	}

	// 成功
	ctx.JSON(http.StatusOK, dto.Success())
}
