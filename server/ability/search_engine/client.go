// -------------------------------------------------
// Package es
// Author: hanzhi
// Date: 2024/12/18
// -------------------------------------------------

package search_engine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gcnote/server/ability/document"
	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/zap"
	"io"
	"log"
)

func IndexCreate(client *elasticsearch.Client, indexName string) (error, int) {
	//client := config.ElasticClient
	body := map[string]interface{}{
		"settings": _defaultSettings(),
		"mappings": _defaultMappings(),
	}
	// 将 body 转换为 JSON
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		log.Fatalf("Error encoding body: %s", err)
	}
	res, err := client.Indices.Create(
		indexName,
		client.Indices.Create.WithBody(&buf),
		client.Indices.Create.WithContext(context.Background()),
	)
	if err != nil {
		return err, 500
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			zap.S().Error("close error %v", err)
		}
	}(res.Body)

	return err, res.StatusCode
}

func IndexExist(client *elasticsearch.Client, indexName string) (error, int) {
	//client := config.ElasticClient
	// 检查索引是否存在
	res, err := client.Indices.Exists(
		[]string{indexName},
		client.Indices.Exists.WithContext(context.Background()),
	)
	if err != nil {
		return err, 500
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			zap.S().Error("close error %v", err)
		}
	}(res.Body)
	return err, res.StatusCode
}

func AddTexts(client *elasticsearch.Client, indexName string, documents []*document.Document, embedding [][]float64) error {
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
			client.Index.WithDocumentID(doc.Metadata["doc_id"]),
		)
		if err != nil {
			return err
		}
		fmt.Println(i, index)
	}
	return nil
}
