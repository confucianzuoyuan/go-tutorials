package chatroom

import (
	"container/list"
	"time"
)

// 保存历史消息的条数
const archiveSize = 20

const (
	//msg
	EVENT_TYPE_MSG = iota

	EVENT_TYPE_JOIN
	EVENT_TYPE_TYPING
	EVENT_TYPE_LEAVE
	EVENT_TYPE_IMAGE

)

// 聊天室事件定义
type Event struct {
	// 事件类型
	Type      	int
	// 用户名
	User      	string
	// 时间戳
	Timestamp 	int64
	// 事件内容
	Text		string
}


// 用户订阅
type Subscription struct {

	// 历史事件
	Archive []Event

	// 事件接收通道
	// 用户从这个通道接收消息
	NewMsg <-chan Event

}

var (
	// 接收订阅事件的通道
	// 用户加入聊天室后要把历史事件推送给用户
	subscribe = make(chan (chan<- Subscription), 10)

	// 用户取消订阅通道
	// 把通道中的历史事件释放
	// 并把用户从聊天室用户列表中删除
	unsubscribe = make(chan (<-chan Event), 10)

	// 聊天室的消息推送入口
	publish = make(chan Event, 10)
)

// 取消订阅
func (s Subscription) Cancel() {
	unsubscribe <- s.NewMsg // 将用户从聊天室列表中移除
}

func newEvent(typ int , user, msg string) Event {
	return Event{typ, user, time.Now().UnixNano(), msg}
}

// 用户订阅聊天室入口函数
// 返回用户订阅的对象，用户根据对象中的属性读取历史消息和即时消息
func Subscribe() Subscription {
	resp := make(chan Subscription)
	subscribe <- resp
	return <-resp
}

// 用来向聊天室发送用户消息
func Join(user string) {
	publish <- newEvent(EVENT_TYPE_JOIN, user, "**join the room**")
}

func Say(user, message string) {
	publish <- newEvent(EVENT_TYPE_MSG, user, message)
}

func Leave(user string) {
	publish <- newEvent(EVENT_TYPE_LEAVE, user, " **leave the room**")
}

func Typing(user string)  {
	publish <- newEvent(EVENT_TYPE_TYPING,user,"**typing**")
}


// 处理聊天室中的事件
func chatroom() {

	// 历史消息列表
	archive := list.New()

	// 订阅者列表
	subscribers := list.New()

	for {
		select {

		// 当有新的订阅者加入时
		// 拿到用户订阅对象
		// 加入历史事件表
		case ch := <- subscribe:

			var events []Event

			//把历史事件加入events
			for e := archive.Front(); e != nil; e = e.Next() {
				events = append(events, e.Value.(Event))
			}

			subscriber := make(chan Event, 10)

			subscribers.PushBack(subscriber)

			ch <- Subscription{events, subscriber}

		//当有新的消息时
		case event := <-publish:

			//收到消息时，把消息推送给所有订阅者
			for ch := subscribers.Front(); ch != nil; ch = ch.Next() {
				ch.Value.(chan Event) <- event
			}

			//推送消息后，限制本地只保存10条历史消息

			if archive.Len() >= archiveSize {
				archive.Remove(archive.Front())
			}
			archive.PushBack(event)

		//当有取消订阅事件时
		case unsub := <-unsubscribe:

			//找到取消订阅的订阅者
			for ch := subscribers.Front(); ch != nil; ch = ch.Next() {
				if ch.Value.(chan Event) == unsub {
					subscribers.Remove(ch)
					break
				}
			}
		}
	}
}

// 开启goroutine loop chatroom
func init() {
	go chatroom()
}

