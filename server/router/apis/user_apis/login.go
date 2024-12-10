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
	"gcnote/server/router/middleware"
	"gcnote/server/router/wrench"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"time"
)

// Login
// @Summary      用户登录接口
// @Description  用户通过用户名和密码进行登录，成功后返回JWT令牌。
// @ID           user-login
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request  body    dto.LoginRequest  true  "登录请求体"
// @Success      200      {object} dto.BaseResponse "成功响应，返回JWT令牌(code:0)"
// @Failure      400      {object} dto.BaseResponse "参数错误(code:40000)"
// @Failure      401      {object} dto.BaseResponse "密码错误(code:40102)"
// @Failure      404      {object} dto.BaseResponse "记录不存在(code:40001)"
// @Failure		 500      {object} dto.BaseResponse "服务器内部错误(code:50000)"
// @Router       /user/login [post]
func Login(context *gin.Context) {
	// 参数校验
	var req dto.LoginRequest
	err := context.ShouldBindJSON(&req)
	if err != nil {
		context.JSON(http.StatusOK, dto.Fail(dto.ParamsErrCode))
		return
	}
	// 逻辑处理，查询是否在数据库，如果在就登录成功，否则失败
	userSearch := model.User{
		UserName: req.UserName,
	}
	tx := config.DB.Where("user_name =?", userSearch.UserName).First(&userSearch)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		zap.S().Errorf("Login query user:%+v err:%v", userSearch, tx.Error)
		context.JSON(http.StatusOK, dto.Fail(dto.InternalErrCode))
		return
	}
	// 然后是查看用户是否存在
	if userSearch.ID == 0 {
		context.JSON(http.StatusOK, dto.FailWithMessage(dto.RecordNotFoundErrCode, "用户账号不存在"))
		return
	}

	// 密码校验
	if wrench.CheckPassword(userSearch.Password, req.Password) != nil {
		context.JSON(http.StatusOK, dto.Fail(dto.UserPasswordErrCode))
		return
	}
	// 这代表成功了，然后提示登录成功
	// 初始化一个jwt
	newJwt := middleware.NewJWT()
	claims := jwt.MapClaims{
		"sub":  userSearch.UserId,                     // 用户ID
		"name": userSearch.UserName,                   // 用户名
		"exp":  time.Now().Add(time.Hour * 24).Unix(), // 过期时间为一天
	}
	// 生成JWT
	token, err := newJwt.GenerateJWT(claims)
	if err != nil {
		zap.S().Infof("[CreateToken] 生成token失败")
		context.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成token失败",
		})
		return
	}
	context.JSON(http.StatusOK, dto.SuccessWithData(token))
	return
}
