// -------------------------------------------------
// Package kb_apis
// Author: hanzhi
// Date: 2024/12/30
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

// RecentDocs
// @Summary 获取最近文档
// @Description 获取用户最近修改或创建的文件
// @ID         recent-docs
// @Tags       index
// @Accept     json
// @Produce    json
// @Param      request  body      dto.RecentDocsRequest true "请求体"
// @Success    200      {object}  dto.BaseResponse{data=[]dto.RecentDocResponse} "成功响应，返回文档信息列表"
// @Failure    400      {object}  dto.BaseResponse "参数错误(code:40000)"
// @Failure    401      {object}  dto.BaseResponse "Token错误(code:40101)"
// @Failure    500      {object}  dto.BaseResponse "服务器内部错误(code:50000)"
// @Router     /index/recent_docs [post]
func RecentDocs(ctx *gin.Context) {
	var req dto.RecentDocsRequest
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

	// 查询最近的文件
	var kbFileList []model.KBFile
	tx := config.DB.Model(&model.KBFile{}).Where("user_id = ?", currentUserId)

	switch req.Mode {
	case "modified":
		// 按最近修改时间排序
		tx = tx.Order("updated_at DESC")
	case "created":
		// 按最近创建时间排序
		tx = tx.Order("created_at DESC")
	default:
		ctx.JSON(http.StatusBadRequest, dto.Fail(dto.ParamsErrCode))
		return
	}

	if err := tx.Limit(20).Find(&kbFileList).Error; err != nil {
		zap.S().Errorf("failed to retrieve recent docs: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	// 检查结果是否为空
	if len(kbFileList) == 0 {
		zap.S().Infof("No records found for user_id: %s, mode: %s", currentUserId, req.Mode)
		ctx.JSON(http.StatusOK, dto.SuccessWithData([]dto.RecentDocResponse{}))
		return
	}

	// 格式化响应数据
	var responseList []dto.RecentDocResponse
	for _, file := range kbFileList {
		responseList = append(responseList, dto.RecentDocResponse{
			FileId:     file.KBFileId,
			IndexId:    file.IndexId,
			FileName:   file.KBFileName,
			UserId:     file.UserId,
			CreatedAt:  file.CreatedAt,
			ModifiedAt: file.UpdatedAt,
		})
	}

	// 返回结果
	ctx.JSON(http.StatusOK, dto.SuccessWithData(responseList))
}
