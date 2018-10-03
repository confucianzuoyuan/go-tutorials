package main

import (
	"chat-room/routes"
	"github.com/gin-gonic/gin"
	"net/http"
)

func index() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	}
}
func demo() gin.HandlerFunc {
	return func(c *gin.Context) {
		demo := c.Query("demo")
		user := c.Query("user")
		switch demo {
		case "refresh":
			c.Redirect(http.StatusMovedPermanently, "/refresh?user=" + user)
		case "longpolling":
			c.Redirect(http.StatusMovedPermanently, "/longpolling/room?user=" + user)
		case "websocket":
			c.Redirect(http.StatusMovedPermanently,"/websocket/room?user=" + user)
		default:
			c.Redirect(http.StatusMovedPermanently, "/websocket/room?user=" + user)
		}
	}
}

func main() {

	s := gin.Default()

	s.GET("/", index())

	s.GET("/demo", demo())


	routes.Refresh(s)

	routes.LongPolling(s)

	routes.WebSocket(s)


	s.LoadHTMLGlob("./templates/*")

	s.Static("/static", "./static")

	s.Run()

}
