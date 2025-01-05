// -------------------------------------------------
// Package dto
// Author: hanzhi
// Date: 2025/1/1
// -------------------------------------------------

package dto

// WebSocket operation types
const (
	InsertOperation = 1
	DeleteOperation = 2
)

// Operation represents a single WebSocket operation
type Operation struct {
	Type     int    `json:"op_type"`   // Operation type: insert, delete
	Position int    `json:"position"`  // Operation position
	Content  string `json:"content"`   // Operation content
	SenderId string `json:"sender_id"` // Operation sender ID
}

// WebSocketResponse represents a WebSocket response
type WebSocketResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Data    Operation `json:"data,omitempty"`
}

// CreateShareFileRequest 创建分享文件的请求
type CreateShareFileRequest struct {
	KBFileId string `json:"kb_file_id" binding:"required"`
}

// DeleteShareFileRequest 删除分享文件的请求
type DeleteShareFileRequest struct {
	ShareFileId string `json:"share_file_id" binding:"required"`
}

// ShareFileExistResponse 分享文件是否存在的响应
type ShareFileExistResponse struct {
	ShareFileId string `json:"shareFileId"`
	Exist       bool   `json:"exist"`
}

// ShareFileInfo 分享文件信息
type ShareFileInfo struct {
	ShareFileId string `json:"shareFileId"`
	IndexId     string `json:"indexId"`
	KBFileId    string `json:"kbFileId"`
	FileName    string `json:"fileName"`
	CreatedAt   string `json:"createdAt"`
	Password    string `json:"password"`
}

// ShareFileListResponse 用户分享文件列表响应
type ShareFileListResponse struct {
	Total int64           `json:"total"`
	List  []ShareFileInfo `json:"list"`
}
