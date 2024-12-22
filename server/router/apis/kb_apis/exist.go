// -------------------------------------------------
// Package kb_apis
// Author: hanzhi
// Date: 2024/12/22
// -------------------------------------------------

package kb_apis

import (
	"errors"
	"gcnote/server/config"
	"gcnote/server/dto"
	"gcnote/server/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
)

// KBFileExist
// @Summary 检查某个文件是否存在
// @Description 检查某个文件是否存在
// @ID check-kb-file
// @Tags index
// @Accept json
// @Produce json
// @Param indexId query string true "知识库ID"
// @Success 200 {object} dto.BaseResponse "成功响应，返回success"
// @Failure 400 {object} dto.BaseResponse "参数错误(code:40000)"
// @Failure 404 {object} dto.BaseResponse "知识库不存在(code:40201)"
// @Failure 500 {object} dto.BaseResponse "服务器内部错误(code:50000)"
// @Router /index/kb_file_exist [get]
func KBFileExist(ctx *gin.Context) {
	var req dto.KBFileRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		zap.S().Debugf("mark 1")
		zap.S().Debugf("params : %+v", req)
		ctx.JSON(http.StatusBadRequest, dto.Fail(dto.ParamsErrCode))
		return
	}
	claims, _ := ctx.Get("claims")

	zap.S().Debugf("claims: %v", claims)
	currentUser := claims.(jwt.MapClaims)
	currentUserId := currentUser["sub"].(string)
	if currentUserId == "" {
		zap.S().Debugf("currentUserId is empty")
		ctx.JSON(http.StatusBadRequest, dto.Fail(dto.ParamsErrCode))
		return
	}

	var index model.Index
	if err := config.DB.Where("index_id = ?", req.KBFileId).First(&index).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, dto.Fail(dto.KBFileNotExistErrCode))
		} else {
			ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		}
		return
	}

	ctx.JSON(http.StatusOK, dto.Success())
}
