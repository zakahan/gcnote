// -------------------------------------------------
// Package kb_apis
// Author: hanzhi
// Date: 2024/12/11
// -------------------------------------------------

package index_apis

import (
	"errors"
	"gcnote/server/ability/search_engine"
	"gcnote/server/cache"
	"gcnote/server/config"
	"gcnote/server/dto"
	"gcnote/server/model"
	"gcnote/server/router/wrench"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"os"
	"path/filepath"
)

// CreateIndex
// @Summary 创建知识库
// @Description 创建知识库(MYSQL + 文件系统已完成)
// @ID			create-index
// @Tags		index
// @Accept		json
// @Produce		json
// @Param		request		body		dto.IndexCreateRequest true "登录请求体"
// @Success 	200			{object} 	dto.BaseResponse	"成功响应，返回success"
// @Failure		400			{object}	dto.BaseResponse	"参数错误(code:40000)"
// @Failure		401			{object}	dto.BaseResponse	"Token错误(code:40101)"
// @Failure		409			{object}	dto.BaseResponse	"该知识库已存在(code:40200)"
// @Failure		500			{object}	dto.BaseResponse	"服务器内部错误(code:50000)"
// @Router 		/index/create_index [post]
func CreateIndex(ctx *gin.Context) {
	var req dto.IndexCreateRequest
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
	// 验证IndexName的有效性
	if !wrench.ValidateIndexName(req.IndexName) {
		zap.S().Debugf("Unaccept index name %v", req.IndexName)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.IndexNameErrCode))
	}

	zap.S().Debugf("claims: %v", claims)
	currentUser := claims.(jwt.MapClaims)
	currentUserId := currentUser["sub"].(string)
	if currentUserId == "" {
		zap.S().Debugf("currentUserId is empty")
		ctx.JSON(http.StatusBadRequest, dto.Fail(dto.ParamsErrCode))
		return
	}
	// 获取了IndexId之后，创建知识库
	indexId := wrench.IdGenerator()
	indexNew := model.Index{
		IndexId:   indexId,
		IndexName: req.IndexName,
		UserId:    currentUserId,
	}
	// 查看是否有同名的，如果有，那么就无法创建知识库（
	// 这个有问题，是允许重名的！不允许的应该是同一个用户目录下，不允许重名！）
	indexSearch := model.Index{}

	tx := config.DB.Where("index_name = ? AND user_id = ? ",
		req.IndexName, currentUserId).First(&indexSearch)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		zap.S().Errorf("Create index %v, err: %v", indexNew.IndexName, tx.Error)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}
	if tx.RowsAffected != 0 { // 行数不为零，即存在同名知识库 那就不允许
		ctx.JSON(http.StatusConflict, dto.Fail(dto.IndexExistErrCode))
		return
	}
	// ------------------------------------------------------------------
	// 开始事务
	tx = config.DB.Begin()
	if tx.Error != nil {
		zap.S().Errorf("Failed to begin transaction, err: %v", tx.Error)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}
	// 成功创建
	err = tx.Create(&indexNew).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		zap.S().Errorf("Create index %v, err: %v", indexNew.IndexName, tx.Error)
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}
	// ------------------------------------------------------------------
	// 文件系统创建
	kbPath := config.PathCfg.KnowledgeBasePath
	// 检查路径是否存在
	_, err = os.Stat(kbPath)
	if os.IsNotExist(err) {
		zap.S().Infof("路径 %v 不存在，执行创建操作", kbPath)
		err = os.MkdirAll(kbPath, os.ModePerm)
		zap.S().Infof("路径 %v 创建成功。", kbPath)
		if err != nil {
			zap.S().Errorf("创建路径 %v 失败， err: %v", kbPath, err)
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, dto.FailWithMessage(dto.InternalErrCode, "create file dir error"))
			return
		}
	} else if err != nil {
		zap.S().Errorf("检查路径 %s 时出错: %v\n", kbPath, err)
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, dto.FailWithMessage(dto.InternalErrCode, "create file dir error"))
		return
	}

	// 创建名为 xy 的文件夹
	xyPath := filepath.Join(kbPath, indexId)
	err = os.Mkdir(xyPath, os.ModePerm)
	if err != nil {
		zap.S().Errorf("创建文件夹 %s 失败: %v\n", xyPath, err)
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, dto.FailWithMessage(dto.InternalErrCode, "create file dir error"))
		return
	}
	// ------------------------------------------------------------------
	// 创建es索引
	err, code := search_engine.IndexExist(config.ElasticClient, "gcnote-"+indexId)
	if err != nil || code != 404 {
		if err != nil {
			zap.S().Errorf("Failed to check exist of index, err: %v", err)
		}
		if code != 404 {
			zap.S().Errorf("index is already exist, code == %v", code)
		}
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}
	// 只有404才创建
	err = search_engine.IndexCreate(config.ElasticClient, "gcnote-"+indexId)
	if err != nil {
		zap.S().Errorf("Failed to create index, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	// ----------------------------------------------------

	err = tx.Commit().Error
	if err != nil {
		zap.S().Errorf("Failed to commit transaction, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	// 设置缓存
	err = cache.SetIndexInfo(ctx, indexNew)
	if err != nil {
		zap.S().Errorf("Failed to set index cache, err: %v", err)
	}
	// 刷新用户的index列表缓存
	_, err = cache.RefreshUserIndexList(ctx, currentUserId)
	if err != nil {
		zap.S().Errorf("Failed to refresh user index list cache, err: %v", err)
	}

	zap.S().Infof("Create index %v done.", indexNew.IndexName)
	ctx.JSON(http.StatusOK, dto.Success())
}
