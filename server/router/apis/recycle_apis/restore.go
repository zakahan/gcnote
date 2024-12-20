// -------------------------------------------------
// Package recycle_apis
// Author: hanzhi
// Date: 2024/12/17
// -------------------------------------------------

package recycle_apis

import (
	"gcnote/server/ability/embeds"
	"gcnote/server/ability/search_engine"
	"gcnote/server/ability/splitter"
	"gcnote/server/config"
	"gcnote/server/dto"
	"gcnote/server/model"
	"gcnote/server/router/wrench"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"net/http"
	"os"
	"path/filepath"
)

// RestoreRecycleFile
// @Summary      从回收站恢复文件
// @Description  将回收站中的文件恢复到正常知识库系统中，并确保知识库存在
// @ID           restore-recycle-file
// @Tags         recycle
// @Accept       json
// @Produce      json
// @Param        request  body      dto.RecycleRequest true "回收站文件恢复请求体"
// @Success      200      {object}  dto.BaseResponse   "成功响应，返回success"
// @Failure      400      {object}  dto.BaseResponse   "参数错误(code:40000)"
// @Failure      401      {object}  dto.BaseResponse   "Token错误(code:40101)"
// @Failure      404      {object}  dto.BaseResponse   "知识库或文件不存在(code:40201)"
// @Failure      500      {object}  dto.BaseResponse   "服务器内部错误(code:50000)"
// @Router       /recycle/restore [post]
func RestoreRecycleFile(ctx *gin.Context) {
	// 解析请求体
	var req dto.RecycleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		zap.S().Debugf("Invalid parameters: %+v", req)
		ctx.JSON(http.StatusBadRequest, dto.Fail(dto.ParamsErrCode))
		return
	}

	// 验证用户身份
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
		ctx.JSON(http.StatusBadRequest, dto.Fail(dto.ParamsErrCode))
		return
	}

	// 检查文件是否存在于回收站
	var recycleFile model.Recycle
	if err := config.DB.Where("kb_file_id = ?", req.KBFileId).First(&recycleFile).Error; err != nil {
		zap.S().Errorf("Recycle file not found: %v", req.KBFileId)
		ctx.JSON(http.StatusNotFound, dto.Fail(dto.IndexNotExistErrCode))
		return
	}

	// 检查对应的知识库是否存在
	// 不负责检查是否存在知识库，应该是先检查是否存在这个知识库，如果不存在，前端调用创建操作。不要耦合在这里。

	// 文件从回收站恢复到知识库目录
	recycleFilePath := filepath.Join(config.PathCfg.RecycleBinPath, recycleFile.SourceIndexId, recycleFile.KBFileId)
	knowledgeBaseFilePath := filepath.Join(config.PathCfg.KnowledgeBasePath, recycleFile.SourceIndexId, recycleFile.KBFileId)

	// 确保目标目录存在
	if err := os.MkdirAll(filepath.Dir(knowledgeBaseFilePath), os.ModePerm); err != nil {
		zap.S().Errorf("Failed to create knowledge base directory: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	// 移动文件
	src := recycleFilePath
	dst := knowledgeBaseFilePath
	// 复制文件到要求的路径
	err := wrench.CopyDir(src, dst)
	if err != nil {
		zap.S().Errorf("Failed to create recycle dir, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	// 删除原来的文件夹
	err = wrench.RemoveContents(src)
	if err != nil {
		zap.S().Errorf("Delete kb id :%v err:%v", req.KBFileId, err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	// 在知识库中创建文件记录
	newFile := model.KBFile{
		IndexId:    recycleFile.SourceIndexId,
		KBFileId:   recycleFile.KBFileId,
		KBFileName: recycleFile.KBFileName,
	}
	if err := config.DB.Create(&newFile).Error; err != nil {
		zap.S().Errorf("Failed to create KBFile record: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	// 删除回收站中的记录
	if err := config.DB.Delete(&recycleFile).Error; err != nil {
		zap.S().Errorf("Failed to delete recycle record: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	// 恢复到ElasticSearch中
	filePath := filepath.Join(dst, recycleFile.KBFileName+".md")
	// 读取filePath为文件字符串
	content, err := os.ReadFile(filePath)
	if err != nil {
		zap.S().Errorf("Failed to read file: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}
	mdString := string(content)
	chunks := splitter.SplitMarkdown(mdString, 512)
	docList := splitter.Chunk2Doc(chunks, recycleFile.KBFileId, recycleFile.SourceIndexId)
	embedList, err := embeds.RandEmbedding(docList) // fixme 之后换成正常的Embedding服务
	err = search_engine.AddDocuments(config.ElasticClient, "gcnote-"+recycleFile.SourceIndexId, docList, embedList)
	if err != nil {
		zap.S().Errorf("Failed to add document into index, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}
	// 返回成功响应
	ctx.JSON(http.StatusOK, dto.Success())
}
