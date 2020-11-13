# Thearte

**Thearte** (intelligently misspelled from 'theatre') is a lightweight, minimal implementation
of the actor model system in golang. Now, depending upon your background, you may ask:

- Why to embedded `redux` in golang?
- Why to create an `event emitter` in golang with javascript already has?
- Why an actor model when golang is CSP based? or
- Why all this "bunga bunga" to do something that is relatively easy to create by
  employing the built-in communication & synchronization primitives that go already
  offers?
  
There is no magic. The underlying mechanism are go channels. I did NOT code this
because I think that using channels in a standard fashion tend to couple things together, but because complex
UML diagrams of interdependent pieces, or let's say large codebases, do not get along
with go. When your code starts filling up with strategies, senders, adapters, factories, 
clients, interceptors... it starts to feel natural just keep small and compact and connect all these pieces together
by leveraging, you guessed, channels. It's the same reason as using pub/sub systems to segregate
business domain logic into microservices, where logic is decoupled because it needs to be scaled. If
that same concept applied to the internals of such microservices, your code may look like this
example savaged from the `examples` folder.

```golang
package main

import (
	"context"
	"log"
	"time"

	"github.com/sonirico/thearte/pkg"
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

    // Start the mandanga
	stage.Dispatch(thearte.NewAction("ping", nil))

	go func() {
		<-time.After(time.Second * 10)
		log.Println("canceling")
		cancel()
		<-time.After(time.Second * 10)
		done <- struct{}{}
	}()
	<- done
}
```

## Disclaimer

Either way, **this project is a prototype** so not yet ready to use in production. There
is no code coverage nor intend to add it soon. Use it at your own risk, although we both
know you wouldn't use it ;P
