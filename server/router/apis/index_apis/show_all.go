// -------------------------------------------------
// Package index_apis
// Author: hanzhi
// Date: 2024/12/11
// -------------------------------------------------

package index_apis

import (
	"gcnote/server/config"
	"gcnote/server/dto"
	"gcnote/server/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"net/http"
)

// ShowIndexes
// @Summary 展示同一个用户下的所有知识库
// @Description 展示知识库(这个是知识库的摘要展示，如果要看更细节的比如每个表有哪些文件，需要另外的接口)
// @ID			show-all-index
// @Tags		index
// @Accept		json
// @Produce		json
// @Param		request		body		dto.IndexCRUDRequest true "登录请求体"
// @Success 	200			{object} 	dto.BaseResponse	"成功响应，返回success"
// @Failure		400			{object}	dto.BaseResponse	"参数错误(code:40000)"
// @Failure		401			{object}	dto.BaseResponse	"Token错误(code:40101)"
// @Failure		500			{object}	dto.BaseResponse	"服务器内部错误(code:50000)"
// @Router 		/index/show_indexes [get]
func ShowIndexes(ctx *gin.Context) {
	claims, exists := ctx.Get("claims")
	if !exists {
		zap.S().Infof("Unable to get the claims")
		ctx.JSON(http.StatusUnauthorized, dto.Fail(dto.UserTokenErrCode))
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
	// 查看当前user_id 的所有表内信息，
	indexModel := model.Index{}
	tx := config.DB.Model(&indexModel).Where(
		"user_id = ?",
		currentUserId,
	)
	var indexList []model.Index
	err := tx.Find(&indexList).Error
	if err != nil {
		zap.S().Errorf("failed to find index: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}
	ctx.JSON(http.StatusOK, dto.SuccessWithData(indexList))

}
