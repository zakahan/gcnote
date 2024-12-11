// -------------------------------------------------
// Package kb_apis
// Author: hanzhi
// Date: 2024/12/11
// -------------------------------------------------

package index_apis

import (
	"errors"
	"fmt"
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
func CreateIndex(ctx *gin.Context) {
	var req dto.IndexCRUDRequset
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		zap.S().Debugf("mark 1")
		zap.S().Debugf("params : %+v", req)
		ctx.JSON(http.StatusOK, dto.Fail(dto.ParamsErrCode))
		return
	}
	claims, exists := ctx.Get("claims")
	if !exists {
		zap.S().Infof("Unable to get the claims")
		ctx.JSON(http.StatusOK, dto.Fail(dto.UserTokenErrCode))
		zap.S().Debugf("mark 2")
		return
	}
	zap.S().Debugf("claims: %v", claims)
	currentUser := claims.(jwt.MapClaims)
	currentUserId := currentUser["sub"].(string)
	if currentUserId == "" {
		zap.S().Debugf("currentUserId is empty")
		ctx.JSON(http.StatusOK, dto.Fail(dto.ParamsErrCode))
		return
	}
	// 获取了UserId之后，创建知识库
	indexId := wrench.IdGenerator()
	indexNew := model.Index{
		IndexId:   indexId,
		IndexName: req.IndexName,
		UserId:    currentUserId,
	}
	// 查看是否有同名的，如果有，那么就无法创建知识库
	indexSearch := model.Index{}

	tx := config.DB.Where("index_name = ?", indexSearch.IndexName).First(&indexSearch)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		zap.S().Errorf("Create index %v, err: %v", indexNew.IndexName, tx.Error)
		ctx.JSON(http.StatusOK, dto.Fail(dto.InternalErrCode))
		return
	}
	if tx.RowsAffected != 0 { // 行数不为零，即存在同名知识库 那就不允许
		ctx.JSON(http.StatusOK, dto.Fail(dto.IndexExistErrCode))
		return
	}
	// 成功创建
	tx = config.DB.Create(&indexNew)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		zap.S().Errorf("Create index %v, err: %v", indexNew.IndexName, tx.Error)
		ctx.JSON(http.StatusOK, dto.Fail(dto.InternalErrCode))
		return
	}
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
			// 回滚mysql表
			ctx.JSON(http.StatusOK, dto.FailWithMessage(dto.InternalErrCode, "create file dir error"))
			return
		}
	} else if err != nil {
		zap.S().Errorf("检查路径 %s 时出错: %v\n", kbPath, err)
		// 回滚mysql表
		ctx.JSON(http.StatusOK, dto.FailWithMessage(dto.InternalErrCode, "create file dir error"))
		return
	}

	// 创建名为 xy 的文件夹
	xyPath := filepath.Join(kbPath, indexId)
	if err := os.Mkdir(xyPath, os.ModePerm); err != nil {
		fmt.Printf("创建文件夹 %s 失败: %v\n", xyPath, err)
		ctx.JSON(http.StatusOK, dto.FailWithMessage(dto.InternalErrCode, "create file dir error"))
		return
	}
	zap.S().Infof("Create index %v done.", indexNew.IndexName)
	ctx.JSON(http.StatusOK, dto.Success())
}

func rollBackCreate(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, dto.FailWithMessage(dto.InternalErrCode, "create file dir error"))
	return
}
