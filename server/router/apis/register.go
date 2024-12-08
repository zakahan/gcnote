// -------------------------------------------------
// Package apis
// Author: hanzhi
// Date: 2024/12/8
// -------------------------------------------------

package apis

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

func Register(context *gin.Context) {
	// 参数校验
	var req dto.RegisterRequest
	err := context.ShouldBindJSON(&req) // 检查
	if err != nil {
		context.JSON(http.StatusOK, wrench.Fail(wrench.ParamsErrCode))
		return
	}
	// 逻辑处理
	userNew := model.User{
		UserName: req.UserName,
		Password: wrench.HashPassword(req.Password),
		Email:    req.Email,
	}

	//
	userSearch := model.User{}

	// 1. 检查是否存在相同用户名，如果存在则无法注册
	tx := config.DB.Where("user_name = ?", userNew.UserName).First(&userSearch)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {

		zap.S().Errorf("Register query user: %v err: %v", userNew.UserName, tx.Error)
		context.JSON(http.StatusOK, wrench.Fail(wrench.InternalErrCode))
		return
	}
	// tx.RowsAffected
	if tx.RowsAffected != 0 { // 行数不为0
		// 如果用户已经存在
		context.JSON(http.StatusOK, wrench.Fail(wrench.UserExistsErrCode))
		return
	}

	// 2. 判断是否存在相同的电子邮箱，如果存在这返回失败
	tx = config.DB.Where("email = ?", userNew.Email).First(&userSearch)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {

		zap.S().Errorf("Register create user: %v err: %v", userNew, tx.Error)
		context.JSON(http.StatusOK, wrench.Fail(wrench.InternalErrCode))
		return
	}
	if tx.RowsAffected != 0 { // 行数不为0
		// 如果用户已经存在
		context.JSON(http.StatusOK, wrench.Fail(wrench.UserEmailExistsErrCode))
		return
	}
	// 成功创建
	tx = config.DB.Create(&userNew)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {

		zap.S().Errorf("Register create user: %v err: %v", userNew, tx.Error)
		context.JSON(http.StatusOK, wrench.Fail(wrench.InternalErrCode))
		return
	}
	context.JSON(http.StatusOK, wrench.Success())
}

func Login(context *gin.Context) {
	// 参数校验
	var req dto.LoginRequest
	err := context.ShouldBindJSON(&req)
	if err != nil {
		context.JSON(http.StatusOK, wrench.Fail(wrench.ParamsErrCode))
		return
	}
	// 逻辑处理，查询是否在数据库，如果在就登录成功，否则失败
	userSearch := model.User{
		UserName: req.UserName,
	}
	tx := config.DB.Where("user_name =?", userSearch.UserName).First(&userSearch)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		zap.S().Errorf("Login query user:%+v err:%v", userSearch, tx.Error)
		context.JSON(http.StatusOK, wrench.Fail(wrench.InternalErrCode))
		return
	}
	// 然后是查看用户是否存在
	if userSearch.ID == 0 {
		context.JSON(http.StatusOK, wrench.Fail(wrench.RecordNotFoundErrCode))
		return
	}

	// 密码校验
	if wrench.CheckPassword(userSearch.Password, req.Password) != nil {
		context.JSON(http.StatusOK, wrench.Fail(wrench.UserPasswordErrCode))
		return
	}
	// 这代表成功了，然后提示登录成功
	// .... 没写完
}
