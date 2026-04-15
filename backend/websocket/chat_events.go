package websocket

import (
	"encoding/json"
)

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func (s *ChatSocket) Listen() {
	defer s.Conn.Close()

	for {
		_, msg, err := s.Conn.ReadMessage()
		if err != nil {
			s.Hub.RemoveSocket(s.Conn.RemoteAddr().String())
			break
		}

		var message Message
		json.Unmarshal(msg, &message)

		switch message.Type {

		// typing
		case "typing":
			s.handleTyping(message.Data)

		// send message
		case "send-message":
			s.handleSendMessage(message.Data)

		// message seen
		case "message-seen":
			s.handleSeen(message.Data)
		}
	}
}

func (s *ChatSocket) handleTyping(data interface{}) {
	m := data.(map[string]interface{})

	receiverId := m["receiverId"].(string)

	sockets := s.Hub.GetSockets(receiverId)

	for _, socketID := range sockets {
		_ = socketID
	}
}

func (s *ChatSocket) handleSendMessage(data interface{}) {
	m := data.(map[string]interface{})

	receiverId := m["receiverId"].(string)

	sockets := s.Hub.GetSockets(receiverId)

	for _, socketID := range sockets {
		_ = socketID // emit receive-message
	}
}

func (s *ChatSocket) handleSeen(data interface{}) {
	m := data.(map[string]interface{})

	receiverId := m["receiverId"].(string)

	sockets := s.Hub.GetSockets(receiverId)

	for _, socketID := range sockets {
		_ = socketID // emit message-seen
	}
}
