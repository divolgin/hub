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

func NewRequestMessage(topic, reply string, req interface{}, withReply bool, serializer BusSerializer) (*Message, error) {
	data, err := serializer.Serialize(req)
	if err != nil {
		return nil, err
	}

	return &Message{
		ID:        uuid.New(),
		Topic:     topic,
		Reply:     reply,
		IsReponse: false,
		Payload: Payload{
			Error: "",
			Data:  data,
		},
	}
}
