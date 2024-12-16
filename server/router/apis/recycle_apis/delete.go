// -------------------------------------------------
// Package kb_apis
// Author: hanzhi
// Date: 2024/12/15
// -------------------------------------------------

package recycle_apis

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

// DeleteRecycleFile
// *
func DeleteRecycleFile(ctx *gin.Context) {
	var req dto.RecycleRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		zap.S().Debugf("mark 1")
		zap.S().Debugf("params : %+v", req)
		ctx.JSON(http.StatusBadRequest, dto.Fail(dto.ParamsErrCode))
		return
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

	// 开始删除操作
	var recycleFile model.Recycle
	tx := config.DB.Model(&recycleFile).Delete(
		"kb_file_id", req.KBFileId)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			zap.S().Errorf(
				"User with user_id %v source_index_name %v not found",
				currentUserId, req.KBFileId)
			ctx.JSON(
				http.StatusNotFound,
				dto.FailWithMessage(dto.IndexNotExistErrCode, "try to rename the index but record note found."),
			)
			return
		} else {
			zap.S().Errorf(
				"Unexcept Error: User with user_id %v source_index_name %v , ",
				currentUserId, req.KBFileId)
			ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
			return
		}
	}
	// 删除文件
	path := config.PathCfg.RecycleBinPath
	err = wrench.RemoveContents(filepath.Join(path, req.IndexId, req.KBFileId))
	if err != nil {
		zap.S().Errorf("Failed to remove file, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	ctx.JSON(http.StatusOK, dto.Success())

}
