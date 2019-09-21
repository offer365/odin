package node

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestRunRpcServer(t *testing.T) {
	node := NewNode("odin0", "127.0.0.1")
	go RunRpcServer("1111", node)

	time.Sleep(1 * time.Second)

	for i := 0; i < 100; i++ {
		ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*800)
		node, err := GetRemoteNode(ctx, "odin0", "127.0.0.1", "11111")
		fmt.Println(err)
		byt, err := json.Marshal(node)
		fmt.Println(string(byt))
		time.Sleep(1 * time.Second)
	}
}
