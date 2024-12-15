// -------------------------------------------------
// Package kb_apis
// Author: hanzhi
// Date: 2024/12/13
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

/*
1. 获取request
2. 获取userId
3. 验证index存在（所以说redis还是得搞，格式就是user_id/index_name/把）  redis最后补上
4. 在kb_file表新建表项（事务
5. 创建对应的文件夹（kb_file_id/image, kb_file_id/file.md)

*/

func CreateKBFile(ctx *gin.Context) {
	var req dto.KBFileRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		zap.S().Debugf("mark 1")
		zap.S().Debugf("params : %+v", req)
		ctx.JSON(http.StatusOK, dto.Fail(dto.ParamsErrCode))
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
		ctx.JSON(http.StatusInternalServerError, dto.FailWithMessage(dto.IndexNotExistErrCode,
			"知识库"+req.IndexName+"不存在"))
	}
	// 开始事务
	tx = config.DB.Begin()
	if tx.Error != nil {
		zap.S().Errorf("Failed to begin transaction, err:%v", tx.Error)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}
	// 创建一个KBFile
	KBId := wrench.IdGenerator()
	KBFileNew := model.KBFile{
		IndexId:    indexSearch.IndexId,
		KBFileId:   KBId,
		KBFileName: req.KBFileName,
	}
	err = tx.Create(&KBFileNew).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		zap.S().Errorf("Create index %v, Error: %v", KBFileNew.KBFileName, tx.Error)
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}
	// 文件系统
	kbPath := config.PathCfg.KnowledgeBasePath
	kbFilePath := filepath.Join(kbPath, indexSearch.IndexId, KBId)
	err = os.Mkdir(kbFilePath, os.ModePerm)
	if err != nil {
		zap.S().Errorf("Create KBFile Dir %s Error: %v\n", kbFilePath, err)
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, dto.FailWithMessage(dto.InternalErrCode, "create file dir error"))
		return
	}
	// 生成文件夹 image
	imageDirPath := filepath.Join(kbFilePath, "image")
	err = os.Mkdir(imageDirPath, os.ModePerm)
	if err != nil {
		zap.S().Errorf("Create Image Dir %s Error: %v\n", imageDirPath, err)
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, dto.FailWithMessage(dto.InternalErrCode, "create file dir error"))
		return
	}
	// 创建md文件
	mdPath := filepath.Join(kbFilePath, req.KBFileName+".md")
	file, err := os.Create(mdPath)
	if err != nil {
		zap.S().Errorf("Create File %s Error: %v\n", mdPath, err)
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, dto.FailWithMessage(dto.InternalErrCode, "create file dir error"))
		return
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			zap.S().Errorf("Close File Error %s Error: %v\n", mdPath, err)
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, dto.FailWithMessage(dto.InternalErrCode, "create file dir error"))
			return
		}
	}(file)

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		zap.S().Errorf("Failed to commit transaction, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	zap.S().Infof("Create file name %v done.", KBFileNew.KBFileName)
	ctx.JSON(http.StatusOK, dto.Success())

}
