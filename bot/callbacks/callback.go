package callbacks

import (
	"time"
)

var callbacks = make(map[int64]Callback)

func GetCallback(chatID int64) (Callback, bool) {
	a, b := callbacks[chatID]
	return a, b
}

type Callback interface {
	Start(timeout time.Duration) error
	Finish()
	Cancel()
}

type genericCallback struct {
	chatID      int64
	channel     *chan bool
	execSuccess func() error
	execFailure func() error
}

func NewCallback(chatID int64, execSuccess func() error, execFailure func() error) Callback {
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

func (g *genericCallback) Start(timeout time.Duration) error {
	if g.channel != nil {
		return nil
	}
	ch := make(chan bool, 1)
	g.channel = &ch
	defer close(*g.channel)
	defer delete(callbacks, g.chatID)
	select {
	case res := <-*g.channel:
		if !res {
			err := g.execFailure()
			if err != nil {
				return err
			}
			return nil
		}
		err := g.execSuccess()
		if err != nil {
			return err
		}
		return nil
	case <-time.After(timeout):
		err := g.execFailure()
		if err != nil {
			return err
		}
		return nil
	}
}
