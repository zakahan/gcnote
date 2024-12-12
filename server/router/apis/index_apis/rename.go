// -------------------------------------------------
// Package index_apis
// Author: hanzhi
// Date: 2024/12/11
// -------------------------------------------------

package index_apis

import (
	"gcnote/server/dto"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"net/http"
)

// RenameIndex
// @Summary 知识库重命名
// @Description 展示知识库(这个是知识库的摘要展示，如果要看更细节的比如每个表有哪些文件，需要另外的接口)
// @ID			rename-index
// @Tags		index
// @Accept		json
// @Produce		json
// @Param		request		body		dto.IndexCRUDRequset true "登录请求体"
// @Success 	200			{object} 	dto.BaseResponse	"成功响应，返回success"
// @Failure		200			{object}	dto.BaseResponse	"参数错误(code:40000)"
// @Failure		200			{object}	dto.BaseResponse	"Token错误(code:40101)"
// @Failure		200			{object}	dto.BaseResponse	"服务器内部错误(code:50000)"
// @Router 		/index/rename [get]
func RenameIndex(ctx *gin.Context) {
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
}
