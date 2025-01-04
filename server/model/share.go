// -------------------------------------------------
// Package model
// Author: hanzhi
// Date: 2024/12/10
// -------------------------------------------------

package model

import "gorm.io/gorm"

// Share  分享记录（都分享给过谁）
type Share struct {
	gorm.Model
	ShareId        string
	UserId         string
	IndexId        string
	KBFileId       string
	ShareToUserId  string
	PermissionType string
}

// ShareFile 分享表
type ShareFile struct {
	gorm.Model
	OwnerUserId string
	IndexId     string
	KBFileId    string
}
