// -------------------------------------------------
// Package user_apis
// Author: hanzhi
// Date: 2024/12/10
// -------------------------------------------------

package user_apis

//
//// Info
//// @Summary 用户信息
//// @Description 获取用户信息并展示
//// @ID           user-info
//// @Tags         user
//// @Accept       json
//// @Produce      json
//func Info(context *gin.Context) {
//	//@Param        request  none    dto.LoginRequest  true  "登录请求体"
//	claim, _ := context.Get("claim")
//	currentUser := claim.(jwt.MapClaims)
//	currentUserId := currentUser["sub"].(float64)
//	if currentUserId == 0 {
//		context.JSON(http.StatusOK, dto.Fail(dto.ParamsErrCode))
//		return
//	}
//}
