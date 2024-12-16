// -------------------------------------------------
// Package dto
// Author: hanzhi
// Date: 2024/12/16
// -------------------------------------------------

package dto

type RecycleRequest struct {
	KBFileId string `json:"kb_file_id" binding:"required"`
	IndexId  string `json:"index_id" binding:"required"`
}
