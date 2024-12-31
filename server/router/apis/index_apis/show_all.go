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
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"net/http"
)

// ShowIndexes
// @Summary 展示同一个用户下的所有知识库
// @Description 展示知识库(这个是知识库的摘要展示，如果要看更细节的比如每个表有哪些文件，需要另外的接口)
// @ID			show-all-index
// @Tags		index
// @Accept		json
// @Produce		json
// @Success 	200			{object} 	dto.BaseResponse	"成功响应，返回success"
// @Failure		400			{object}	dto.BaseResponse	"参数错误(code:40000)"
// @Failure		401			{object}	dto.BaseResponse	"Token错误(code:40101)"
// @Failure		500			{object}	dto.BaseResponse	"服务器内部错误(code:50000)"
// @Router 		/index/show_indexes [get]
func ShowIndexes(ctx *gin.Context) {
	claims, exists := ctx.Get("claims")
	if !exists {
		zap.S().Infof("Unable to get the claims")
		ctx.JSON(http.StatusUnauthorized, dto.Fail(dto.UserTokenErrCode))
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

	// 先尝试从缓存获取
	indexList, err := cache.GetUserIndexList(ctx, currentUserId)
	if err == nil {
		// 缓存命中，直接返回
		ctx.JSON(http.StatusOK, dto.SuccessWithData(indexList))
		return
	}
	if !errors.Is(err, redis.Nil) {
		// 如果是其他错误，记录日志
		zap.S().Errorf("Failed to get index list from cache: %v", err)
	}

	// 缓存未命中或发生错误，从数据库查询
	indexModel := model.Index{}
	tx := config.DB.Model(&indexModel).Where(
		"user_id = ?",
		currentUserId,
	)
	indexList = make([]model.Index, 0)
	err = tx.Find(&indexList).Error
	if err != nil {
		zap.S().Errorf("failed to find index: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	// 更新缓存
	err = cache.SetUserIndexList(ctx, currentUserId, indexList)
	if err != nil {
		zap.S().Errorf("Failed to set index list cache: %v", err)
	}

	ctx.JSON(http.StatusOK, dto.SuccessWithData(indexList))
}
