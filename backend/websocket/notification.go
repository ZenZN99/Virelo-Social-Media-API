package websocket

import (
	"github.com/gorilla/websocket"
)

type NotificationGateway struct {
	Hub         *Hub
	Connections map[string]*websocket.Conn // socketId -> conn
}

func (n *NotificationGateway) SendNotification(receiverId string, payload interface{}) {

	socketIDs := n.Hub.GetSockets(receiverId)

	if len(socketIDs) == 0 {
		return
	}

	for _, socketID := range socketIDs {

		conn, ok := n.Connections[socketID]
		if !ok {
			continue
		}

		_ = conn.WriteJSON(map[string]interface{}{
			"event": "receive-notification",
			"data":  payload,
		})
	}
}
