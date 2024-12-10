// -------------------------------------------------
// Package user_apis
// Author: hanzhi
// Date: 2024/12/10
// -------------------------------------------------

package user_apis

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

// Info
// @Summary 用户信息
// @Description 获取用户信息并展示
// @ID           user-info
// @Tags         user
// @Accept       json
// @Produce      json
// @Success		 200	{object} dto.BaseResponse "成功响应，data返回用户信息"
// @Failure		 200	{object} dto.BaseResponse "参数错误(code:40000)"
// @Failure		 200	{object} dto.BaseResponse "用户验证错误(code:40101)"
// @Failure      200	{object} dto.BaseResponse "服务器内部错误(code:50000)"
// @Router		 /user/info [get]
func Info(ctx *gin.Context) {
	claims, exists := ctx.Get("claims")
	if !exists {
		zap.S().Infof("Unable to get the claims")
		ctx.JSON(http.StatusOK, dto.Fail(dto.UserTokenErrCode))
		return
	}
	zap.S().Debugf("claims: %v", claims)
	currentUser := claims.(jwt.MapClaims)
	currentUserId := currentUser["sub"].(string)
	if currentUserId == "" {
		zap.S().Debugf("currentUserId is empty")
		ctx.JSON(http.StatusOK, dto.Fail(dto.ParamsErrCode))
		return
	}
	// 查看缓存
	user, err := cache.GetUserInfo(ctx.Request.Context(), currentUserId)
	if !errors.Is(err, redis.Nil) && err != nil {
		zap.S().Errorf("Info.cache.GetUserInfo  userId:%+v err:%v", currentUserId, err)
		ctx.JSON(http.StatusOK, dto.Fail(dto.InternalErrCode))
		return
	}
	if user != nil {
		zap.S().Debugf("get user info by redis")
		ctx.JSON(http.StatusOK, dto.SuccessWithData(user))
		return
	}
	// 刷新一下redis
	zap.S().Debugf("refresh user info in redis.")
	user, err = cache.RefreshUserInfo(ctx.Request.Context(), currentUserId)
	if err != nil {
		zap.S().Errorf("Info.refreshUserInfoCache  user:%+v err:%v", currentUserId, err)
		ctx.JSON(http.StatusOK, dto.Fail(dto.InternalErrCode))
	}

	ctx.JSON(http.StatusOK, dto.SuccessWithData(user))
}
