package redis

import (
	"fmt"
	"testing"
	"time"
)

func TestSub(t *testing.T) {
	re := NewRedisEntry("172.24.42.83:16378", "")
	fmt.Println(time.Now().Unix())
	msg, err := re.SubscribeAndReceiveMessage("sub_test", 10*time.Second)
	fmt.Println(time.Now().Unix())
	fmt.Println(msg)
	fmt.Println(err)
}
