// -------------------------------------------------
// Package index_apis
// Author: hanzhi
// Date: 2024/12/11
// -------------------------------------------------

package index_apis

import (
	"errors"
	"gcnote/server/cache"
	"gcnote/server/config"
	"gcnote/server/dto"
	"gcnote/server/model"
	"gcnote/server/router/wrench"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
)

// RenameIndex
// @Summary      知识库重命名
// @Description  检查存在+重命名知识库(检查目标名称是否已存在，防止重名)
// @ID           rename-index
// @Tags         index
// @Accept       json
// @Produce      json
// @Param        request  body      dto.IndexRenameRequest true "重命名请求体"
// @Success      200      {object}  dto.BaseResponse       "成功响应，返回success"
// @Failure      400      {object}  dto.BaseResponse       "参数错误(code:40000)、新索引名称不符合规范(code:40202)"
// @Failure      401      {object}  dto.BaseResponse       "Token错误(code:40101)"
// @Failure      404      {object}  dto.BaseResponse       "原索引记录不存在(code:40201)"
// @Failure      409      {object}  dto.BaseResponse       "目标索引名称已存在(code:40203)"
// @Failure      500      {object}  dto.BaseResponse       "服务器内部错误(code:50000)"
// @Router       /index/rename_index [post]

func RenameIndex(ctx *gin.Context) {
	var req dto.IndexRenameRequest
	zap.S().Debugf("req params : %v", req)

	// 参数绑定
	if err := ctx.ShouldBindJSON(&req); err != nil {
		zap.S().Debugf("Invalid parameters: %+v", req)
		ctx.JSON(http.StatusBadRequest, dto.Fail(dto.ParamsErrCode))
		return
	}

	// 检查目标名称是否合法
	if !wrench.ValidateIndexName(req.DestIndexName) {
		zap.S().Debugf("Invalid index name: %v", req.DestIndexName)
		ctx.JSON(http.StatusBadRequest, dto.FailWithMessage(dto.IndexNameErrCode, "dest index name not allow."))
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

	// 检查是否存在同名的目标知识库
	var count int64
	if err := config.DB.Model(&model.Index{}).
		Where("index_name = ? AND user_id = ?", req.DestIndexName, currentUserId).
		Count(&count).Error; err != nil {
		zap.S().Errorf("Failed to query existing index name: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	if count > 0 {
		zap.S().Infof("Target index name already exists: %v", req.DestIndexName)
		ctx.JSON(http.StatusConflict, dto.FailWithMessage(dto.IndexExistErrCode, "Target index name already exists."))
		return
	}

	// 更新索引名称
	tx := config.DB.Model(&model.Index{}).
		Where("index_id = ? AND user_id = ?", req.IndexId, currentUserId).
		Update("index_name", req.DestIndexName)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			zap.S().Errorf("Index not found for user_id %v, index_id %v", currentUserId, req.IndexId)
			ctx.JSON(http.StatusNotFound, dto.FailWithMessage(dto.IndexNotExistErrCode, "Index not found."))
			return
		}
		zap.S().Errorf("Unexpected error while renaming index: %v", tx.Error)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	// 更新缓存
	// 1. 刷新index信息缓存
	_, err := cache.RefreshIndexInfo(ctx, req.IndexId)
	if err != nil {
		zap.S().Errorf("Failed to refresh index cache: %v", err)
	}
	// 2. 刷新用户的index列表缓存
	_, err = cache.RefreshUserIndexList(ctx, currentUserId)
	if err != nil {
		zap.S().Errorf("Failed to refresh user index list cache: %v", err)
	}

	zap.S().Infof("Successfully renamed index %v to %v for user %v", req.IndexId, req.DestIndexName, currentUserId)
	ctx.JSON(http.StatusOK, dto.Success())
}
