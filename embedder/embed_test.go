package embedder

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestNewEmbed(t *testing.T) {
	embed := NewEmbed("etcd")
	err := embed.Init(context.Background(), WithID("odin0"), WithDir("../disk"), WithIP("127.0.0.1"), WithClientPort("21389"), WithPeerPort("21390"), WithCluster([]string{"127.0.0.1"}))
	fmt.Println(err)
	t.Error(err)
	ready := make(chan struct{})
	go func() {
		err := embed.Run(ready)
		fmt.Println(err)
		t.Error(err)
	}()
	select {
	case <-ready:
		err = embed.SetAuth("root", "613f#8d164df4ACPF49@93a510df49!66f98b*d6")
		fmt.Println(err)
		t.Error(err)
	}
	time.Sleep(10 * time.Second)

}
