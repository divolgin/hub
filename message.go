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

func NewReqMessage(topic Topic, req interface{}, serializer BusSerializer, withReply bool) (*Message, error) {
	// reply string
	reply := ""
	if withReply {
		reply = topic.ResUnique()
	}

	// encode request
	b, err := serializer.Serialize(req)
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
			Data:  b,
		},
	}
}
