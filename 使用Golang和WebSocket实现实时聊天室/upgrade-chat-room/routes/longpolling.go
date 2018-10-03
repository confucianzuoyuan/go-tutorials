package routes

import (
	"chat-room/routes/chatroom"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func LongPolling(server *gin.Engine)  {
	r := server.Group("/longpolling")

	// Long polling demo

	r.GET("/room", func(c *gin.Context) {
		user := c.Query("user")
		chatroom.Join(user)
		c.HTML(http.StatusOK,"longpolling.html", struct {
			User string
		}{user})
	})


	r.POST("/room/messages", func(c *gin.Context) {
		user := c.PostForm("user")
		message := c.PostForm("message")
		chatroom.Say(user, message)
	})

	r.GET("/room/leave", func(c *gin.Context) {
		user := c.Query("user")
		chatroom.Leave(user)
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	r.GET("/room/messages", func(c *gin.Context) {

		lastReceived,_ := strconv.ParseInt(c.Query("lastReceived"),10,64)
		subscription := chatroom.Subscribe()

		defer subscription.Cancel()

		// See if anything is new in the archive.
		var events []chatroom.Event
		for _, event := range subscription.Archive {
			if event.Timestamp > lastReceived {
				events = append(events, event)
			}
		}

		// If we found one, grand.
		if len(events) > 0 {
			c.JSON(http.StatusOK,events)
			return
		}
		// Else, wait for something new.
		event := <-subscription.NewMsg

		c.JSON(http.StatusOK,[]chatroom.Event{event})

		return
	})


}





