package diff

import (
	"testing"
	"fmt"
)

func TestReplay(_ *testing.T) {
	log, err := NewACPCLog("test-game.log")
	if err != nil {
		panic(err)
	}
	c := make(chan *ACPCLog)
	go log.Replay(c)
	fmt.Println(log.hands[4])
	for i := 0; i < 20; i++ {
		fmt.Println(<-c)
	}
}
