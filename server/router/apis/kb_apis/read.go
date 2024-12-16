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
)

func ReadKBFile(ctx *gin.Context) {
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

	// 然后读取文件就行了
	/*
		我是这样想的，
		markdown文本是给一个个的字符串（呃，看来我得修改一下图片资源的部分了。
		然后图片的话专门做一些图片服务器
	*/
}
