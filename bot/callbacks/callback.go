package callbacks

import (
	"fmt"
	"time"
)

var callbacks = make(map[int64]Callback)

func GetCallback(chatID int64) (Callback, bool) {
	a, b := callbacks[chatID]
	return a, b
}

type Callback interface {
	Start(timeout time.Duration)
	Finish()
	Cancel()
}

type genericCallback struct {
	chatID      int64
	channel     *chan bool
	execSuccess func()
	execFailure func()
}

func NewCallback(chatID int64, execSuccess func(), execFailure func()) Callback {
	cb := genericCallback{execSuccess: execSuccess, execFailure: execFailure, chatID: chatID}
	callbacks[chatID] = &cb
	return &cb
}

func (g *genericCallback) Finish() {
	*g.channel <- true
}

func (g *genericCallback) Cancel() {
	*g.channel <- false
}

func (g *genericCallback) Start(timeout time.Duration) {
	if g.channel != nil {
		return
	}
	ch := make(chan bool, 1)
	g.channel = &ch
	defer close(*g.channel)
	defer delete(callbacks, g.chatID)
	select {
	case res := <-*g.channel:
		if !res {
			fmt.Println("---------------------IM DED")
			g.execFailure()
			return
		}
		g.execSuccess()
		return
	case <-time.After(timeout):
		g.execFailure()
		return
	}
}
