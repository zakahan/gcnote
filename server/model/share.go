// -------------------------------------------------
// Package model
// Author: hanzhi
// Date: 2024/1/4
// -------------------------------------------------

package model

import "gorm.io/gorm"

// ShareFile 分享文件表
type ShareFile struct {
	gorm.Model
	ShareFileId string // 分享文件ID，与原KBFile的ID保持一致
	IndexId     string // 来源知识库ID
	KBFileId    string // 来源知识库文件ID
	FileName    string // 文件名
	UserId      string // 分享文件的所有者
	Password    string // 文件密码，可以选择不设置，那样的话就是所有人都能用了。看的时候检查这个存在与否，不存在直接读
}
