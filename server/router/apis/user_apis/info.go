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
func Info(ctx *gin.Context) {
	//@Param        request  none    dto.LoginRequest  true  "登录请求体"
	claim, _ := ctx.Get("claim")
	currentUser := claim.(jwt.MapClaims)
	currentUserId := currentUser["sub"].(string)
	if currentUserId == "" {

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
		ctx.JSON(http.StatusOK, dto.SuccessWithData(user))
		return
	}
}
