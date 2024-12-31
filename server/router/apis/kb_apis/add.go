// -------------------------------------------------
// Package kb_apis
// Author: hanzhi
// Date: 2024/12/15
// -------------------------------------------------

package kb_apis

import (
	"errors"
	"gcnote/server/ability/convert"
	"gcnote/server/ability/embeds"
	"gcnote/server/ability/search_engine"
	"gcnote/server/ability/splitter"
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
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// AddKBFile
// @Summary		导入文件
// @Description	 将文件导入对应的知识库
// @ID			 add_kb_file
// @Tags		 index
// @Accept       json
// @Produce      json
// @Param		request		body		dto.KBFileAddRequest true "文档添加请求体"
// @Success		 200	{object} 		dto.BaseResponse "成功"
// @Failure		 400	{object} 		dto.BaseResponse "KBFileName无效等参数问题(40300)"
// @Failure		 401	{object} 		dto.BaseResponse "未授权，用户未登录(40101)"
// @Failure		 409	{object} 		dto.BaseResponse "知识库不存在（40201）"
// @Failure      500	{object} 		dto.BaseResponse "服务器内部错误(code:50000)"
// @Router		 /index/add_file [post]
func AddKBFile(ctx *gin.Context) {
	var req dto.KBFileAddRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Fail(dto.ParamsErrCode))
		zap.S().Debugf("mark 1: params: %v", req)
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
	index, err := cache.GetIndexInfo(ctx, req.IndexId)
	if errors.Is(err, redis.Nil) {
		// 缓存未命中，从数据库查询
		index, err = cache.RefreshIndexInfo(ctx, req.IndexId)
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
		ctx.JSON(http.StatusConflict, dto.FailWithMessage(dto.IndexNotExistErrCode,
			"知识库"+req.IndexId+"不存在"))
		return
	}

	// 正式开始文件处理
	// 获取文件名和文件扩展名
	fileName := req.File.Filename
	fileExt := filepath.Ext(fileName)
	// 获取上传的文件
	file, err := req.File.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			dto.FailWithMessage(dto.ParamsErrCode, "param file error."))
		return
	}
	defer func(file multipart.File) {
		err = file.Close()
		if err != nil {
			zap.S().Errorf("Close File Error %s Error: %v\n", fileName, err)
			ctx.JSON(http.StatusInternalServerError, dto.FailWithMessage(dto.InternalErrCode, "create file dir error"))
			return
		}
	}(file)
	// 保存到本地tmp目录里面先
	tmpDirPath := config.PathCfg.TempDirPath
	tmpFileDirPath := filepath.Join(tmpDirPath, wrench.IdGenerator())
	tmpFilePath := filepath.Join(tmpFileDirPath, fileName)
	// 保存到这里
	dst := filepath.Join(tmpFileDirPath, fileName)
	err = ctx.SaveUploadedFile(req.File, dst)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.KBFileAddFileErrCode))
	}
	// 导入没问题 那么我就后续处理了
	// ------------------------------------

	// 开始sql操作
	KBFileNew := model.KBFile{
		UserId:     currentUserId,
		KBFileId:   wrench.IdGenerator(),
		KBFileName: strings.TrimSuffix(fileName, fileExt),
		IndexId:    req.IndexId,
	}
	// 检查是否存在同名的文件（kbFileId）
	var kbFile = model.KBFile{}
	tx := config.DB.Model(&KBFileNew).Where("kb_file_name = ? AND index_id = ?",
		KBFileNew.KBFileName, KBFileNew.IndexId).First(&kbFile)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		zap.S().Debugf("kb_file_name: %s is exist", KBFileNew.KBFileName)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
		// 继续往下走
	}
	if tx.RowsAffected != 0 { // 行数不为0 ，说明已经存在了
		// 如果文件已经存在了，那就给他重命名一下啊？
		// 算了，把麻烦交给用户把
		ctx.JSON(http.StatusConflict, dto.Fail(dto.KBFileExistErrCode))
		return
	}
	// -------------------------------------------------
	// 一切问题都不存在了，那我就开始操作。
	// 开始事务，处理
	// 开始事务
	tx = config.DB.Begin()
	if tx.Error != nil {
		zap.S().Errorf("Failed to begin transaction, err:%v", tx.Error)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}
	err = tx.Create(&KBFileNew).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		zap.S().Errorf("Create kbfile %v, Error: %v", KBFileNew.KBFileName, tx.Error)
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}
	// 文件系统操作
	kbPath := config.PathCfg.KnowledgeBasePath
	kbDirPath := filepath.Join(kbPath, KBFileNew.IndexId, KBFileNew.KBFileId)
	err = os.Mkdir(kbDirPath, os.ModePerm)
	if err != nil {
		zap.S().Errorf("Create KBFile Dir %s Error: %v\n", kbDirPath, err)
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, dto.FailWithMessage(dto.InternalErrCode, "create file dir error"))
		return
	}
	// 文件导入操作
	var mdString string
	_, mdString, err = convert.AutoConvert(tmpFilePath, kbDirPath, fileExt) // 第二个 mdString
	if err != nil {
		zap.S().Errorf("Convert File Error: %v", err)
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	// 删除临时文件夹 tmpFileDirPath
	err = wrench.RemoveContents(tmpFileDirPath)
	if err != nil {
		zap.S().Errorf("Failed to remove temp dir, err: %v", err)
	}

	// 将mdString切片，并提交到es中
	chunks := splitter.SplitMarkdown(mdString, 512)
	docList := splitter.Chunk2Doc(chunks, KBFileNew.KBFileId, KBFileNew.IndexId)
	embedList, err := embeds.RandEmbedding(docList) // fixme 之后换成正常的Embedding服务
	err = search_engine.AddDocuments(config.ElasticClient, "gcnote-"+KBFileNew.IndexId, docList, embedList)
	if err != nil {
		zap.S().Errorf("Failed to add document into index, err: %v", err)
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		zap.S().Errorf("Failed to commit transaction, err: %v", err)
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	// 更新缓存
	// 1. 设置kb文件信息缓存
	err = cache.SetKBInfo(ctx, KBFileNew)
	if err != nil {
		zap.S().Errorf("Failed to set kb file cache: %v", err)
	}
	// 2. 刷新index的kb文件列表缓存
	_, err = cache.RefreshIndexKBList(ctx, req.IndexId)
	if err != nil {
		zap.S().Errorf("Failed to refresh index kb list cache: %v", err)
	}
	// 3. 刷新用户的最近访问kb列表缓存
	_, err = cache.RefreshRecentKBList(ctx, currentUserId)
	if err != nil {
		zap.S().Errorf("Failed to refresh user recent kb list cache: %v", err)
	}

	zap.S().Infof("Create file name %v done.", KBFileNew.KBFileName)
	ctx.JSON(http.StatusOK, dto.Success())
}
