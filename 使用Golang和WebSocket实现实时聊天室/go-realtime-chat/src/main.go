package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket" // 常用的websocket包
)

var clients = make(map[*websocket.Conn]bool) // 已连接的客户端，map类型
var broadcast = make(chan Message)           // 用来广播消息的通道

// 配置upgrader, 也就是握手，这里允许任何客户端进行连接
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 定义消息类型
type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

func main() {
	// 创建一个简单的文件服务器, 用来托管静态文件，或者叫前端页面
	fs := http.FileServer(http.Dir("../public"))
	// 首页的url路由
	http.Handle("/", fs)

	// 配置websocket路由
	http.HandleFunc("/ws", handleConnections)

	// 监听从客户端发送到服务端的消息
	go handleMessages()

	// 本地监听8000端口，log错误
	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// 处理连接
func handleConnections(w http.ResponseWriter, r *http.Request) {
	// 建立连接
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// 做好清理工作
	defer ws.Close()

	// 将新的连接添加到clients里面
	clients[ws] = true

	for {
		var msg Message
		// 读取客户端发送来的消息，并反序列化json，然后映射为msg。
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		// 将接收到的消息广播出去
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		// 从broadcast通道获取下一条要广播的信息
		msg := <-broadcast
		// 将msg发送给所有连接到服务端的客户端
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
