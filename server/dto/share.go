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
