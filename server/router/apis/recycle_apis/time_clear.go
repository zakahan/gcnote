// -------------------------------------------------
// Package recycle_apis
// Author: hanzhi
// Date: 2024/12/17
// -------------------------------------------------

package recycle_apis

import (
	"gcnote/server/config"
	"gcnote/server/dto"
	"gcnote/server/model"
	"gcnote/server/router/wrench"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"path/filepath"
	"time"
)

// CleanupOldRecycleFiles
// @Summary      定时清理回收站
// @Description  删除回收站中所有超过30天未更新的文件记录及其对应的物理文件
// @ID           cleanup-old-recycle-files
// @Tags         recycle
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.BaseResponse  "成功响应，返回success"
// @Failure      500  {object}  dto.BaseResponse  "服务器内部错误(code:50000)"
// @Router       /recycle/cleanup [post]
func CleanupOldRecycleFiles(ctx *gin.Context) {
	// 获取30天前的时间点
	cutoffDate := time.Now().AddDate(0, 0, -30)
	//cutoffDate := time.Now().Add(-5 * time.Minute)		// 测试用的

	// 查询需要清理的回收站记录
	var oldRecycleFiles []model.Recycle
	err := config.DB.Where("updated_at < ?", cutoffDate).Find(&oldRecycleFiles).Error
	if err != nil {
		zap.S().Errorf("Failed to query old recycle bin files: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	// 删除文件系统中的对应文件
	for _, file := range oldRecycleFiles {
		filePath := filepath.Join(config.PathCfg.RecycleBinPath, file.SourceIndexId, file.KBFileId)
		if err := wrench.RemoveContents(filePath); err != nil {
			zap.S().Errorf("Failed to delete file %s for cleanup: %v", file.KBFileId, err)
			ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
			return
		}
	}

	// 删除数据库中的过期记录
	err = config.DB.Where("updated_at < ?", cutoffDate).Delete(&model.Recycle{}).Error
	if err != nil {
		zap.S().Errorf("Failed to delete old recycle bin records: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, dto.Success())
}
