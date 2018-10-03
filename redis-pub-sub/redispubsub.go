package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

var (
	gStore      *Store
	gPubSubConn *redis.PubSubConn
	gRedisConn  = func() (redis.Conn, error) {
		return redis.Dial("tcp", ":6379")
	}
)

func init() {
	gStore = &Store{
		Users: make([]*User, 0, 1),
	}
}

// User holds  a websocket connection and unique ID
type User struct {
	ID   string
	conn *websocket.Conn
}

// Store holds the collection of users connected through websocket
type Store struct {
	Users []*User
	sync.Mutex
}

// Message gets exchanged between users through redis pub/sub messaging
// Users may have websocket connections to different nodes and stored in
// different instances of this application
type Message struct {
	DeliveryID string `json:"id"`
	Content    string `json:"content"`
}

func (s *Store) newUser(conn *websocket.Conn) *User {
	u := &User{
		ID:   uuid.NewV4().String(),
		conn: conn,
	}

	if err := gPubSubConn.Subscribe(u.ID); err != nil {
		panic(err)
	}
	s.Lock()
	defer s.Unlock()

	s.Users = append(s.Users, u)
	return u
}

var serverAddress = ":8080"

func main() {
	gRedisConn, err := gRedisConn()
	if err != nil {
		panic(err)
	}
	defer gRedisConn.Close()

	gPubSubConn = &redis.PubSubConn{Conn: gRedisConn}
	defer gPubSubConn.Close()

	// deliver messages runs on another goroutine to binded to redis pub/sub
	// messaging responsible of delivering messages to users through websocket conn
	go deliverMessages()

	http.HandleFunc("/ws", wsHandler)

	log.Printf("server started at %s\n", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, nil))
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrader error %s\n" + err.Error())
		return
	}
	u := gStore.newUser(conn)
	log.Printf("user %s joined\n", u.ID)

	for {
		var m Message

		if err := u.conn.ReadJSON(&m); err != nil {
			log.Printf("error on ws. message %s\n", err)
		}

		if c, err := gRedisConn(); err != nil {
			log.Printf("error on redis conn. %s\n", err)
		} else {
			c.Do("PUBLISH", m.DeliveryID, string(m.Content))
		}
	}
}

func deliverMessages() {
	for {
		switch v := gPubSubConn.Receive().(type) {
		case redis.Message:
			gStore.findAndDeliver(v.Channel, string(v.Data))
		case redis.Subscription:
			log.Printf("subscription message: %s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			log.Println("error pub/sub on connection, delivery has stopped")
			return
		}
	}
}

func (s *Store) findAndDeliver(userID string, content string) {
	m := Message{
		Content: content,
	}

	for _, u := range s.Users {
		if u.ID == userID {
			if err := u.conn.WriteJSON(m); err != nil {
				log.Printf("error on message delivery through ws. e: %s\n", err)
			} else {
				log.Printf("user %s found at our store, message sent\n", userID)
			}
			return
		}
	}

	log.Printf("user %s not found at our store\n", userID)
}
