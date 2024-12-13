// -------------------------------------------------
// Package user_apis
// Author: hanzhi
// Date: 2024/12/10
// -------------------------------------------------

package user_apis

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

// UpdateUserName
// @Summary 更新用户名
// @Description 获取用户名，并修改之
// @ID 			user-update-user-name
// @Tags        user
// @Accept		json
// @Produce 	json
// @Param   	request  body    dto.UpdateUserNameRequest  true  "更新用户名请求体"
// @success 	200		{object} dto.BaseResponse "成功响应，返回更新成功的消息"
// @Failure 	200		{object} dto.BaseResponse "参数错误(code:40000)"
// @Failure 	200		{object} dto.BaseResponse "用户验证错误(code:40101)"
// @Failure     200		{object} dto.BaseResponse "服务器内部错误(code:50000)"
// @Router	/user/update_user_name [post]
func UpdateUserName(ctx *gin.Context) {
	var req dto.UpdateUserNameRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, dto.FailWithMessage(dto.ParamsErrCode, "get user_name:"+req.UserName))
		return
	}
	zap.S().Debugf("req.UserName get as %v\n", req.UserName)
	// jwt get
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
	// 逻辑处理部分
	userInfo := model.User{}
	//tx := config.DB.Where("user_id = ?", currentUserId).First(&userInfo)
	tx := config.DB.Model(&userInfo).Where(
		"user_id = ?", currentUserId,
	).Update("user_name", req.UserName)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			zap.S().Errorf("User with user_id %v not found", currentUserId)
			ctx.JSON(
				http.StatusOK,
				dto.FailWithMessage(dto.RecordNotFoundErrCode, "try to update but record note found."),
			)
		} else {
			zap.S().Errorf("Failed to get user, id: %v", currentUserId)
			ctx.JSON(http.StatusOK, dto.Fail(dto.InternalErrCode))
		}
	}

	// 刷新缓存
	_, err = cache.RefreshUserInfo(ctx.Request.Context(), currentUserId)
	if err != nil {
		zap.S().Errorf("Update.refreshUserInfoCache, UserID: %v err: %v", currentUserId, tx.Error)
		ctx.JSON(http.StatusOK, dto.Fail(dto.ParamsErrCode))
		return
	}

	ctx.JSON(http.StatusOK, dto.Success())

}

// UpdatePassword
// @Summary 更新用户密码
// @Description 获取用户id，并修改密码
// @ID 		user-update-password
// @Tags   user
// @Accept	json
// @Produce json
// @Param   request  body    dto.UpdatePasswordRequest  true  "更新用户名请求体"
// @success 200		{object} dto.BaseResponse "成功响应，返回更新成功的消息"
// @Failure 200		{object} dto.BaseResponse "参数错误(code:40000)"
// @Failure 200		{object} dto.BaseResponse "用户验证错误(code:40101)"
// @Failure 200		{object} dto.BaseResponse "服务器内部错误(code:50000)"
// @Router	/user/update_password [post]
func UpdatePassword(ctx *gin.Context) {
	// 这个怎么改，才能让jwt失效？
	var req dto.UpdatePasswordRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, dto.FailWithMessage(dto.ParamsErrCode, "get password:"+req.Password))
		return
	}
	zap.S().Debugf("req.Password get as %v\n", req.Password)
	// jwt get
	claims, exists := ctx.Get("claims")
	if !exists {
		zap.S().Infof("Unable to get the claims")
		ctx.JSON(http.StatusOK, dto.Fail(dto.UserTokenErrCode))
		return
	}
	zap.S().Debugf("claims: %v", claims)
	currentUser := claims.(jwt.MapClaims)
	currentUserId := currentUser["sub"].(string)
	// 逻辑处理部分
	userInfo := model.User{}
	//tx := config.DB.Where("user_id = ?", currentUserId).First(&userInfo)
	tx := config.DB.Model(&userInfo).Where("user_id = ?", currentUserId).Update("password", wrench.HashPassword(req.Password))

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			zap.S().Errorf("User with user_id %v not found", currentUserId)
			ctx.JSON(
				http.StatusOK,
				dto.FailWithMessage(dto.RecordNotFoundErrCode, "try to update but record note found."),
			)
		} else {
			zap.S().Errorf("Failed to get user, id: %v", currentUserId)
			ctx.JSON(http.StatusOK, dto.Fail(dto.InternalErrCode))
		}
	}

	// 刷新缓存
	_, err = cache.RefreshUserInfo(ctx.Request.Context(), currentUserId)
	if err != nil {
		zap.S().Errorf("Update.refreshUserInfoCache, UserID: %v err: %v", currentUserId, tx.Error)
		ctx.JSON(http.StatusOK, dto.Fail(dto.ParamsErrCode))
		return
	}

	ctx.JSON(http.StatusOK, dto.Success())

}
