// -------------------------------------------------
// Package kb_apis
// Author: hanzhi
// Date: 2024/12/16
// -------------------------------------------------

package kb_apis

import (
	"gcnote/server/ability/splitter"
	"gcnote/server/config"
	"gcnote/server/dto"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"path/filepath"
)

func ReadFile(ctx *gin.Context) {
	var req dto.KBFileReadRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Fail(dto.ParamsErrCode))
		zap.S().Debugf("mark 1: params: %v", req)
		return
	}

	// 读取文件
	fileDir := filepath.Join(config.PathCfg.KnowledgeBasePath, req.IndexId, req.KBFileId)
	// 文件
	filePath := filepath.Join(fileDir, req.KBFileName+".md")
	// 首先检查fileDir和filePath是否存在
	if _, err := os.Stat(fileDir); os.IsNotExist(err) {
		zap.S().Debugf("fileDir: %s is not exist", fileDir)
		ctx.JSON(http.StatusNotFound, dto.Fail(dto.KBFileNotExistErrCode))
		return
	}
	// 读取文件，转为字符串
	content, err := os.ReadFile(filePath)
	if err != nil {
		zap.S().Errorf("Failed to read file: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}
	// 将content转为切片然后读取
	chunks := splitter.SplitMarkdownEasy(string(content))
	resultData := splitter.ChunkRead(chunks, config.PathCfg.ImageServerURL, req.IndexId, req.KBFileId)
	// 修改图片路径为url路径

	ctx.JSON(http.StatusOK, dto.SuccessWithData(resultData))
	return

}
