// -------------------------------------------------
// Package recycle_apis
// Author: hanzhi
// Date: 2024/12/17
// -------------------------------------------------

package recycle_apis

import (
	"errors"
	"gcnote/server/cache"
	"gcnote/server/config"
	"gcnote/server/dto"
	"gcnote/server/model"
	"gcnote/server/router/wrench"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"net/http"
	"path/filepath"
)

// ClearUserRecycleBin
// @Summary      清空指定用户的回收站
// @Description  删除指定用户回收站中的所有文件，包括数据库记录和文件系统文件
// @ID           clear-user-recycle-bin
// @Tags         recycle
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.BaseResponse  "成功响应，返回success"
// @Failure      400  {object}  dto.BaseResponse  "参数错误(code:40000)"
// @Failure      401  {object}  dto.BaseResponse  "Token错误(code:40101)"
// @Failure      500  {object}  dto.BaseResponse  "服务器内部错误(code:50000)"
// @Router       /recycle/clear [get]
func ClearUserRecycleBin(ctx *gin.Context) {
	// 验证用户身份
	claims, exists := ctx.Get("claims")
	if !exists {
		zap.S().Infof("Unable to get the claims")
		ctx.JSON(http.StatusUnauthorized, dto.Fail(dto.UserTokenErrCode))
		return
	}

	// 获取当前用户ID
	currentUser := claims.(jwt.MapClaims)
	currentUserId, ok := currentUser["sub"].(string)
	if !ok || currentUserId == "" {
		zap.S().Errorf("Invalid user ID in claims")
		ctx.JSON(http.StatusBadRequest, dto.Fail(dto.ParamsErrCode))
		return
	}

	// 从缓存获取用户的回收站文件列表
	userRecycleFiles, err := cache.GetUserRecycleList(ctx, currentUserId)
	if errors.Is(err, redis.Nil) {
		// 缓存未命中，从数据库查询
		userRecycleFiles, err = cache.RefreshUserRecycleList(ctx, currentUserId)
		if err != nil {
			zap.S().Errorf("Failed to get user recycle list: %v", err)
			ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
			return
		}
	} else if err != nil {
		zap.S().Errorf("Failed to get user recycle list from cache: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	// 删除回收站文件系统中的所有文件
	for _, file := range userRecycleFiles {
		filePath := filepath.Join(config.PathCfg.RecycleBinPath, file.SourceIndexId, file.KBFileId)
		if err := wrench.RemoveContents(filePath); err != nil {
			zap.S().Errorf("Failed to delete file %s for user %s: %v", file.KBFileId, currentUserId, err)
			ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
			return
		}

		// 删除每个文件的缓存
		err = cache.DelRecycleInfo(ctx, file.KBFileId)
		if err != nil {
			zap.S().Errorf("Failed to delete recycle cache for file %s: %v", file.KBFileId, err)
		}
	}

	// 删除数据库中的用户回收站记录
	err = config.DB.Where("user_id = ?", currentUserId).Delete(&model.Recycle{}).Error
	if err != nil {
		zap.S().Errorf("Failed to delete recycle bin records for user %s: %v", currentUserId, err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	// 刷新用户的回收站列表缓存
	_, err = cache.RefreshUserRecycleList(ctx, currentUserId)
	if err != nil {
		zap.S().Errorf("Failed to refresh user recycle list cache: %v", err)
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, dto.Success())
}
