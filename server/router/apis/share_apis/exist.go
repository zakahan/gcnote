// -------------------------------------------------
// Package share_apis
// Author: hanzhi
// Date: 2024/1/4
// -------------------------------------------------

package share_apis

import (
	"gcnote/server/config"
	"gcnote/server/dto"
	"gcnote/server/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

// CheckShareFileExist 检查分享文件是否存在
// @Summary 检查分享文件是否存在
// @Description 检查指定ID的分享文件是否存在
// @Tags Share
// @Accept json
// @Produce json
// @Param shareFileId path string true "分享文件ID"
// @Success 200 {object} dto.BaseResponse{data=dto.ShareFileExistResponse}
// @Router /share/exist [get]
func CheckShareFileExist(ctx *gin.Context) {
	// 1. 获取用户ID
	claims, exists := ctx.Get("claims")
	if !exists {
		zap.S().Infof("Unable to get the claims")
		ctx.JSON(http.StatusUnauthorized, dto.Fail(dto.UserTokenErrCode))
		return
	}
	currentUser := claims.(jwt.MapClaims)
	userId := currentUser["sub"].(string)

	// 2. 获取分享文件ID
	shareFileId := ctx.Query("share_file_id")
	if shareFileId == "" {
		ctx.JSON(http.StatusOK, dto.Fail(dto.ParamsErrCode))
		return
	}

	// 3. 查询数据库
	var shareFile model.ShareFile
	result := config.DB.Where("share_file_id = ? AND user_id = ?", shareFileId, userId).First(&shareFile)

	response := dto.ShareFileExistResponse{
		ShareFileId: shareFileId,
		Exist:       result.Error == nil,
	}

	ctx.JSON(http.StatusOK, dto.SuccessWithData(response))
}
