// -------------------------------------------------
// Package dto
// Author: hanzhi
// Date: 2024/12/8
// -------------------------------------------------

package dto

type IndexCRUDRequest struct {
	IndexName string `json:"index_name" binding:"required"`
}

type IndexRenameRequest struct {
	SourceIndexName string `json:"source_index_name" binding:"required"`
	DestIndexName   string `json:"dest_index_name" binding:"required"`
}
