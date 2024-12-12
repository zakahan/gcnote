// -------------------------------------------------
// Package user_apis
// Author: hanzhi
// Date: 2024/12/10
// -------------------------------------------------

package user_apis

import (
	"gcnote/server/cache"
	"gcnote/server/config"
	"gcnote/server/dto"
	"gcnote/server/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"net/http"
)

// Delete
// @Summary 删除用户
// @Description 获取用户id，并删除用户
// @ID 		user-delete
// @Tags    user
// @Accept	json
// @Produce json
// @success 200		{object} dto.BaseResponse "成功响应，返回更新成功的消息"
// @Failure 200		{object} dto.BaseResponse "参数错误(code:40000)"
// @Failure 200		{object} dto.BaseResponse "用户验证错误(code:40101)"
// @Failure 200		{object} dto.BaseResponse "服务器内部错误(code:50000)"
// @Router	/user/delete [post]
func Delete(ctx *gin.Context) {
	// 这个接口实际上应该是不会开放给用户的，所以不会出现 用户没了，但是jwt还在所以能用的情况
	// 根据 jwt，取出当前用户信息
	claims, _ := ctx.Get("claims") // 是否需要remove这个变量？
	currentUser := claims.(jwt.MapClaims)
	currentUserId := currentUser["sub"].(string)
	if currentUserId == "" {
		ctx.JSON(http.StatusOK, dto.FailWithMessage(dto.ParamsErrCode, "token is empty"))
		return
	}

	// 逻辑处理
	userInfo := model.User{}
	tx := config.DB.Model(&userInfo).Where("user_id = ?", currentUserId).Delete(&userInfo)
	if tx.Error != nil {
		zap.S().Errorf("Delete  userId:%v err:%v", currentUserId, tx.Error)
		ctx.JSON(http.StatusOK, dto.Fail(dto.InternalErrCode))
		return
	}

	// 更新redis
	err := cache.DelUserInfo(ctx.Request.Context(), currentUserId)
	if err != nil {
		zap.S().Errorf("Delete.DelUserInfo  userId:%v err:%v", currentUserId, tx.Error)
		ctx.JSON(http.StatusOK, dto.Fail(dto.InternalErrCode))
		return
	}
	ctx.JSON(http.StatusOK, dto.Success())
}
