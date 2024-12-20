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
	"gcnote/server/ability/document"
	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/zap"
	"io"
	"log"
)

func IndexCreate(client *elasticsearch.Client, indexName string) error {
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
		return err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			zap.S().Error("close error %v", err)
		}
	}(res.Body)

	return err
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

func IndexDelete(client *elasticsearch.Client, indexName string) error {
	//client := config.ElasticClient
	_, err := client.Indices.Delete([]string{indexName})
	if err != nil {
		zap.S().Errorf("删除出现错误 %v", err)
		return err
	}
	//fmt.Println(response)
	return nil
}

func ShowAllTexts(client *elasticsearch.Client, indexName string) ([]*document.Document, error) {
	query := `{ "query": { "match_all": {} } }`
	return BaseSearch(query, client, indexName)
}
