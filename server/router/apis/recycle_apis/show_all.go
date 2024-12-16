// -------------------------------------------------
// Package recycle_apis
// Author: hanzhi
// Date: 2024/12/16
// -------------------------------------------------

package recycle_apis

import (
	"gcnote/server/config"
	"gcnote/server/dto"
	"gcnote/server/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"net/http"
)

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

	var recycleTable []model.Recycle
	err := config.DB.Find(&recycleTable).Error
	if err != nil {
		zap.S().Errorf("failed to find recycle table: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}
	ctx.JSON(http.StatusOK, dto.SuccessWithData(recycleTable))
}
