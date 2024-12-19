// -------------------------------------------------
// Package search_engine
// Author: hanzhi
// Date: 2024/12/19
// -------------------------------------------------

package search_engine

import (
	"fmt"
	"gcnote/server/ability/document"
	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/zap"
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
	return BaseSearch(query, client, indexName)
}
