// -------------------------------------------------
// Package dto
// Author: hanzhi
// Date: 2024/12/8
// -------------------------------------------------

package dto

import "mime/multipart"

type KBFileRequest struct {
	IndexName  string `json:"index_name" binding:"required"`
	KBFileName string `json:"kb_file_nane" binding:"required"`
}

type KBFileRenameRequest struct {
	IndexName      string `json:"index_name" binding:"required"`
	SourceFileName string `json:"source_kb_file_nane" binding:"required"`
	DestFileName   string `json:"dest_kb_file_nane" binding:"required"`
}

// 都用form标签
type KBFileAddRequest struct {
	IndexName  string                `form:"index_name" binding:"required"`
	KBFileName string                `form:"kb_file_name" binding:"required"`
	File       *multipart.FileHeader `form:"file" binding:"required"`
}
