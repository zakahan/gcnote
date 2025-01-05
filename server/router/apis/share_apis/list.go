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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

// ListShareFiles 获取用户的所有分享文件
// @Summary 获取用户的所有分享文件
// @Description 获取当前用户的所有分享文件列表
// @Tags Share
// @Accept json
// @Produce json
// @Success 200 {object} dto.BaseResponse{data=dto.ShareFileListResponse}
// @Router /share/info [get]
func ListShareFiles(ctx *gin.Context) {
	// 1. 获取用户ID
	claims, exists := ctx.Get("claims")
	if !exists {
		zap.S().Infof("Unable to get the claims")
		ctx.JSON(http.StatusUnauthorized, dto.Fail(dto.UserTokenErrCode))
		return
	}
	currentUser := claims.(jwt.MapClaims)
	userId := currentUser["sub"].(string)

	// 2. 查询数据库
	var shareFiles []model.ShareFile
	var total int64

	if err := config.DB.Model(&model.ShareFile{}).Where("user_id = ?", userId).Count(&total).Error; err != nil {
		zap.S().Errorf("Failed to count share files: %v", err)
		ctx.JSON(http.StatusOK, dto.Fail(dto.InternalErrCode))
		return
	}

	if err := config.DB.Where("user_id = ?", userId).Find(&shareFiles).Error; err != nil {
		zap.S().Errorf("Failed to query share files: %v", err)
		ctx.JSON(http.StatusOK, dto.Fail(dto.InternalErrCode))
		return
	}

	// 3. 构建响应
	shareFileInfos := make([]dto.ShareFileInfo, 0, len(shareFiles))
	for _, file := range shareFiles {
		shareFileInfos = append(shareFileInfos, dto.ShareFileInfo{
			ShareFileId: file.ShareFileId,
			IndexId:     file.IndexId,
			KBFileId:    file.KBFileId,
			FileName:    file.FileName,
			Password:    file.Password,
			CreatedAt:   file.CreatedAt.Format(time.RFC3339),
		})
	}

	response := dto.ShareFileListResponse{
		Total: total,
		List:  shareFileInfos,
	}

	ctx.JSON(http.StatusOK, dto.SuccessWithData(response))
}
