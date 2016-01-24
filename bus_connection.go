package hub

type BusConnection interface {
	Subscribe(topic string) (*Subscription, error)
	Listen(topic string) (*Subscription, error)
	Publish(message *Message) error
	Unsubscribe(subscriptionIds ...string)
	Request(message *Message) error
	ServiceNameIsSet() bool
}
