// -------------------------------------------------
// Package search_engine
// Author: hanzhi
// Date: 2024/12/18
// -------------------------------------------------

package search_engine

import (
	"fmt"
	"github.com/google/uuid"
	"strconv"
)

type Document struct {
	PageContent string
	Metadata    map[string]string
}

type DocType int

const (
	TEXT = iota
	IMAGE
	TABLE
)

func (docType DocType) String() string {
	return [...]string{"TEXT", "IMAGE", "TABLE"}[docType]
}

func showDocumentExample() {
	var doc Document
	doc = Document{
		PageContent: "这里是样例文本",
		Metadata: map[string]string{
			"doc_id":     uuid.New().String(),
			"kb_file_id": uuid.New().String(),
			"index_id":   uuid.New().String(),
			"type":       strconv.Itoa(TEXT),
			"image_path": "",
		},
	}
	fmt.Println(doc)
}
