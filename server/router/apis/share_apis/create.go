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

// CreateShareFile 创建分享文件
// @Summary 创建分享文件
// @Description 基于已存在的知识库文件创建分享文件
// @Tags Share
// @Accept json
// @Produce json
// @Param request body dto.CreateShareFileRequest true "创建分享文件请求"
// @Success 200 {object} dto.BaseResponse
// @Router /share/create [post]
func CreateShareFile(ctx *gin.Context) {
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
	var req dto.CreateShareFileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, dto.Fail(dto.ParamsErrCode))
		return
	}

	// 3. 检查原始知识库文件是否存在
	var kbFile model.KBFile
	if err := config.DB.Where("kb_file_id = ? AND user_id = ?", req.KBFileId, userId).First(&kbFile).Error; err != nil {
		ctx.JSON(http.StatusOK, dto.Fail(dto.KBFileNotExistErrCode))
		return
	}

	// 4. 创建分享文件记录
	shareFilePassword := wrench.RandStringBytes(16)
	shareFile := model.ShareFile{
		ShareFileId: kbFile.KBFileId, // 使用原文件的ID
		IndexId:     kbFile.IndexId,  // 记录来源知识库
		KBFileId:    kbFile.KBFileId, // 记录来源文件
		FileName:    kbFile.KBFileName,
		UserId:      userId,
		Password:    shareFilePassword,
	}

	// 5. 创建分享文件目录
	sharePath := filepath.Join(config.PathCfg.ShareFileDirPath, kbFile.KBFileId)
	if err := os.MkdirAll(filepath.Dir(sharePath), 0755); err != nil {
		ctx.JSON(http.StatusOK, dto.FailWithMessage(dto.InternalErrCode, "创建分享文件目录失败"))
		return
	}

	// 6. 复制文件内容
	kbPath := config.PathCfg.KnowledgeBasePath
	srcPath := filepath.Join(kbPath, kbFile.IndexId, kbFile.KBFileId)
	err := wrench.CopyDir(srcPath, sharePath)
	if err != nil {
		ctx.JSON(http.StatusOK, dto.FailWithMessage(dto.InternalErrCode, "复制文件目录失败"))
		return
	}

	// 7. 保存到数据库
	if err = config.DB.Create(&shareFile).Error; err != nil {
		zap.S().Errorf("保存到数据库ShareFile表出错：%v", err)
		ctx.JSON(http.StatusOK, dto.FailWithMessage(dto.InternalErrCode, "保存分享文件记录失败"))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessWithData(map[string]string{
		"password": shareFilePassword,
	}))
}
