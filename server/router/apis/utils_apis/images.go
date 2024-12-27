// -------------------------------------------------
// Package utils_apis
// Author: hanzhi
// Date: 2024/12/27
// -------------------------------------------------

package utils_apis

import (
	"gcnote/server/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func GetImage(ctx *gin.Context) {
	// 从URL路径中获取参数
	indexID := ctx.Param("index_id")
	kbFileID := ctx.Param("kb_file_id")
	imageNameWithExt := ctx.Param("image_name")

	// 分离文件名和扩展名
	parts := strings.Split(imageNameWithExt, ".")
	if len(parts) != 2 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image name format"})
		return
	}

	imageName := parts[0]
	ext := parts[1]

	// 构建图片文件的完整路径
	imagePath := filepath.Join(config.PathCfg.KnowledgeBasePath, indexID, kbFileID, "images", imageName+"."+ext)

	// 读取图片文件
	imageBytes, err := os.ReadFile(imagePath)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	// 根据文件扩展名设置Content-Type
	var contentType string
	switch ext {
	case "jpg", "jpeg":
		contentType = "image/jpeg"
	case "png":
		contentType = "image/png"
	case "gif":
		contentType = "image/gif"
	case "bmp":
		contentType = "image/bmp"
	default:
		contentType = "application/octet-stream"
	}

	// 返回图片数据
	ctx.Header("Content-Type", contentType)
	ctx.Data(http.StatusOK, contentType, imageBytes)
}
