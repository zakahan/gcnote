// -------------------------------------------------
// Package index_apis
// Author: hanzhi
// Date: 2024/12/11
// -------------------------------------------------

package index_apis

import (
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
// @Summary 删除知识库
// @Description 删除知识库(MYSQL + 文件系统已完成 + 更新redis)
// @ID			delete-index
// @Tags		index
// @Accept		json
// @Produce		json
// @Param		request		body		dto.IndexCRUDRequest true "登录请求体"
// @Success 	200			{object} 	dto.BaseResponse	"成功响应，返回success"
// @Failure		400			{object}	dto.BaseResponse	"参数错误(code:40000)"
// @Failure		401			{object}	dto.BaseResponse	"Token错误(code:40101)"
// @Failure		404			{object}	dto.BaseResponse	"该知识库不存在(code:40201)"
// @Failure		500			{object}	dto.BaseResponse	"服务器内部错误(code:50000)"
// @Router 		/index/delete_index [post]
func DeleteIndex(ctx *gin.Context) {
	var req dto.IndexCRUDRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		zap.S().Debugf("mark 1")
		zap.S().Debugf("params : %+v", req)
		ctx.JSON(http.StatusBadRequest, dto.Fail(dto.ParamsErrCode))
		return
	}
	claims, exists := ctx.Get("claims")
	if !exists {
		zap.S().Infof("Unable to get the claims")
		ctx.JSON(http.StatusUnauthorized, dto.Fail(dto.UserTokenErrCode))
		zap.S().Debugf("mark 2")
		return
	}
	zap.S().Debugf("claims: %v", claims)
	currentUser := claims.(jwt.MapClaims)
	currentUserId := currentUser["sub"].(string)
	if currentUserId == "" {
		zap.S().Debugf("currentUserId is empty")
		ctx.JSON(http.StatusBadRequest, dto.Fail(dto.ParamsErrCode))
		return
	}
	// 查看是否有这个index，然后删除
	indexModel := model.Index{}
	tx := config.DB.Model(&indexModel).Where(
		"index_name = ? AND user_id = ?",
		req.IndexName, currentUserId,
	)
	// 不需要再验证是否只有一条了，因为已经再create的时候验证过了
	err = tx.First(&indexModel).Error
	if err != nil {
		zap.S().Errorf("failed to find index: %v", err)
		ctx.JSON(http.StatusNotFound, dto.Fail(dto.IndexNotExistErrCode))
		return
	}
	indexId := indexModel.IndexId
	zap.S().Debugf("flag of index , index_id is %v", indexId)

	indexInfo := model.Index{}

	// 开始事务
	tx = config.DB.Begin()
	// 删除
	err = tx.Model(&indexInfo).Where(
		"index_id = ?",
		indexId,
	).Delete(&indexInfo).Error
	if err != nil {
		zap.S().Errorf("Delete index id :%v err:%v", indexId, tx.Error)
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	// 文件系统创建
	kbPath := config.PathCfg.KnowledgeBasePath
	xyPath := filepath.Join(kbPath, indexId)
	// 检查路径是否存在
	_, err = os.Stat(xyPath)
	if os.IsNotExist(err) {
		zap.S().Infof("路径 %v 不存在，无法删除操作", xyPath)
	}
	// 删除xyPath以及里面的所有文件
	err = os.RemoveAll(xyPath)
	if err != nil {
		zap.S().Errorf("删除路径 %v 失败， err: %v", xyPath, err)
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	err = tx.Commit().Error
	if err != nil {
		zap.S().Errorf("Failed to commit transaction, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}
	zap.S().Infof("Delete index name %v , user id : %v done.", req.IndexName, currentUserId)
	ctx.JSON(http.StatusOK, dto.Success())
}
