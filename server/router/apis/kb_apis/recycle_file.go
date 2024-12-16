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

// RecycleKBFile
// @Summary 删除知识库里的文件（fixme:新写一个吧，写错了，应该是收到回收站的，这个留给回收站了）
// @Description 删除知识库内部的文件(MYSQL + 文件系统已完成 + redis不管了)
// @ID			recycle-kb-file
// @Tags		index
// @Accept		json
// @Produce		json
// @Param		request		body		dto.KBFileRequest true "文档请求体"
// @Success 	200			{object} 	dto.BaseResponse	"成功响应，返回success"
// @Failure		400			{object}	dto.BaseResponse	"参数错误(code:40000)"
// @Failure		401			{object}	dto.BaseResponse	"Token错误(code:40101)"
// @Failure		404			{object}	dto.BaseResponse	"该知识库不存在(code:40201)"
// @Failure		500			{object}	dto.BaseResponse	"服务器内部错误(code:50000)"
// @Router 		/index/recycle_file [post]
func RecycleKBFile(ctx *gin.Context) {
	var req dto.KBFileRequest
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
	err = config.DB.Where("index_id = ? AND kb_file_name = ?", index.IndexId, req.KBFileName).First(&kbFile).Error
	if err != nil {
		zap.S().Errorf("Failed to find kb file, err: %v", err)
		ctx.JSON(http.StatusOK, dto.FailWithMessage(dto.KBFileNotExistErrCode, "该文件不存在"))
		return
	}

	// 回收文件
	src := filepath.Join(config.PathCfg.KnowledgeBasePath, kbFile.IndexId, kbFile.KBFileId)
	dst := filepath.Join(config.PathCfg.RecycleBinPath, kbFile.IndexId, kbFile.KBFileId)
	err = os.MkdirAll(dst, os.ModePerm)
	if err != nil {
		zap.S().Errorf("Failed to create recycle dir, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	// 复制文件到要求的路径
	err = wrench.CopyDir(src, dst)
	if err != nil {
		zap.S().Errorf("Failed to create recycle dir, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	// 删除原来的文件夹
	err = wrench.RemoveContents(src)
	if err != nil {
		zap.S().Errorf("Delete kb_file name :%v err:%v", req.KBFileName, err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}
	// 然后sql操作
	recycleModel := model.Recycle{
		SourceIndexId: kbFile.IndexId,
		KBFileId:      kbFile.KBFileId,
		KBFileName:    kbFile.KBFileName,
	}

	// 首先删除旧表
	tx := config.DB.Begin()
	tx.Model(&kbFile).Where(
		"kb_file_id = ?",
		kbFile.KBFileId,
	).Delete(&kbFile)
	if tx.Error != nil {
		zap.S().Errorf("Delete kb_file name :%v err:%v", req.KBFileName, tx.Error)
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}
	// 然后放到回收站里
	tx.Create(&recycleModel)
	if tx.Error != nil {
		zap.S().Errorf("recycle kb_file name :%v err:%v", req.KBFileName, tx.Error)
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}
	// 提交事务
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
