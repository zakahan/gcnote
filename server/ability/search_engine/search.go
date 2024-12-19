// -------------------------------------------------
// Package search_engine
// Author: hanzhi
// Date: 2024/12/19
// -------------------------------------------------

package search_engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gcnote/server/ability/document"
	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/zap"
	"io"
	"log"
	"strings"
)

func BaseSearch(query string, client *elasticsearch.Client, indexName string) ([]*document.Document, error) {
	response, err := client.Search(
		client.Search.WithIndex(indexName),
		client.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			zap.S().Error("close error %v", err)
		}
	}(response.Body)

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, err
	}

	// 获取 hits 字段
	hitsMap, ok := result["hits"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected 'hits' to be a map")
	}

	// 获取 hits 字段中的 hits 切片
	hitsSlice, ok := hitsMap["hits"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected 'hits.hits' to be a slice")
	}

	docSlice := make([]*document.Document, len(hitsSlice))
	for i := range hitsSlice {
		//score := hitsSlice[i].(map[string]interface{})["_score"].(float64)
		doc := hitsSlice[i].(map[string]interface{})["_source"].(map[string]interface{})
		docP, err := document.ConvertDocument(doc)
		if err != nil {
			return nil, err
		}
		fmt.Println(docP)
		docSlice[i] = docP
	}

	return docSlice, nil
}

func FullTextSearch(client *elasticsearch.Client, indexName string, userQuery string, topK int) ([]*document.Document, error) {
	//query := `{"query": {"match": {}}}`
	query := fmt.Sprintf(
		`{"query": {"match": { "page_content": "%s" } },"size": %v }`,
		userQuery, topK)
	return BaseSearch(query, client, indexName)
}

func VectorSearch(client *elasticsearch.Client, indexName string, queryVector []float64, topK int) ([]*document.Document, error) {
	// 定义 KNN 查询的请求体
	query := map[string]interface{}{
		"knn": map[string]interface{}{
			"field":          "vector",    // dense_vector 字段名
			"query_vector":   queryVector, // 查询向量
			"k":              topK,        // 最近邻数量
			"num_candidates": 100,         // 候选数目，用于提高性能
		},
	}
	// 将查询体转换为 JSON
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}
	queryString := buf.String()
	return BaseSearch(queryString, client, indexName)
}

func KeywordsSearch(client *elasticsearch.Client, indexName string, word string) ([]*document.Document, error) {
	query := fmt.Sprintf(
		`"query":{
                    "match_phrase": {
                        "page_content": "%v"
                    }
                }`, word)
	return BaseSearch(query, client, indexName)
}
