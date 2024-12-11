// -------------------------------------------------
// Package dto
// Author: hanzhi
// Date: 2024/12/8
// -------------------------------------------------

package dto

type IndexCRUDRequset struct {
	IndexName string `json:"index_name" binding:"required"`
}
