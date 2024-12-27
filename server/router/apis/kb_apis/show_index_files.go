// -------------------------------------------------
// Package kb_apis
// Author: hanzhi
// Date: 2024/12/15
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

// ShowIndexFiles
// @Summary 展示同一个用户下的所有知识库
// @Description 展示知识库(这个是知识库的摘要展示，如果要看更细节的比如每个表有哪些文件，需要另外的接口)
// @ID			show-index_files
// @Tags		index
// @Accept		json
// @Produce		json
// @Param		request		body		dto.KBFileShowRequest true "登录请求体"
// @Success 	200			{object} 	dto.BaseResponse	"成功响应，返回success"
// @Failure		400			{object}	dto.BaseResponse	"参数错误(code:40000)"
// @Failure		401			{object}	dto.BaseResponse	"Token错误(code:40101)"
// @Failure		500			{object}	dto.BaseResponse	"服务器内部错误(code:50000)"
// @Router 		/index/show_files [get]
func ShowIndexFiles(ctx *gin.Context) {
	var req dto.KBFileShowRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		// 打印请求体
		zap.S().Debugf("Invalid parameters: %+v", ctx.Request)
		ctx.JSON(http.StatusBadRequest, dto.Fail(dto.ParamsErrCode))
		//zap.S().Debugf("%v",ctx.)
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

	// 获取文件
	kbFileModel := model.KBFile{}
	tx := config.DB.Model(&kbFileModel).Where(
		"index_id = ?",
		req.IndexId,
	)
	// 输出
	var kbFileList []model.KBFile
	err = tx.Find(&kbFileList).Error
	if err != nil {
		zap.S().Errorf("failed to find kb_file list: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}
	ctx.JSON(http.StatusOK, dto.SuccessWithData(kbFileList))

}
