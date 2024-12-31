// -------------------------------------------------
// Package kb_apis
// Author: hanzhi
// Date: 2024/12/15
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
	"gorm.io/gorm"
	"net/http"
	"path/filepath"
)

// DeleteRecycleFile
// @Summary      删除回收站中的文件
// @Description  删除指定回收站中的文件，包括数据库记录和文件系统文件
// @ID           delete-recycle-file
// @Tags         recycle
// @Accept       json
// @Produce      json
// @Param        request  body      dto.RecycleRequest true "回收站文件删除请求体"
// @Success      200      {object}  dto.BaseResponse   "成功响应，返回success"
// @Failure      400      {object}  dto.BaseResponse   "参数错误(code:40000)"
// @Failure      401      {object}  dto.BaseResponse   "Token错误(code:40101)"
// @Failure      404      {object}  dto.BaseResponse   "文件记录未找到(code:40401)"
// @Failure      500      {object}  dto.BaseResponse   "服务器内部错误(code:50000)"
// @Router       /recycle/delete_file [post]
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

	// 从缓存验证回收站文件是否存在
	recycle, err := cache.GetRecycleInfo(ctx, req.KBFileId)
	if errors.Is(err, redis.Nil) {
		// 缓存未命中，从数据库查询
		recycle, err = cache.RefreshRecycleInfo(ctx, req.KBFileId)
		if err != nil {
			zap.S().Errorf("Failed to get recycle info: %v", err)
			ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
			return
		}
	} else if err != nil {
		zap.S().Errorf("Failed to get recycle from cache: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	if recycle == nil || recycle.KBFileId == "" || recycle.UserId != currentUserId {
		ctx.JSON(http.StatusNotFound, dto.Fail(dto.RecycleFileNotExistErrCode))
		return
	}

	// 开始删除操作
	var recycleFile model.Recycle
	tx := config.DB.Where("kb_file_id = ?", req.KBFileId).Delete(&recycleFile)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			zap.S().Errorf(
				"User with user_id %v source_index_name %v not found",
				currentUserId, req.KBFileId)
			ctx.JSON(
				http.StatusNotFound,
				dto.FailWithMessage(dto.RecycleFileNotExistErrCode, "try to rename the index but record note found."),
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

	// 更新缓存
	// 1. 删除回收站文件信息缓存
	err = cache.DelRecycleInfo(ctx, req.KBFileId)
	if err != nil {
		zap.S().Errorf("Failed to delete recycle cache: %v", err)
	}
	// 2. 刷新用户的回收站列表缓存
	_, err = cache.RefreshUserRecycleList(ctx, currentUserId)
	if err != nil {
		zap.S().Errorf("Failed to refresh user recycle list cache: %v", err)
	}

	ctx.JSON(http.StatusOK, dto.Success())
}
