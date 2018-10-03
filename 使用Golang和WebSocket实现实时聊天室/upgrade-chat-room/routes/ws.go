package routes

import (
	"../routes/chatroom"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WebSocket(server *gin.Engine)  {

	r := server.Group("/websocket")

	r.GET("/room", func(c *gin.Context) {
		user := c.Query("user")
		c.HTML(http.StatusOK,"websocket.html", struct {
			User string
		}{user})
	})


	r.GET("/room/socket", func(c *gin.Context) {
		user := c.Query("user")

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			panic(err)
		}

		chatroom.Join(user)
		defer chatroom.Leave(user)

		// Join the room.
		subscription := chatroom.Subscribe()
		defer subscription.Cancel()

		//先把历史消息推送出去
		// Send down the archive.
		for _, event := range subscription.Archive {
			if conn.WriteJSON(&event) != nil {
				// They disconnected
				return
			}
		}

		// In order to select between websocket messages and subscription events, we
		// need to stuff websocket events into a channel.
		newMessages := make(chan string)
		go func() {
			var res = struct {
				Msg string `json:"msg"`
			}{}
			for {
				err := conn.ReadJSON(&res)
				if err != nil {
					close(newMessages)
					return
				}
				newMessages <- res.Msg
			}
		}()

		// Now listen for new events from either the websocket or the chatroom.
		for {
			select {
			case event := <-subscription.NewMsg:
				if conn.WriteJSON(&event) != nil {
					// They disconnected.
					return
				}
			case msg, ok := <-newMessages:
				// If the channel is closed, they disconnected.
				if !ok {
					return
				}
				// Otherwise, say something.
				chatroom.Say(user, msg)
			}
		}
	})
}




