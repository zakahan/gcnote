// -------------------------------------------------
// Package dto
// Author: hanzhi
// Date: 2024/12/8
// -------------------------------------------------

package dto

type UserRegisterRequest struct {
	UserName string `json:"user_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserLoginRequest struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserNameRequest struct {
	UserName string `json:"user_name" binding:"required"`
}

type UpdatePasswordRequest struct {
	Password string `json:"password" binding:"required"`
}
