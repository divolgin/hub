package hub

type Subscription struct {
	ID       string
	Messages chan *Message
}
