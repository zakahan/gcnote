// -------------------------------------------------
// Package share_apis
// Author: hanzhi
// Date: 2025/1/5
// -------------------------------------------------

package share_apis

import (
	"gcnote/server/ability/splitter"
	"gcnote/server/config"
	"gcnote/server/dto"
	"gcnote/server/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"path/filepath"
)

func ReadFile(ctx *gin.Context) {
	var req dto.ShareFileReadRequest
	if err := ctx.ShouldBind(&req); err != nil {
		zap.S()
		ctx.JSON(http.StatusBadRequest, dto.Fail(dto.ParamsErrCode))
		zap.S().Debugf("mark 1: params: %v", req)
		return
	}

	// 检查当前req.shareFileId是否存在
	if req.ShareFileId == "" {
		ctx.JSON(http.StatusBadRequest, dto.Fail(dto.ParamsErrCode))
		zap.S().Debugf("mark 2: params: %v", req)
		return
	}
	// sql查询
	shareFile := model.ShareFile{}
	tx := config.DB.Where("share_file_id = ? ", req.ShareFileId).First(&shareFile)
	if tx.Error != nil {
		ctx.JSON(http.StatusBadRequest, dto.Fail(dto.ShareFileNotExistErrCode))
		zap.S().Debugf("mark 3: params: %v", req)
		return
	}
	if shareFile.Password != req.Password {
		ctx.JSON(http.StatusBadRequest, dto.Fail(dto.SharePasswordErrCode))
		zap.S().Debugf("mark 4: params: %v", req)
		return
	}

	fileDir := filepath.Join(config.PathCfg.ShareFileDirPath, req.ShareFileId)
	filePath := filepath.Join(fileDir, shareFile.FileName+".md")
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
	resultData := splitter.ChunkRead(chunks, config.PathCfg.ImageServerURL, "share", req.ShareFileId)

	srResponse := dto.ShareReadResponse{
		Title:   shareFile.FileName,
		Content: resultData,
	}

	ctx.JSON(http.StatusOK, dto.SuccessWithData(srResponse))
}
