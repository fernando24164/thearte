package thearte

import (
	"sync"
)

type Actor interface {
	Name() string
	When(string, Dispatcher) Actor
	ListSubs() map[string]Subscription
}

type BaseActor struct {
	mailbox chan Action
	name    string
	subs    map[string]Subscription
	L       sync.RWMutex
}

func NewActor(name string) Actor {
	return &BaseActor{name: name, subs: make(map[string]Subscription)}
}

func (a *BaseActor) Name() string {
	return a.name
}

func (a *BaseActor) When(topic string, dispatch Dispatcher) Actor {
	a.L.Lock()
	a.subs[topic] = newSubscription(topic, dispatch)
	a.L.Unlock()
	return a
}

func (a *BaseActor) ListSubs() map[string]Subscription {
	return a.subs
}
