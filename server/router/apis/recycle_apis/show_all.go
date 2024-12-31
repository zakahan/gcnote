// -------------------------------------------------
// Package recycle_apis
// Author: hanzhi
// Date: 2024/12/16
// -------------------------------------------------

package recycle_apis

import (
	"errors"
	"gcnote/server/cache"
	"gcnote/server/dto"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"net/http"
)

// ShowRecycleFiles
// @Summary      获取回收站中的所有文件
// @Description  展示回收站中的所有文件信息
// @ID           show-recycle-files
// @Tags         recycle
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.BaseResponse  						"成功响应，返回文件列表"
// @Failure      400  {object}  dto.BaseResponse                        "参数错误(code:40000)"
// @Failure      401  {object}  dto.BaseResponse                        "Token错误(code:40101)"
// @Failure      500  {object}  dto.BaseResponse                        "服务器内部错误(code:50000)"
// @Router       /recycle/show_files [get]
func ShowRecycleFiles(ctx *gin.Context) {
	// 展示所有文件
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

	// 从缓存获取用户的回收站文件列表
	recycleList, err := cache.GetUserRecycleList(ctx, currentUserId)
	if errors.Is(err, redis.Nil) {
		// 缓存未命中，从数据库查询
		recycleList, err = cache.RefreshUserRecycleList(ctx, currentUserId)
		if err != nil {
			zap.S().Errorf("Failed to get user recycle list: %v", err)
			ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
			return
		}
	} else if err != nil {
		zap.S().Errorf("Failed to get user recycle list from cache: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessWithData(recycleList))
}
