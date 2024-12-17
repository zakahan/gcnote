// -------------------------------------------------
// Package model
// Author: hanzhi
// Date: 2024/12/16
// -------------------------------------------------

package model

import "gorm.io/gorm"

type Recycle struct {
	gorm.Model
	UserId        string
	SourceIndexId string
	KBFileId      string
	KBFileName    string
}
