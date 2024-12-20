// -------------------------------------------------
// Package search_engine
// Author: hanzhi
// Date: 2024/12/19
// -------------------------------------------------

package search_engine

import (
	"gcnote/server/ability/document"
	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/zap"
)

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
