// -------------------------------------------------
// Package share_apis
// Author: hanzhi
// Date: 2025/1/1
// -------------------------------------------------

package share_apis

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	messageSync           = 0
	messageAwareness      = 1
	messageAuth           = 2
	messageQueryAwareness = 3
)

// Room represents a collaborative editing room
type Room struct {
	ID        string
	Clients   map[*Client]bool
	State     []byte // Document state
	mu        sync.RWMutex
	Broadcast chan []byte
}

// Client represents a connected client
type Client struct {
	ID   string
	Room *Room
	Conn *websocket.Conn
	Send chan []byte
}

var (
	// Configure the upgrader
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for now
		},
	}

	// Global rooms management
	rooms      = make(map[string]*Room)
	roomsMutex sync.RWMutex
)

// readVarUint reads a variable-length unsigned integer from a byte slice
func readVarUint(data []byte) (uint64, int, error) {
	var x uint64
	var s uint
	for i, b := range data {
		if b < 0x80 {
			if i > 9 || i == 9 && b > 1 {
				return 0, 0, fmt.Errorf("overflow")
			}
			return x | uint64(b)<<s, i + 1, nil
		}
		x |= uint64(b&0x7f) << s
		s += 7
	}
	return 0, 0, fmt.Errorf("buffer too small")
}

// writeVarUint writes a variable-length unsigned integer to a buffer
func writeVarUint(buf *bytes.Buffer, x uint64) {
	for x >= 0x80 {
		buf.WriteByte(byte(x) | 0x80)
		x >>= 7
	}
	buf.WriteByte(byte(x))
}

// createRoom creates a new room if it doesn't exist
func createRoom(roomID string) *Room {
	roomsMutex.Lock()
	defer roomsMutex.Unlock()

	if room, exists := rooms[roomID]; exists {
		return room
	}

	initialContent, _, err := readInitialContent(roomID)
	if err != nil {
		return nil
	}
	initialState := initializeYDoc(initialContent)

	room := &Room{
		ID:        roomID,
		Clients:   make(map[*Client]bool),
		State:     initialState,
		Broadcast: make(chan []byte),
	}
	rooms[roomID] = room

	// Start room's broadcast handler
	go room.run()

	return room
}

// run handles broadcasting messages to all clients in the room
func (r *Room) run() {
	for {
		message := <-r.Broadcast
		r.mu.RLock()
		for client := range r.Clients {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(r.Clients, client)
			}
		}
		r.mu.RUnlock()
	}
}

// readPump pumps messages from the websocket connection to the room
func (c *Client) readPump() {
	defer func() {
		c.Room.mu.Lock()
		delete(c.Room.Clients, c)
		c.Room.mu.Unlock()
		err := c.Conn.Close()
		if err != nil {
			return
		}
	}()

	for {
		messageType, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				zap.S().Errorf("error: %v", err)
			}
			break
		}

		if messageType == websocket.BinaryMessage {
			// Read message type (varint)
			msgType, _, err := readVarUint(message)
			if err != nil {
				zap.S().Errorf("error reading message type: %v", err)
				continue
			}

			// Handle different message types
			switch msgType {
			case messageSync:
				// Store sync message in room state
				c.Room.mu.Lock()
				c.Room.State = message
				zap.S().Debugf("看一看message: %v", message)
				c.Room.mu.Unlock()
			case messageAwareness:
				// Just broadcast awareness updates
			case messageAuth:
				// Handle auth if needed
			case messageQueryAwareness:
				// Handle awareness query
			}

			// Broadcast the message to all clients in the room
			c.Room.Broadcast <- message
		}
	}
}

// writePump pumps messages from the room to the websocket connection
func (c *Client) writePump() {
	defer func() {
		err := c.Conn.Close()
		if err != nil {
			return
		}
	}()

	for {
		message, ok := <-c.Send
		if !ok {
			err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			if err != nil {
				return
			}
			return
		}

		if err := c.Conn.WriteMessage(websocket.BinaryMessage, message); err != nil {
			return
		}
	}
}

// HandleWebSocket handles websocket requests from clients
func HandleWebSocket(c *gin.Context) {
	roomID := c.Param("room")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room ID is required"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		zap.S().Errorf("Failed to upgrade connection: %v", err)
		return
	}

	room := createRoom(roomID)
	client := &Client{
		ID:   fmt.Sprintf("%p", conn),
		Room: room,
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	room.mu.Lock()
	room.Clients[client] = true
	room.mu.Unlock()

	// Send current document state to new client
	if len(room.State) > 0 {
		client.Send <- room.State
	}

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
