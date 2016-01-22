package hub

type BusConnection interface {
	Subscribe(topic string) (*Subscription, error)
	Listen(topic Topic) (*Subscription, error)
	Publish(message *Message) error
	Unsubscribe(subscriptionIds ...string)
	Request(message Message, res interface{}) error
	ServiceNameIsSet() bool
}
