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
	"gcnote/server/config"
	"go.uber.org/zap"
	"io"
	"log"
)

func IndexCreate(indexName string) (error, int) {
	es := config.ElasticClient
	body := map[string]interface{}{
		"settings": _defaultSettings(),
		"mappings": _defaultMappings(),
	}
	// 将 body 转换为 JSON
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		log.Fatalf("Error encoding body: %s", err)
	}
	res, err := es.Indices.Create(
		indexName,
		es.Indices.Create.WithBody(&buf),
		es.Indices.Create.WithContext(context.Background()),
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

func IndexExist(indexName string) (error, int) {
	es := config.ElasticClient
	// 检查索引是否存在
	res, err := es.Indices.Exists(
		[]string{indexName},
		es.Indices.Exists.WithContext(context.Background()),
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

//func AddTexts(document[]) {}
