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
	"gcnote/server/router/wrench"
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
	zap.S().Debugf("Rename request params : %v", req)

	// 参数绑定
	if err := ctx.ShouldBindJSON(&req); err != nil {
		zap.S().Debugf("Invalid parameters: %+v", req)
		ctx.JSON(http.StatusBadRequest, dto.Fail(dto.ParamsErrCode))
		return
	}

	// 验证新文件名是否合法
	if !wrench.ValidateKBName(req.DestKBFileName) {
		zap.S().Debugf("Invalid KB file name: %v", req.DestKBFileName)
		ctx.JSON(http.StatusBadRequest, dto.FailWithMessage(dto.KBFileNameErrCode, "dest KB file name not allowed."))
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
	kbFile := model.KBFile{
		KBFileName: req.KBFileName,
		KBFileId:   req.KBFileId,
		IndexId:    req.IndexId,
	}

	// 检查原来的系统中是否存在这个名字 不允许重命名
	var existKBFile model.KBFile
	tx := config.DB.Model(&existKBFile).Where("kb_file_name = ? AND index_id = ? ",
		kbFile.KBFileName, kbFile.IndexId).First(&existKBFile)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		zap.S().Debugf("kb_file_name: %s is exist", kbFile.KBFileName)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	// 文件名修改
	// os.Rename(oldPath, newPath)
	documentPath := filepath.Join(config.PathCfg.KnowledgeBasePath, kbFile.IndexId, kbFile.KBFileId, kbFile.KBFileName+".md")
	renameDocumentPath := filepath.Join(config.PathCfg.KnowledgeBasePath, kbFile.IndexId, kbFile.KBFileId, req.DestKBFileName+".md")
	zap.S().Debugf("修改文件，原始路径 %v, 新路径 %v", documentPath, renameDocumentPath)
	err := os.Rename(documentPath, renameDocumentPath)
	if err != nil {
		zap.S().Errorf("Failed to rename file, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.FailWithMessage(dto.InternalErrCode, "文件系统重命名文件失败"))
		return
	}

	// 修改SQL
	err = config.DB.Model(&kbFile).Where("kb_file_id = ?", req.KBFileId).Update("kb_file_name", req.DestKBFileName).Error
	if err != nil {
		zap.S().Errorf("Failed to update kb file, err: %v", err)
		ctx.JSON(http.StatusOK, dto.FailWithMessage(dto.InternalErrCode, "更新文件名失败"))
		return
	}
	// 成功
	ctx.JSON(http.StatusOK, dto.Success())

}
