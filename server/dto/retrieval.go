// -------------------------------------------------
// Package dto
// Author: hanzhi
// Date: 2024/12/20
// -------------------------------------------------

package dto

type RetrievalRequest struct {
	IndexId         string `json:"index_id" binding:"required"`
	Query           string `json:"query" binding:"required"`
	RetrievalMethod string `json:"retrieval_method" binding:"required"`
	TopK            int    `json:"top_k"`
}
