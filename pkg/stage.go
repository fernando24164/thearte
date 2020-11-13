package thearte

import (
	"context"
	"log"
	"sync"
)

type status byte

const (
	_ status = iota
	stopped
	started
	paused
)

type Emitter func(action Action)
type Dispatcher func(ctx context.Context, a Action, e Emitter)

type Stage interface {
	//UnRegisterByName(name string) Stage
	//UnRegister(actor Actor) Stage
	Dispatch(a Action) Stage
	Register(actor ...Actor) Stage
	Start(context.Context)
}

type stage struct {
	bus      chan Action
	actors   map[string]Actor
	actorsMu sync.RWMutex
	topics   map[string][]Subscription
	topicsMu sync.RWMutex
	status   status
	signals  chan status
}

func (s *stage) Register(actors ...Actor) Stage {
	s.actorsMu.Lock()
	for _, actor := range actors {
		s.actors[actor.Name()] = actor
	}
	s.actorsMu.Unlock()
	return s
}

func (s *stage) accept(msg Action) {
	s.bus <- msg
}

func (s *stage) Pause() {
	s.signals <- paused
}

func (s *stage) Dispatch(a Action) Stage {
	s.bus <- a
	return s
}

func (s *stage) Start(ctx context.Context) {
	s.status = started
	s.actorsMu.Lock()
	for _, actor := range s.actors {
		for topic, sub := range actor.ListSubs() {
			s.topics[topic] = append(s.topics[topic], sub)
		}
		//go func(a Actor) {
		//	a.Start(ctx)
		//}(actor)
	}
	s.actorsMu.Unlock()
	for {
		select {
		case <-ctx.Done():
			return
		case state := <-s.signals:
			log.Println("new status")
			s.status = state
		case msg := <-s.bus:
			log.Println("dispatching", msg)
			if s.status != started {
				break
			}
			subs, ok := s.topics[msg.Type()]
			if ok {
				for _, sub := range subs {
					go func(su Subscription) {
						dispatch := su.Dispatcher()
						dispatch(ctx, msg, s.accept)
					}(sub)
				}
			}
		}
	}
}

func NewStage() Stage {
	return &stage{
		bus:    make(chan Action),
		actors: make(map[string]Actor),
		topics: make(map[string][]Subscription),
		status: stopped,
	}
}
