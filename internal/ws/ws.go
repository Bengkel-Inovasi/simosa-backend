package ws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"website/backend/internal/models"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

type Client struct {
	conn *websocket.Conn
}

var (
	clients   = make(map[*Client]bool)
	clientsMu sync.Mutex
)

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Error upgrading to WebSocket: %v\n", err)
		return
	}

	client := &Client{conn: conn}
	clientsMu.Lock()
	clients[client] = true
	clientsMu.Unlock()

	defer func() {
		clientsMu.Lock()
		delete(clients, client)
		clientsMu.Unlock()
		conn.Close()
	}()

	// Keep connection alive
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func BroadcastReading(reading models.SensorReading) {
	message, _ := json.Marshal(map[string]interface{}{
		"type": "reading",
		"data": reading,
	})
	broadcast(message)
}

func BroadcastNotification(text string) {
	message, _ := json.Marshal(map[string]interface{}{
		"type": "notification",
		"data": text,
	})
	broadcast(message)
}

func broadcast(message []byte) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for client := range clients {
		err := client.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			fmt.Printf("Error writing to WebSocket: %v\n", err)
			client.conn.Close()
			delete(clients, client)
		}
	}
}
