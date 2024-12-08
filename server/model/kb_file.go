// -------------------------------------------------
// Package model
// Author: hanzhi
// Date: 2024/12/8
// -------------------------------------------------

package model

import "gorm.io/gorm"

type KBFile struct {
	gorm.Model
	IndexId    string
	KBFileId   string
	KBFileName string
}
