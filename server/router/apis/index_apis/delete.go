// -------------------------------------------------
// Package index_apis
// Author: hanzhi
// Date: 2024/12/11
// -------------------------------------------------

package index_apis

import (
	"gcnote/server/cache"
	"gcnote/server/config"
	"gcnote/server/dto"
	"gcnote/server/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"net/http"
)

// DeleteIndex
// goland别他妈动老子注释！
// redis代码先空着？肯定不行啊
// 不对，index_name和index_id的决断，这个地方我得好好想想，有点不对。
// 可以拼接
func DeleteIndex(ctx *gin.Context) {
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
	// helllllll
	// 查看是否有这个index，然后删除
	indexInfo := model.Index{}
	// 删除
	tx := config.DB.Model(&indexInfo).Where(
		"index_name = ? AND user_id = ?",
		req.IndexName, currentUserId,
	).Delete(&indexInfo)
	if tx.Error != nil {
		zap.S().Errorf("Delete index id :%v err:%v", req.IndexName, tx.Error)
		ctx.JSON(http.StatusOK, dto.Fail(dto.InternalErrCode))
		return
	}
	err = cache.DelIndexInfo(ctx.Request.Context(), req.IndexName)
	if err != nil {
		zap.S().Errorf("Delete.DelIndexInfo  userId:%v err:%v", currentUserId, tx.Error)
		ctx.JSON(http.StatusOK, dto.Fail(dto.InternalErrCode))
		return
	}
	ctx.JSON(http.StatusOK, dto.Success())
}
