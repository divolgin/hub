package hub

import (
// "github.com/pborman/uuid"
)

type Bus struct {
	Connection    BusConnection
	serializer    BusSerializer
	subscriptions []*Subscription
}

func NewBus(bc BusConnection, format SerializationFormat) *Bus {
	return &Bus{
		Connection: bc,
		serializer: format.GetSerializer(),
	}
}

// ListenRequest will subscribe to the given topic response, and invoke
// the handler when messages are received. Processs of the same service
// listening to the same topic will all receive every message.
func (b *Bus) ListenResponse(topic TopicResponse, handle MessageHandler) error {

}

// HandleRequest will subscribe to the given topic request, and
// invoke the handler when messages are received. Processes of the same
// service listening to same topic will be have requests round robbined between
// the group.
func (b *Bus) HandleRequest(topic TopicRequest, handler MessageHandler) error {}

// Request will publish a request to the provided topic and wait for a response.
// If the request produces an error, an error will be returned.
func (b *Bus) Request(topic TopicRequest, req, res interface{}) error {
	return nil
}

// Subscribe will invoke the provided handler with messages directed towards
// the provided topic
func (b *Bus) Subscribe(topic Topic, handler MessageHandler) (string, error) {
	sub, err := b.Connection.Subscribe()
	if err != nil {
		panic(err)
	}

	go func() {
		for message := range sub.Messages {
			context := &Context{
				Bus:     b,
				Message: message,
			}
			go handler(context)
		}
	}()

	b.subscriptions = append(b.subscriptions, &subscription{
		id:      sub.SubscriptionId,
		subject: messageSubject,
		uuid:    uuid.New(),
		handler: handler,
	})

	// return sub.SubscriptionId
	return "", nil
}

// Publish the request to the provided topic.
// This method does not wait for a response, it is fire and forget.
func (b *Bus) Publish(topic Topic, req interface{}) error {
	msg, err := NewReqMessage(topic.String(), req, b.serializer, false)
	if err != nil {
		return err
	}

	return b.Connection.Publish(msg)
}

// Unsubscribe will cancel the subscriptions for each of the
// provided subscriptions ids.
func (b *Bus) Unsubscribe(subscriptionIDs ...string) {
	b.Connection.Unsubscribe(subscriptionIDs...)
}
