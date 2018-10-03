package routes

import (
	"../routes/chatroom"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Refresh(server *gin.Engine)  {
	r := server.Group("/refresh")

	// Refresh
	r.GET("", func(c *gin.Context) {
		user := c.Query("user")
		chatroom.Join(user)
		c.Redirect(http.StatusMovedPermanently, "/refresh/room")
	})


	r.GET("/room", func(c *gin.Context) {
		user := c.Query("user")
		subscription := chatroom.Subscribe()
		defer subscription.Cancel()
		events := subscription.Archive
		for i, _ := range events {
			if events[i].User == user {
				events[i].User = "you"
			}
		}
		data := struct {
			User string
			Events []chatroom.Event
		}{user,events}
		c.HTML(http.StatusOK,"refresh.html",data)
	})

	r.POST("/room", func(c *gin.Context) {
		user := c.PostForm("user")
		message := c.PostForm("message")
		chatroom.Say(user, message)
		c.Redirect(http.StatusMovedPermanently,"/refresh/room")
	})


	r.GET("/room/leave", func(c *gin.Context) {
		user := c.Query("user")
		chatroom.Leave(user)
		c.Redirect(http.StatusMovedPermanently, "/")
	})

}

