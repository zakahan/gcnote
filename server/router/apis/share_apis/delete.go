// -------------------------------------------------
// Package share_apis
// Author: hanzhi
// Date: 2025/1/4
// -------------------------------------------------

package share_apis

import (
	"gcnote/server/config"
	"gcnote/server/dto"
	"gcnote/server/model"
	"gcnote/server/router/wrench"
	"net/http"
	"os"
	"path/filepath"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// DeleteShareFile 删除分享文件
// @Summary 删除分享文件
// @Description 删除分享文件及其所有分享记录
// @Tags Share
// @Accept json
// @Produce json
// @Param request body dto.DeleteShareFileRequest true "删除分享文件请求"
// @Success 200 {object} dto.BaseResponse
// @Router /share/delete [post]
func DeleteShareFile(ctx *gin.Context) {
	// 1. 获取用户ID
	// 获取userId
	claims, exists := ctx.Get("claims")
	if !exists {
		zap.S().Infof("Unable to get the claims")
		ctx.JSON(http.StatusUnauthorized, dto.Fail(dto.UserTokenErrCode))
		zap.S().Debugf("mark 2")
		return
	}
	currentUser := claims.(jwt.MapClaims)
	userId := currentUser["sub"].(string)

	// 2. 解析请求
	var req dto.DeleteShareFileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, dto.Fail(dto.ParamsErrCode))
		return
	}

	// 3. 检查分享文件是否存在
	var shareFile model.ShareFile
	if err := config.DB.Where("share_file_id = ? AND user_id = ?", req.ShareFileId, userId).First(&shareFile).Error; err != nil {
		ctx.JSON(http.StatusOK, dto.Fail(dto.ShareFileNotExistErrCode))
		return
	}

	// 4. 删除分享文件记录
	if err := config.DB.Delete(&shareFile).Error; err != nil {
		zap.S().Errorf("删除分享文件记录错误：错误码 %v", err)
		ctx.JSON(http.StatusOK, dto.FailWithMessage(dto.InternalErrCode, "删除分享文件记录失败"))
		return
	}

	// 5. 删除文件系统中的文件
	fileDirPath := filepath.Join(config.PathCfg.ShareFileDirPath, shareFile.ShareFileId)
	if err := wrench.RemoveContents(fileDirPath); err != nil && !os.IsNotExist(err) {
		zap.S().Errorf("删除文件失败：错误码 %v", err)
		ctx.JSON(http.StatusOK, dto.FailWithMessage(dto.InternalErrCode, "删除文件失败"))
		return
	}

	ctx.JSON(http.StatusOK, dto.Success())
}
