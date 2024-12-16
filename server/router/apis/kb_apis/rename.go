// -------------------------------------------------
// Package kb_apis
// Author: hanzhi
// Date: 2024/12/16
// -------------------------------------------------

package kb_apis

import (
	"errors"
	"gcnote/server/config"
	"gcnote/server/dto"
	"gcnote/server/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"os"
	"path/filepath"
)

// RenameKBFile
// [rewrite]需要重写
// @Summary 重命名文件
// @Description 重命名文件(
// @ID			rename-kb-file
// @Tags		index
// @Accept		json
// @Produce		json
// @Param		request		body		dto.KBFileRenameRequest true "登录请求体"
// @Success 	200			{object} 	dto.BaseResponse	"成功响应，或文件不存在"
// @Failure 	400			{object} 	dto.BaseResponse	"参数错误"
// @Failure 	401			{object} 	dto.BaseResponse	"用户未登录"
// @Failure 	500			{object} 	dto.BaseResponse	"服务器错误"
// @Router 		/index/rename_file [post]
func RenameKBFile(ctx *gin.Context) {
	var req dto.KBFileRenameRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Fail(dto.ParamsErrCode))
		return
	}
	// 获取用户信息
	claims, exists := ctx.Get("claims")
	if !exists {
		zap.S().Infof("Unable to get the claims")
		ctx.JSON(http.StatusUnauthorized, dto.Fail(dto.UserTokenErrCode))
		return
	}
	currentUser := claims.(jwt.MapClaims)
	currentUserId := currentUser["sub"].(string)
	if currentUserId == "" {
		zap.S().Debugf("currentUserId is empty")
		ctx.JSON(http.StatusUnauthorized, dto.Fail(dto.ParamsErrCode))
		return
	}

	// 然后查找index表
	// 查找 Index 表中的 IndexId
	var index model.Index
	if err = config.DB.Where("user_id = ? AND index_name = ?", currentUserId, req.IndexName).First(&index).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Index not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	var kbFile model.KBFile
	err = config.DB.Where("index_id = ? AND kb_file_name = ?", index.IndexId, req.SourceFileName).First(&kbFile).Error
	if err != nil {
		zap.S().Errorf("Failed to find kb file, err: %v", err)
		ctx.JSON(http.StatusOK, dto.FailWithMessage(dto.KBFileNotExistErrCode, "该文件不存在"))
		return
	}

	// 文件名修改
	// os.Rename(oldPath, newPath)
	documentPath := filepath.Join(config.PathCfg.KnowledgeBasePath, kbFile.IndexId, kbFile.KBFileId, kbFile.KBFileName+".md")
	renameDocumentPath := filepath.Join(config.PathCfg.KnowledgeBasePath, kbFile.IndexId, kbFile.KBFileId, req.DestFileName+".md")
	zap.S().Debugf("修改文件，原始路径 %v, 新路径 %v", documentPath, renameDocumentPath)
	err = os.Rename(documentPath, renameDocumentPath)
	if err != nil {
		zap.S().Errorf("Failed to rename file, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.FailWithMessage(dto.InternalErrCode, "文件系统重命名文件失败"))
		return
	}
	// 修改SQL
	err = config.DB.Model(&kbFile).Update("kb_file_name", req.DestFileName).Error
	if err != nil {
		zap.S().Errorf("Failed to update kb file, err: %v", err)
		ctx.JSON(http.StatusOK, dto.FailWithMessage(dto.InternalErrCode, "更新文件名失败"))
		return
	}
	// 成功
	ctx.JSON(http.StatusOK, dto.Success())

}
