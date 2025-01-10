// -------------------------------------------------
// Package kb_apis
// Author: hanzhi
// Date: 2024/12/30
// -------------------------------------------------

package kb_apis

import (
	"bytes"
	"errors"
	"gcnote/server/ability/embeds"
	"gcnote/server/ability/search_engine"
	"gcnote/server/ability/splitter"
	"gcnote/server/cache"
	"gcnote/server/config"
	"gcnote/server/dto"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func UpdateKBFile(ctx *gin.Context) {
	indexId := ctx.PostForm("index_id")
	kbFileId := ctx.PostForm("kb_file_id")
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Fail(dto.ParamsErrCode))
		return
	}
	// 获取userId
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

	// 先从缓存验证index是否存在
	index, err := cache.GetIndexInfo(ctx, indexId)
	if errors.Is(err, redis.Nil) {
		// 缓存未命中，从数据库查询
		index, err = cache.RefreshIndexInfo(ctx, indexId)
		if err != nil {
			zap.S().Errorf("Failed to get index info: %v", err)
			ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
			return
		}
	} else if err != nil {
		zap.S().Errorf("Failed to get index from cache: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	if index == nil || index.IndexId == "" || index.UserId != currentUserId {
		ctx.JSON(http.StatusNotFound, dto.Fail(dto.IndexNotExistErrCode))
		return
	}

	// 从缓存验证kb文件是否存在
	kbFile, err := cache.GetKBInfo(ctx, kbFileId)
	if errors.Is(err, redis.Nil) {
		// 缓存未命中，从数据库查询
		kbFile, err = cache.RefreshKBInfo(ctx, kbFileId)
		if err != nil {
			zap.S().Errorf("Failed to get kb file info: %v", err)
			ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
			return
		}
	} else if err != nil {
		zap.S().Errorf("Failed to get kb file from cache: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	if kbFile == nil || kbFile.KBFileId == "" || kbFile.IndexId != indexId || kbFile.UserId != currentUserId {
		ctx.JSON(http.StatusNotFound, dto.Fail(dto.KBFileNotExistErrCode))
		return
	}

	// 首先清空es相应的文件
	err = search_engine.DeleteByTerm(config.ElasticClient, "gcnote-"+kbFile.IndexId, "kb_file_id", kbFile.KBFileId)
	if err != nil {
		zap.S().Errorf("Failed to delete the document in elasticsearch index, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}
	// 然后是替换文件
	kbDirPath := filepath.Join(config.PathCfg.KnowledgeBasePath, kbFile.IndexId,
		kbFile.KBFileId)
	kbFilePath := filepath.Join(kbDirPath, kbFile.KBFileName+".md")

	// 先切片修改一下文件
	// 读取file为字符串
	src, err := file.Open()
	if err != nil {
		zap.S().Errorf("Failed to read file: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}
	defer func(src multipart.File) {
		err = src.Close()
		if err != nil {
			zap.S().Errorf("Close File Error Error: %v\n", err)
		}
	}(src)
	buffer := new(bytes.Buffer)
	_, err = io.Copy(buffer, src)
	if err != nil {
		zap.S().Errorf("Failed to read file: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
	}
	content := buffer.String()

	// 处理字符串为切片
	chunks := splitter.SplitMarkdownEasy(content)

	mdString, err := splitter.ChunkReadReverse(chunks, config.PathCfg.ImageServerURL, indexId, kbFileId)

	// 将mdString切片，并提交到es中
	chunks = splitter.SplitMarkdown(mdString, 512)
	docList := splitter.Chunk2Doc(chunks, kbFile.KBFileId, kbFile.IndexId)
	embedList, err := embeds.RandEmbedding(docList) // fixme 之后换成正常的Embedding服务
	err = search_engine.AddDocuments(config.ElasticClient, "gcnote-"+kbFile.IndexId, docList, embedList)
	if err != nil {
		zap.S().Errorf("Failed to add document into index, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	// 删除原始文件
	err = os.Remove(kbFilePath)
	if err != nil {
		zap.S().Errorf("Failed to delete the file, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	// 将文件保存到这里(kbFilePath)
	err = os.MkdirAll(kbDirPath, os.ModePerm)
	if err != nil {
		zap.S().Errorf("Failed to create kb dir, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}
	err = os.WriteFile(kbFilePath, []byte(mdString), 0644)
	if err != nil {
		zap.S().Errorf("Failed to write file, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
	}

	// 更新 KBFile 的修改时间
	if err := config.DB.Model(&kbFile).Updates(map[string]interface{}{
		"UpdatedAt": time.Now(),
	}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	// 更新缓存
	// 1. 刷新kb文件信息缓存
	_, err = cache.RefreshKBInfo(ctx, kbFileId)
	if err != nil {
		zap.S().Errorf("Failed to refresh kb file cache: %v", err)
	}
	// 2. 刷新index的kb文件列表缓存
	_, err = cache.RefreshIndexKBList(ctx, indexId)
	if err != nil {
		zap.S().Errorf("Failed to refresh index kb list cache: %v", err)
	}
	// 3. 刷新用户的最近访问kb列表缓存
	_, err = cache.RefreshRecentKBList(ctx, currentUserId)
	if err != nil {
		zap.S().Errorf("Failed to refresh user recent kb list cache: %v", err)
	}

	ctx.JSON(http.StatusOK, dto.Success())
}
