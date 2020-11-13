package main

import (
	"context"
	"log"
	"time"

	thearte "github.com/sonirico/thearte/pkg"
)

type pingActor struct {
	thearte.Actor
}

func newPingActor() *pingActor {
	return &pingActor{Actor: thearte.NewActor("pinger")}
}

func (p *pingActor) Ping(ctx context.Context, a thearte.Action, emit thearte.Emitter) {
	time.Sleep(time.Second)
	emit(thearte.NewAction("ping", nil))
}

type pongActor struct {
	thearte.Actor
}

func newPongActor() *pongActor {
	return &pongActor{Actor: thearte.NewActor("ponger")}
}

func (p *pongActor) Pong(ctx context.Context, a thearte.Action, emit thearte.Emitter) {
	time.Sleep(time.Second)
	emit(thearte.NewAction("pong", nil))
}

func main() {
	done := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	stage := thearte.NewStage()
	ping := newPingActor()
	pong := newPongActor()

	ping.When("pong", ping.Ping)
	pong.When("ping", pong.Pong)

	stage.Register(ping, pong)
	go stage.Start(ctx)

	stage.Dispatch(thearte.NewAction("ping", nil))

	go func() {
		<-time.After(time.Second * 10)
		log.Println("canceling")
		cancel()
		<-time.After(time.Second * 10)
		done <- struct{}{}
	}()
	<-done
}
