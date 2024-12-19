// -------------------------------------------------
// Package search_engine
// Author: hanzhi
// Date: 2024/12/19
// -------------------------------------------------

package search_engine

import (
	"encoding/json"
	"fmt"
	"gcnote/server/ability/document"
	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/zap"
	"io"
	"strings"
)

func IndexDelete(client *elasticsearch.Client, indexName string) {
	//client := config.ElasticClient
	response, err := client.Indices.Delete([]string{indexName})
	if err != nil {
		fmt.Printf("删除出现错误 %v", err)
		zap.S().Errorf("删除出现错误 %v", err)
	}
	fmt.Println(response)
}

func ShowAllTexts(client *elasticsearch.Client, indexName string) ([]*document.Document, error) {
	query := `{ "query": { "match_all": {} } }`
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
