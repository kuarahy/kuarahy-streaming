// internal/obs/manager.go
package obs

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type OBSManager struct {
	upgrader    websocket.Upgrader
	clients     map[*websocket.Conn]bool
	clientsLock sync.Mutex
	config      Config
}

func NewManager(cfg Config) *OBSManager {
	return &OBSManager{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		clients: make(map[*websocket.Conn]bool),
		config:  cfg,
	}
}

func (m *OBSManager) StartWSServer() error {
	http.HandleFunc(m.config.WSEndpoint, m.handleWebSocket)
	return http.ListenAndServe(m.config.WSAddress, nil)
}

func (m *OBSManager) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := m.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	m.clientsLock.Lock()
	m.clients[conn] = true
	m.clientsLock.Unlock()

	defer func() {
		m.clientsLock.Lock()
		delete(m.clients, conn)
		m.clientsLock.Unlock()
		conn.Close()
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (m *OBSManager) SendEvent(event Event) {
	msg, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshaling event: %v", err)
		return
	}

	m.clientsLock.Lock()
	defer m.clientsLock.Unlock()

	for client := range m.clients {
		if err := client.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Printf("Error sending to client: %v", err)
			client.Close()
			delete(m.clients, client)
		}
	}
}
