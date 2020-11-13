package thearte

import "fmt"

type Action interface {
	Type() string
	Data() interface{}
	String() string
}

type action struct {
	id   string
	data interface{}
}

func (a *action) String() string {
	return fmt.Sprintf("{type=%s,data=%v}", a.id, a.data)
}

func (a *action) Data() interface{} {
	return a.data
}

func (a *action) Type() string {
	return a.id
}

func NewAction(id string, data interface{}) Action {
	return &action{id: id, data: data}
}
