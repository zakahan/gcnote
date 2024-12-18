// -------------------------------------------------
// Package kb_apis
// Author: hanzhi
// Date: 2024/12/18
// -------------------------------------------------

package kb_apis

import (
	"gcnote/server/config"
	"gcnote/server/dto"
	"gcnote/server/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"net/http"
)

// SearchKBFiles
// @Summary 搜索文件
// @Description 根据文件名和索引ID搜索文件
// @ID         search-kb_files
// @Tags       index
// @Accept     json
// @Produce    json
// @Param      request  body      KBFileCSRequest true "搜索请求体"
// @Success    200      {object}  dto.BaseResponse "成功响应，返回搜索结果"
// @Failure    400      {object}  dto.BaseResponse "参数错误(code:40000)"
// @Failure    401      {object}  dto.BaseResponse "Token错误(code:40101)"
// @Failure    500      {object}  dto.BaseResponse "服务器内部错误(code:50000)"
// @Router     /index/search_file [post]
func SearchKBFiles(ctx *gin.Context) {
	var req dto.KBFileSearchRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Fail(dto.ParamsErrCode))
		return
	}

	// 获取用户身份信息
	claims, exists := ctx.Get("claims")
	if !exists {
		zap.S().Infof("Unable to get the claims")
		ctx.JSON(http.StatusUnauthorized, dto.Fail(dto.UserTokenErrCode))
		return
	}

	currentUser := claims.(jwt.MapClaims)
	currentUserId := currentUser["sub"].(string)
	if currentUserId == "" {
		ctx.JSON(http.StatusBadRequest, dto.Fail(dto.ParamsErrCode))
		return
	}

	// 查询文件
	var kbFileList []model.KBFile
	tx := config.DB.Model(&model.KBFile{})

	if req.IsFuzzySearch {
		// 模糊查询
		tx = tx.Where("index_id = ? AND kb_file_name LIKE ?", req.IndexId, "%"+req.KBFileName+"%")
	} else {
		// 精确查询
		tx = tx.Where("index_id = ? AND kb_file_name = ?", req.IndexId, req.KBFileName)
	}

	if err := tx.Find(&kbFileList).Error; err != nil {
		zap.S().Errorf("failed to find kb_file list: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	// 返回结果
	ctx.JSON(http.StatusOK, dto.SuccessWithData(kbFileList))
}
