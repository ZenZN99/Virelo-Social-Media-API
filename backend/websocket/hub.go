package websocket

import (
	"sync"
)

type OnlineUser struct {
	UserID   string
	SocketID string
}

type Hub struct {
	mu    sync.Mutex
	users map[string][]string // userId -> socketIds
}

func NewHub() *Hub {
	return &Hub{
		users: make(map[string][]string),
	}
}

// Add user socket
func (h *Hub) AddUser(userID, socketID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.users[userID] = append(h.users[userID], socketID)
}

// Remove socket
func (h *Hub) RemoveSocket(socketID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for userID, sockets := range h.users {
		newSockets := []string{}
		for _, s := range sockets {
			if s != socketID {
				newSockets = append(newSockets, s)
			}
		}
		if len(newSockets) == 0 {
			delete(h.users, userID)
		} else {
			h.users[userID] = newSockets
		}
	}
}

// Get sockets
func (h *Hub) GetSockets(userID string) []string {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.users[userID]
}