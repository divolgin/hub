package hub

type BusConnection interface {
	Subscribe(topic Topic) (*BusSubscription, error)
	Listen(topic Topic) (*BusSubscription, error)
	Publish(message *Message) error
	Unsubscribe(subscriptionIds ...string)
	Request(message Message, res interface{}) error
	ServiceNameIsSet() bool
}

type BusSubscription struct {
	SubscriptionId string // unique id, used to unsubscribe
	Messages       chan *Message
}
