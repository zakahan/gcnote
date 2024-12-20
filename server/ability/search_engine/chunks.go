// -------------------------------------------------
// Package search_engine
// Author: hanzhi
// Date: 2024/12/19
// -------------------------------------------------

package search_engine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gcnote/server/ability/document"
	"github.com/elastic/go-elasticsearch/v8"
	"strings"
)

func AddDocuments(client *elasticsearch.Client, indexName string, documents []*document.Document, embedding [][]float64) error {
	//client := config.ElasticClient
	for i, doc := range documents {
		// todo 这里应该将document转为一个map，然后再map里添加vector - embedding
		docMap := doc.ToMap()
		docMap["vector"] = embedding[i]

		// 首先将documents转为
		jsonData, err := json.Marshal(docMap)
		if err != nil {
			return err
		}
		index, err := client.Index(
			indexName,
			bytes.NewReader(jsonData),
			client.Index.WithContext(context.Background()),
			//client.Index.WithDocumentID(doc.Metadata["doc_id"]),
		)
		if err != nil {
			return err
		}
		fmt.Println(i, index)
	}
	return nil
}

func DeleteByTerm(client *elasticsearch.Client, indexName string, key string, value string) error {
	// 构建查询体
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"metadata." + key: value,
			},
		},
	}

	// 将查询体序列化为 JSON 字符串
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return err
	}

	// 使用 DeleteByQuery API 删除文档
	res, err := client.DeleteByQuery(
		[]string{indexName},                  // 索引列表
		strings.NewReader(string(queryJSON)), // 查询体
		client.DeleteByQuery.WithContext(context.Background()),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// 检查响应状态码
	if res.IsError() {
		return fmt.Errorf("failed to delete documents: %s", res.String())
	}

	return nil
}
