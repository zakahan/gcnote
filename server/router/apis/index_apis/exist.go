// -------------------------------------------------
// Package index_apis
// Author: hanzhi
// Date: 2024/12/22
// -------------------------------------------------

package index_apis

import (
	"errors"
	"gcnote/server/cache"
	"gcnote/server/config"
	"gcnote/server/dto"
	"gcnote/server/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
)

// IndexExist
// @Summary 检查知识库是否存在
// @Description 检查知识库是否存在
// @ID check-index
// @Tags index
// @Accept json
// @Produce json
// @Param indexId query string true "知识库ID"
// @Success 200 {object} dto.BaseResponse "成功响应，返回success"
// @Failure 400 {object} dto.BaseResponse "参数错误(code:40000)"
// @Failure 404 {object} dto.BaseResponse "知识库不存在(code:40201)"
// @Failure 500 {object} dto.BaseResponse "服务器内部错误(code:50000)"
// @Router /index/exist [get]
func IndexExist(ctx *gin.Context) {
	var req dto.IndexRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		zap.S().Debugf("mark 1")
		zap.S().Debugf("params : %+v", req)
		ctx.JSON(http.StatusBadRequest, dto.Fail(dto.ParamsErrCode))
		return
	}
	claims, _ := ctx.Get("claims")

	zap.S().Debugf("claims: %v", claims)
	currentUser := claims.(jwt.MapClaims)
	currentUserId := currentUser["sub"].(string)
	if currentUserId == "" {
		zap.S().Debugf("currentUserId is empty")
		ctx.JSON(http.StatusBadRequest, dto.Fail(dto.ParamsErrCode))
		return
	}

	// 先尝试从缓存获取
	_, err = cache.GetIndexInfo(ctx, req.IndexId)
	if err == nil {
		// 缓存命中，直接返回
		ctx.JSON(http.StatusOK, dto.Success())
		return
	}
	if !errors.Is(err, redis.Nil) {
		// 如果是其他错误，记录日志
		zap.S().Errorf("Failed to get index from cache: %v", err)
	}

	// 缓存未命中或发生错误，从数据库查询
	var dbIndex model.Index
	if err := config.DB.Where("index_id = ?", req.IndexId).First(&dbIndex).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, dto.Fail(dto.IndexNotExistErrCode))
		} else {
			ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		}
		return
	}

	// 更新缓存
	err = cache.SetIndexInfo(ctx, dbIndex)
	if err != nil {
		zap.S().Errorf("Failed to set index cache: %v", err)
	}

	ctx.JSON(http.StatusOK, dto.Success())
}
