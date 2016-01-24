package hub

import (
	"github.com/pborman/uuid"
)

type Message struct {
	ID         string
	Topic      string
	Reply      string
	IsResponse bool
	Payload    Payload
}

type Payload struct {
	Error string `json:"error"`
	Data  []byte `json:"data"`
}

func NewDefaultMessage(opts ...func(m *Message)) *Message {
	msg := &Message{
		ID: uuid.New(),
	}

	for _, f := range opts {
		f(msg)
	}

	return msg
}
