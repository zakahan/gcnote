// -------------------------------------------------
// Package model
// Author: hanzhi
// Date: 2024/12/10
// -------------------------------------------------

package model

import "gorm.io/gorm"

type Share struct {
	gorm.Model
	ShareId        string
	IndexId        string
	KBFileId       string
	ShareToUserId  string
	PermissionType string
}
