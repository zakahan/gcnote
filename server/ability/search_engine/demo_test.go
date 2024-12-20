// -------------------------------------------------
// Package search_engine
// Author: hanzhi
// Date: 2024/12/18
// -------------------------------------------------

package search_engine

import (
	"fmt"
	"gcnote/server/ability/document"
	"gcnote/server/ability/embeds"
	"github.com/elastic/go-elasticsearch/v8"
	"testing"
)

func TestShowDocumentExample(t *testing.T) {
	document.ShowDocumentExample()
}

func TestSearchAll(t *testing.T) {
	// 更新后不支持了，因为循环调用
	indexName := "gcnote"
	esCfg := elasticsearch.Config{
		Addresses: []string{"http://10.1.69.142:8581"},
		Username:  "elastic",
		Password:  "ebyte_zxcvbnm",
	}

	elasticClient, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return
	}
	_, err = ShowAllTexts(elasticClient, indexName)
	if err != nil {
		return
	}

}

func TestAdd(t *testing.T) {
	// 更新后不支持了，因为循环调用
	indexName := "gcnote"
	esCfg := elasticsearch.Config{
		Addresses: []string{"http://10.1.69.142:8581"},
		Username:  "elastic",
		Password:  "ebyte_zxcvbnm",
	}

	elasticClient, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return
	}
	docList := []*document.Document{document.ShowDocumentExample()}
	embedList, err := embeds.RandEmbedding(docList)
	if err != nil {
		return
	}
	err = AddDocuments(elasticClient, indexName, docList, embedList)
	if err != nil {
		return
	}
}

func TestSearch(t *testing.T) {
	// 更新后不支持了，因为循环调用
	indexName := "gcnote"
	esCfg := elasticsearch.Config{
		Addresses: []string{"http://10.1.69.142:8581"},
		Username:  "elastic",
		Password:  "ebyte_zxcvbnm",
	}

	elasticClient, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return
	}
	response, err := FullTextSearch(elasticClient, indexName, "文本", 1)
	if err != nil {
		return
	}
	fmt.Println(response)

}
