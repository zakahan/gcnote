// -------------------------------------------------
// Package kb_apis
// Author: hanzhi
// Date: 2024/12/15
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
	"path/filepath"
)

// DeleteKBFile
// *
func DeleteKBFile(ctx *gin.Context) {
	var req dto.KBFileRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		zap.S().Debugf("mark 1")
		zap.S().Debugf("params : %+v", req)
		ctx.JSON(http.StatusBadRequest, dto.Fail(dto.ParamsErrCode))
		return
	}

	// 验证KBName的有效性
	if !wrench.ValidateKBName(req.KBFileName) {
		zap.S().Debugf("Unaccept index name %v", req.KBFileName)
		ctx.JSON(http.StatusBadRequest, dto.Fail(dto.KBFileNameErrCode))
	}

	claims, exists := ctx.Get("claims")
	if !exists {
		zap.S().Infof("Unable to get the claims")
		ctx.JSON(http.StatusUnauthorized, dto.Fail(dto.UserTokenErrCode))
		zap.S().Debugf("mark 2")
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

	// 查看知识库是否存在
	// fixme 这里要改为redis
	indexSearch := model.Index{}
	tx := config.DB.Where("index_name = ? AND user_id = ?",
		req.IndexName, currentUserId,
	).First(&indexSearch)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		zap.S().Errorf("Create index %v, err: %v", req.IndexName, tx.Error)
		ctx.JSON(http.StatusOK, dto.Fail(dto.InternalErrCode))
		return
	}
	if tx.RowsAffected == 0 {
		ctx.JSON(http.StatusConflict, dto.FailWithMessage(dto.IndexNotExistErrCode,
			"知识库"+req.IndexName+"不存在"))
	}

	if tx.Error != nil {
		zap.S().Errorf("Failed to begin transaction, err:%v", tx.Error)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}
	// 创建一个KBFile
	kbFileModel := model.KBFile{
		IndexId: indexSearch.IndexId,
		//KBFileId:   空的,
		KBFileName: req.KBFileName,
	}
	// 检查是否存在同名的文件（kbFileId）
	var kbFile = model.KBFile{}
	tx = config.DB.Model(&kbFile).Where("kb_file_name = ? AND index_id = ? ",
		kbFileModel.KBFileName, kbFileModel.IndexId).First(&kbFile) // fixme: 这里的代码有问题，需要解决
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		zap.S().Debugf("kb_file_name: %s is exist", kbFileModel.KBFileName)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
		// 继续往下走
	}
	if tx.RowsAffected == 0 { // 行数为0 ，说明不存在了，说明请求有错误
		ctx.JSON(http.StatusBadRequest, dto.Fail(dto.KBFileNotExistErrCode))
		return
	}

	// 开始事务
	tx = config.DB.Begin()
	// SQL删除
	tx = tx.Model(&kbFile).Where(
		"kb_file_id = ?",
		kbFile.KBFileId,
	).Delete(&kbFile)
	if tx.Error != nil {
		zap.S().Errorf("Delete kb_file name :%v err:%v", req.KBFileName, tx.Error)
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}
	// 从文件系统删除
	kbPath := config.PathCfg.KnowledgeBasePath
	indexPath := filepath.Join(kbPath, kbFile.KBFileId)
	err = wrench.RemoveContents(indexPath)
	if err != nil {
		zap.S().Errorf("Delete kb_file name :%v err:%v", req.KBFileName, tx.Error)
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	err = tx.Commit().Error
	if err != nil {
		zap.S().Errorf("Failed to commit transaction, err: %v", err)
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}
	zap.S().Infof("Delete kb_file name %v , user id : %v done.", req.IndexName, currentUserId)
	ctx.JSON(http.StatusOK, dto.Success())
}
