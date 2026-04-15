package websocket

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatSocket struct {
	Conn   *websocket.Conn
	UserID string
	Hub    *Hub
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewChatSocket(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		userID := c.Query("userId")

		socket := &ChatSocket{
			Conn:   conn,
			UserID: userID,
			Hub:    hub,
		}

		hub.AddUser(userID, conn.RemoteAddr().String())

		go socket.Listen()
	}
}
