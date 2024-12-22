// -------------------------------------------------
// Package dto
// Author: hanzhi
// Date: 2024/12/8
// -------------------------------------------------

package dto

import "mime/multipart"

type KBFileCreateRequest struct {
	KBFileName string `json:"kb_file_name" binding:"required"`
	IndexId    string `json:"index_id" binding:"required"`
}

type KBFileRenameRequest struct {
	IndexId        string `json:"index_id" binding:"required"`          // 知识库索引ID
	KBFileId       string `json:"kb_file_id" binding:"required"`        // 原文件ID
	KBFileName     string `json:"kb_file_name" binding:"required"`      // 原文件名
	DestKBFileName string `json:"dest_kb_file_name" binding:"required"` // 新文件名

}

type KBFileRequest struct {
	KBFileId string `json:"kb_file_id" binding:"required"`
}

// 都用form标签
type KBFileAddRequest struct {
	IndexId string `form:"index_id" binding:"required"`
	//KBFileName string                `form:"kb_file_name" binding:"required"`
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type KBFileShowRequest struct {
	IndexId string `json:"index_id" binding:"required"`
}

type KBFileUDRequest struct {
	KBFileId   string `json:"kb_file_id" binding:"required"`
	KBFileName string `json:"kb_file_name" binding:"required"`
	IndexId    string `json:"index_id" binding:"required"`
}

type KBFileSearchRequest struct {
	KBFileName    string `json:"kb_file_name" binding:"required"`     // 文件名，可选
	IndexId       string `json:"index_id" binding:"required"`         // 索引ID，必填
	IsFuzzySearch bool   `json:"is_fuzzy_search" binding:"omitempty"` // 是否模糊搜索，可选，默认false
}
