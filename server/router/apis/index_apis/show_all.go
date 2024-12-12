// -------------------------------------------------
// Package index_apis
// Author: hanzhi
// Date: 2024/12/11
// -------------------------------------------------

package index_apis

import (
	"gcnote/server/dto"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"net/http"
)

func ShowIndexes(ctx *gin.Context) {
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
	// 查看当前user_id 的所有表内信息，
}
