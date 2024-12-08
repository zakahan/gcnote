// -------------------------------------------------
// Package model
// Author: hanzhi
// Date: 2024/12/8
// -------------------------------------------------

package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserId   string
	UserName string
	Password string
}
