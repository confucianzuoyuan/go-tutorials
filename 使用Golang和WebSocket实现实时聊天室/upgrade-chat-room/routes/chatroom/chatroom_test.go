package chatroom

import (
	"fmt"
	"testing"
	"time"
)

func TestJoin(t *testing.T) {

	subsc := Subscribe()

	go func(subscription Subscription) {
		for ev := range subscription.NewMsg {
			fmt.Println(ev.User,ev.Text)
		}
	}(subsc)

	Join("awc")
	Typing("awc")
	Join("heiheihei")
	Join("csk")

	Say("awc","nice day!")
	Typing("csk")
	Say("csk","yeah")
	Say("heiheihei","bye!")
	Leave("heiheihei")

	time.Sleep(1 * time.Second)


}
