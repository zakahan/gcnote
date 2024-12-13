// -------------------------------------------------
// Package index_apis
// Author: hanzhi
// Date: 2024/12/11
// -------------------------------------------------

package index_apis

import (
	"errors"
	"gcnote/server/config"
	"gcnote/server/dto"
	"gcnote/server/model"
	"gcnote/server/router/wrench"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
)

// RenameIndex
// @Summary 知识库重命名
// @Description 重命名知识库(
// @ID			rename-index
// @Tags		index
// @Accept		json
// @Produce		json
// @Param		request		body		dto.IndexCRUDRequest true "登录请求体"
// @Success 	200			{object} 	dto.BaseResponse	"成功响应，返回success"
// @Failure		200			{object}	dto.BaseResponse	"参数错误(code:40000)"
// @Failure		200			{object}	dto.BaseResponse	"新索引名称不符合规范(code:40202)"
// @Failure		200			{object}	dto.BaseResponse	"Token错误(code:40101)"
// @Failure		200			{object}	dto.BaseResponse	"原索引记录不存在(code:40201)"
// @Failure		200			{object}	dto.BaseResponse	"服务器内部错误(code:50000)"
// @Router 		/index/rename [post]
func RenameIndex(ctx *gin.Context) {
	var req dto.IndexRenameRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		zap.S().Debugf("mark 1")
		zap.S().Debugf("params : %+v", req)
		ctx.JSON(http.StatusOK, dto.Fail(dto.ParamsErrCode))
		return
	}
	// 检查destName是否符合要求
	if !wrench.ValidateIndexName(req.DestIndexName) {
		zap.S().Debugf("Unaccept index name %v", req.DestIndexName)
		ctx.JSON(http.StatusOK, dto.FailWithMessage(dto.IndexNameErrCode, "dest index name not allow."))
	}
	claims, exists := ctx.Get("claims")
	if !exists {
		zap.S().Infof("Unable to get the claims")
		ctx.JSON(http.StatusOK, dto.Fail(dto.UserTokenErrCode))
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
	// 首先查看是否存在被命名的知识库，然后查看是否存在新名称的知识库，如果存在就拒绝（同一用户）
	// 然后更新即可
	// 应该不需要文件系统操作
	indexModel := model.Index{}
	tx := config.DB.Model(&indexModel).Where(
		"index_name = ? AND user_id = ?",
		req.SourceIndexName, currentUserId,
	).Update("index_name = ?", req.SourceIndexName)
	err = tx.First(&indexModel).Error
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			zap.S().Errorf(
				"User with user_id %v source_index_name %v not found",
				currentUserId, req.SourceIndexName)
			ctx.JSON(
				http.StatusOK,
				dto.FailWithMessage(dto.IndexNotExistErrCode, "try to rename the index but record note found."),
			)
		} else {
			zap.S().Errorf(
				"Unexcept Error: User with user_id %v source_index_name %v , ",
				currentUserId, req.SourceIndexName)
			ctx.JSON(http.StatusOK, dto.Fail(dto.InternalErrCode))
		}
	}
	ctx.JSON(http.StatusOK, dto.Success())

}
