// -------------------------------------------------
// Package dto
// Author: hanzhi
// Date: 2024/12/27
// -------------------------------------------------

package dto

type ImageRequest struct {
	IndexId  string `params:"index_id" binding:"required"`
	KBFileId string `params:"kb_file_id" binding:"required"`
}
