package hub

type BusConnection interface {
	Subscribe(topic Topic) (*BusSubscription, error)
	Listen(topic Topic) (*BusSubscription, error)
	Publish(message *Message) error
	Unsubscribe(subscriptionIds ...string)
	Request(message Message, res interface{}) error
	ServiceNameIsSet() bool
}
