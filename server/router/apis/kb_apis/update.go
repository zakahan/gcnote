// -------------------------------------------------
// Package kb_apis
// Author: hanzhi
// Date: 2024/12/30
// -------------------------------------------------

package kb_apis

import (
	"bytes"
	"gcnote/server/ability/embeds"
	"gcnote/server/ability/search_engine"
	"gcnote/server/ability/splitter"
	"gcnote/server/config"
	"gcnote/server/dto"
	"gcnote/server/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"io"
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

	// 获取了indexId, kbFileId, UserId
	// 查看这个是否存在
	kbFile := model.KBFile{
		UserId:   currentUserId,
		IndexId:  indexId,
		KBFileId: kbFileId,
	}

	// 查询 KBFile 是否存在
	if err := config.DB.Where(&kbFile).First(&kbFile).Error; err != nil {
		ctx.JSON(http.StatusNotFound, dto.Fail(dto.RecordNotFoundErrCode))
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
	defer src.Close()
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
	ctx.JSON(http.StatusOK, dto.Success())

}
