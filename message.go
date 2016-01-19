package hub

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

func NewMessage(topic Topic, req interface{}) Message {
	// return Message{topic.String(), request}
	return Message{}
}
