package thearte

type Subscription interface {
	Dispatcher() Dispatcher
	Topic() string
}

type subscription struct {
	topic      string
	dispatcher Dispatcher
}

func (s *subscription) Dispatcher() Dispatcher {
	return s.dispatcher
}

func (s *subscription) Topic() string {
	return s.topic
}

func newSubscription(topic string, dispatcher Dispatcher) Subscription {
	return &subscription{topic: topic, dispatcher: dispatcher}
}
