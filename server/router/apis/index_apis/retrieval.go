// -------------------------------------------------
// Package index_apis
// Author: hanzhi
// Date: 2024/12/20
// -------------------------------------------------

package index_apis

import (
	"gcnote/server/ability/document"
	"gcnote/server/ability/embeds"
	"gcnote/server/ability/search_engine"
	"gcnote/server/config"
	"gcnote/server/dto"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"net/http"
)

// RetrievalIndex
// @Summary 检索文件
// @Description 根据文档切片进行索引
// @ID         retrieval_index
// @Tags       index
// @Accept     json
// @Produce    json
// @Param      request  body      dto.RetrievalRequest true "搜索请求体"
// @Success    200      {object}  dto.BaseResponse "成功响应，返回搜索结果"
// @Failure    400      {object}  dto.BaseResponse "参数错误(code:40000)"
// @Failure    401      {object}  dto.BaseResponse "Token错误(code:40101)"
// @Failure    500      {object}  dto.BaseResponse "服务器内部错误(code:50000)"
// @Router     /index/retrieval [post]
func RetrievalIndex(ctx *gin.Context) {
	var req dto.RetrievalRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Fail(dto.ParamsErrCode))
		return
	}
	allowMethod := []string{"vector_search", "full_text_search", "keyword_search", "multi_keyword_search"}
	flag := false
	for _, value := range allowMethod {
		if req.RetrievalMethod == value {
			flag = true
			break
		}
	}
	if !flag {
		ctx.JSON(http.StatusBadRequest, dto.FailWithMessage(dto.ParamsErrCode, "retrieval method not allown"))
		return
	}
	if req.TopK == 0 {
		req.TopK = 10 // 默认为10
	}

	// 获取用户信息
	claims, exists := ctx.Get("claims")
	if !exists {
		zap.S().Infof("Unable to get the claims")
		ctx.JSON(http.StatusUnauthorized, dto.Fail(dto.UserTokenErrCode))
		return
	}
	currentUser := claims.(jwt.MapClaims)
	currentUserId := currentUser["sub"].(string)
	if currentUserId == "" {
		zap.S().Debugf("currentUserId is empty")
		ctx.JSON(http.StatusUnauthorized, dto.Fail(dto.ParamsErrCode))
		return
	}

	// 调用es进行搜索
	var docList []*document.Document
	switch req.RetrievalMethod {
	case "vector_search":
		// 获取向量
		var vector []float64
		vector, err = embeds.QueryRandEmbedding(req.Query)
		if err != nil {
			zap.S().Errorf("Failed to get vector, err: %v", err)
			ctx.JSON(http.StatusBadRequest, dto.FailWithMessage(dto.ParamsErrCode, "Failed to get vector"))
		}
		docList, err = search_engine.VectorSearch(config.ElasticClient, "gcnote-"+req.IndexId, vector, req.TopK)
	case "full_text_search":
		docList, err = search_engine.FullTextSearch(config.ElasticClient, "gcnote-"+req.IndexId, req.Query, req.TopK)

	case "keyword_search":
		docList, err = search_engine.KeywordsSearch(config.ElasticClient, "gcnote-"+req.IndexId, req.Query)

	}
	if err != nil {
		zap.S().Errorf("Failed to search, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}
	var docMapList []map[string]interface{}
	for _, doc := range docList {
		docMapList = append(docMapList, doc.ToMap())
	}
	ctx.JSON(http.StatusOK, dto.SuccessWithData(docMapList))
}
