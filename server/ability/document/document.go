// -------------------------------------------------
// Package search_engine
// Author: hanzhi
// Date: 2024/12/18
// -------------------------------------------------

package document

import (
	"fmt"
	"github.com/google/uuid"
)

type Document struct {
	PageContent string            `json:"page_content"`
	Metadata    map[string]string `json:"metadata"`
}

func (d *Document) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"page_content": d.PageContent,
		"metadata":     d.Metadata,
	}
}

func ConvertDocument(doc map[string]interface{}) (*Document, error) {
	// 验证并转换 pageContent
	pageContent, ok := doc["page_content"].(string)
	if !ok {
		return nil, fmt.Errorf("page_content is not a string")
	}

	// 验证并转换 metadata
	metadata, ok := doc["metadata"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("metadata is not a map")
	}

	// 创建新的 map[string]string 来存储字符串类型的 metadata
	metadataStrMap := make(map[string]string)
	for key, value := range metadata {
		// 确保每个值都是字符串类型
		if strValue, ok := value.(string); ok {
			metadataStrMap[key] = strValue
		} else {
			return nil, fmt.Errorf("metadata value for key '%s' is not a string", key)
		}
	}

	// 创建 Document 实例并赋值
	document := Document{
		PageContent: pageContent,
		Metadata:    metadataStrMap,
	}

	return &document, nil
}

type DocType int

const (
	TEXT DocType = iota
	IMAGE
	TABLE
)

func (docType DocType) String() string {
	return [...]string{"TEXT", "IMAGE", "TABLE"}[docType]
}

func ShowDocumentExample() *Document {
	var doc Document
	doc = Document{
		PageContent: "这是文本内容。",
		Metadata: map[string]string{
			"doc_id":     uuid.New().String(),
			"kb_file_id": uuid.New().String(),
			"index_id":   uuid.New().String(),
			"type":       TEXT.String(),
			"image_path": "",
		},
	}
	fmt.Println(doc)
	return &doc
}
