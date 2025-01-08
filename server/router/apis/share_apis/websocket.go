// -------------------------------------------------
// Package share_apis
// Author: hanzhi
// Date: 2025/1/1
// -------------------------------------------------

package share_apis

import (
	"encoding/json"
	"gcnote/server/dto"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// WebSocketManager manages all connected clients
type WebSocketManager struct {
	clients    map[*websocket.Conn]bool
	operations []dto.Operation
	mu         sync.Mutex
}

var (
	// Global operation queue
	operationQueue []dto.Operation
	// WebSocket upgrader
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins in development, should be restricted in production
		},
	}
	// Global manager instance
	globalManager = newWebSocketManager()
)

// newWebSocketManager creates a new WebSocket manager
func newWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		clients:    make(map[*websocket.Conn]bool),
		operations: make([]dto.Operation, 0),
	}
}

// addConnection adds a WebSocket connection to the manager
func (manager *WebSocketManager) addConnection(conn *websocket.Conn) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	manager.clients[conn] = true
}

// removeConnection removes a WebSocket connection from the manager
func (manager *WebSocketManager) removeConnection(conn *websocket.Conn) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	delete(manager.clients, conn)
}

// broadcastOperation broadcasts an operation to all clients
func (manager *WebSocketManager) broadcastOperation(op dto.Operation) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	// Add operation to queue
	operationQueue = append(operationQueue, op)

	// Marshal operation to JSON
	opJson, err := json.Marshal(op)
	if err != nil {
		zap.S().Errorf("Error marshalling operation: %v", err)
		return
	}

	// Broadcast to all clients
	for conn := range manager.clients {
		err := conn.WriteMessage(websocket.TextMessage, opJson)
		if err != nil {
			zap.S().Errorf("Error sending message to %v: %v", conn.RemoteAddr(), err)
			conn.Close()
			delete(manager.clients, conn)
		}
	}
}

// handleWebSocketConnection handles individual WebSocket connections
func handleWebSocketConnection(manager *WebSocketManager, conn *websocket.Conn, documentId, clientId string) {

	defer conn.Close()

	manager.addConnection(conn)
	zap.S().Infof("New WebSocket connection: %v", conn.RemoteAddr())
	zap.S().Debugf("建立连接documentId %v", documentId)
	zap.S().Debugf("建立连接clientId %v", clientId)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			zap.S().Errorf("Error reading message: %v", err)
			manager.removeConnection(conn)
			return
		}

		var op dto.Operation
		if err := json.Unmarshal(msg, &op); err != nil {
			zap.S().Errorf("Error unmarshaling operation: %v", err)
			continue
		}

		operationQueue = append(operationQueue, op)
		manager.broadcastOperation(op)
		manager.applyOperations()
	}
}

// applyOperations processes the operation queue
func (manager *WebSocketManager) applyOperations() {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	for _, op := range operationQueue {
		zap.S().Infof("Applying operation: %+v", op)
		// Here you can implement operation transformation logic if needed
	}
	operationQueue = []dto.Operation{}
}

// HandleWebSocket
// @Summary WebSocket connection endpoint
// @Description Establishes a WebSocket connection for real-time collaboration
// @ID handle-websocket
// @Tags share
// @Accept json
// @Produce json
// @Success 101 {string} string "Switching Protocols to WebSocket"
// @Failure 400 {object} dto.BaseResponse "Bad Request"
// @Router /share/ws [get]
func HandleWebSocket(c *gin.Context) {
	documentId := c.Param("documentId")
	clientId := c.Query("clientId")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		zap.S().Errorf("Failed to upgrade connection: %v", err)
		c.JSON(http.StatusBadRequest, dto.Fail(dto.InternalErrCode))
		return
	}

	go handleWebSocketConnection(globalManager, conn, documentId, clientId)
}
