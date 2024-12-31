// -------------------------------------------------
// Package kb_apis
// Author: hanzhi
// Date: 2024/12/13
// -------------------------------------------------

package kb_apis

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
	"os"
	"path/filepath"
)

/*
1. 获取request
2. 获取userId
3. 验证index存在（所以说redis还是得搞，格式就是user_id/index_name/把）  redis最后补上
4. 在kb_file表新建表项（事务
5. 创建对应的文件夹（kb_file_id/images, kb_file_id/file.md)

*/

// CreateKBFile
// @Summary     创建知识库文件
// @Description 创建一个空白的文件，要求是在知识库路径里创建一个名为id的文件夹和 images目录以及.md文件
// @ID			 create_KB_file
// @Tags 		 index
// @Accept       json
// @Produce      json
// @Param		request		body		dto.KBFileCreateRequest true "文档请求体"
// @Success		 200	{object} 		dto.BaseResponse "成功"
// @Failure		 400	{object} 		dto.BaseResponse "KBFileName无效等参数问题(40300)"
// @Failure		 401	{object} 		dto.BaseResponse "未授权，用户未登录(40101)"
// @Failure		 409	{object} 		dto.BaseResponse "知识库不存在（40201）"
// @Failure      500	{object} 		dto.BaseResponse "服务器内部错误(code:50000)"
// @Router		 /index/create_file [post]
func CreateKBFile(ctx *gin.Context) {
	var req dto.KBFileCreateRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		zap.S().Debugf("mark 1")
		zap.S().Debugf("params : %+v", req)
		ctx.JSON(http.StatusBadRequest, dto.Fail(dto.ParamsErrCode))
		return
	}

	// 验证KBName的有效性
	if !wrench.ValidateKBName(req.KBFileName) {
		zap.S().Debugf("Unaccept index name %v", req.KBFileName)
		ctx.JSON(http.StatusBadRequest, dto.Fail(dto.KBFileNameErrCode))
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

	// 创建一个KBFile
	KBFileId := wrench.IdGenerator()
	KBFileNew := model.KBFile{
		UserId:     currentUserId,
		IndexId:    req.IndexId,
		KBFileId:   KBFileId,
		KBFileName: req.KBFileName,
	}
	// 检查是否存在同名的文件（kbFileId）
	var kbFile = model.KBFile{}
	tx := config.DB.Model(&kbFile).Where("kb_file_name = ? AND index_id = ? ",
		KBFileNew.KBFileName, KBFileNew.IndexId).First(&kbFile) // fixme: 这里的代码有问题，需要解决
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

	// 开始事务
	tx = config.DB.Begin()
	err = tx.Create(&KBFileNew).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		zap.S().Errorf("Create kbfile %v, Error: %v", KBFileNew.KBFileName, tx.Error)
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}
	// 文件系统
	kbPath := config.PathCfg.KnowledgeBasePath
	kbDirPath := filepath.Join(kbPath, index.IndexId, KBFileId)
	err = os.Mkdir(kbDirPath, os.ModePerm)
	if err != nil {
		zap.S().Errorf("Create KBFile Dir %s Error: %v\n", kbDirPath, err)
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, dto.FailWithMessage(dto.InternalErrCode, "create file dir error"))
		return
	}
	// 生成文件夹 images
	imagesDirPath := filepath.Join(kbDirPath, "images")
	err = os.Mkdir(imagesDirPath, os.ModePerm)
	if err != nil {
		zap.S().Errorf("Create Image Dir %s Error: %v\n", imagesDirPath, err)
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, dto.FailWithMessage(dto.InternalErrCode, "create file dir error"))
		return
	}
	// 创建md文件
	mdPath := filepath.Join(kbDirPath, req.KBFileName+".md")
	file, err := os.Create(mdPath)
	if err != nil {
		zap.S().Errorf("Create File %s Error: %v", mdPath, err)
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, dto.FailWithMessage(dto.InternalErrCode, "create file dir error"))
		return
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			zap.S().Errorf("Close File Error %s Error: %v\n", mdPath, err)
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, dto.FailWithMessage(dto.InternalErrCode, "create file dir error"))
			return
		}
	}(file)

	// 不需要对ES做操作，因为是空的

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		zap.S().Errorf("Failed to commit transaction, err: %v", err)
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
