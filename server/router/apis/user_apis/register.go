// -------------------------------------------------
// Package user_apis
// Author: hanzhi
// Date: 2024/12/9
// -------------------------------------------------

package user_apis

import (
	"errors"
	"gcnote/server/config"
	"gcnote/server/dto"
	"gcnote/server/model"
	"gcnote/server/router/wrench"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
)

// Register
// @Summary      用户注册接口
// @Description  用户通过用户名和Email进行注册。注册成功后返回Success标记
// @ID           user-register
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request  body    dto.UserRegisterRequest  true  "注册请求体"
// @Success      200      {object} dto.BaseResponse "成功响应(code:0)"
// @Failure      400      {object} dto.BaseResponse "参数错误 (code: 40000)"
// @Failure      409      {object} dto.BaseResponse "用户已存在 (code: 40100)、电子邮箱已存在 (code: 40103)"
// @Failure      500      {object} dto.BaseResponse "内部服务器错误 (code: 50000)"
// @Router       /user/register [post]
func Register(ctx *gin.Context) {
	// 参数校验
	var req dto.UserRegisterRequest
	err := ctx.ShouldBindJSON(&req) // 检查
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Fail(dto.ParamsErrCode))
		return
	}
	// 逻辑处理
	userId := wrench.IdGenerator()
	userNew := model.User{
		UserId:   userId,
		UserName: req.UserName,
		Email:    req.Email,
		Password: wrench.HashPassword(req.Password),
	}

	//
	userSearch := model.User{}

	// 1. 检查是否存在相同用户名，如果存在则无法注册
	tx := config.DB.Where("user_name = ?", userNew.UserName).First(&userSearch)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		zap.S().Errorf("Register query user: %v err: %v", userNew.UserName, tx.Error)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}
	// tx.RowsAffected
	if tx.RowsAffected != 0 { // 行数不为0
		// 如果用户已经存在
		ctx.JSON(http.StatusConflict, dto.Fail(dto.UserExistsErrCode))
		return
	}

	// 2. 判断是否存在相同的电子邮箱，如果存在这返回失败
	tx = config.DB.Where("email = ?", userNew.Email).First(&userSearch)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		zap.S().Errorf("Register create user: %v err: %v", userNew, tx.Error)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}
	if tx.RowsAffected != 0 { // 行数不为0
		// 如果用户已经存在
		ctx.JSON(http.StatusConflict, dto.Fail(dto.UserEmailExistsErrCode))
		return
	}
	// 成功创建
	tx = config.DB.Create(&userNew)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		zap.S().Errorf("Register create user: %v err: %v", userNew, tx.Error)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}
	zap.S().Infof("Register create user %v done.", userNew.UserName)
	ctx.JSON(http.StatusOK, dto.Success())
}
