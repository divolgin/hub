package hub

import (
	"fmt"
	"sync"
	"time"

	"github.com/pborman/uuid"
)

type Bus struct {
	Connection        BusConnection
	serializer        BusSerializer
	subscriptionsLock sync.RWMutex
	subscriptions     map[string]*Subscription
	DefaultTimeout    time.Duration
}

func NewBus(bc BusConnection, format SerializationFormat) *Bus {
	return &Bus{
		Connection: bc,
		serializer: format.GetSerializer(),
	}
}

// Request will publish a request to the provided topic and wait for a response.
// If the request produces an error, an error will be returned.
func (b *Bus) Request(topic Topic, req, res interface{}) error {
	msg, err := NewRequestMessage(topic.Req().String(), topic.ResUnique().String(), req, true, b.serializer)
	if err != nil {
		return err
	}

	// subscribe to response
	sub, err := b.Connection.Subscribe(msg.Reply)
	if err != nil {
		return err
	}
	defer func(sub *Subscription) {
		b.Unsubscribe(sub.ID)
	}(sub)

	// send request
	err := b.Connection.Publish(msg)
	if err != nil {
		return err
	}

	// get response or timeout
	select {
	case msg := <-sub.Messages:
		if len(msg.Payload.Error) > 0 {
			return fmt.Errorf(msg.Payload.Error)
		}
		err := b.serializer.Deserialize(msg.Payload.Data, res)
		if err != nil {
			return fmt.Errorf("Error deserializing response: %s", err.Error())
		}
		return nil
	case <-time.After(b.DefaultTimeout):
		return fmt.Errorf("Request timed out")
	}
}

// Subscribe will invoke the provided handler with messages directed towards
// the provided topic
func (b *Bus) Subscribe(topic Topic, handler MessageHandler) (string, error) {
	sub, err := b.Connection.Subscribe(topic.String())
	if err != nil {
		return "", err
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

	// save subscription in map
	b.subscriptionsLock.Lock()
	b.subscriptions[sub.ID] = sub
	b.subscriptionsLock.Unlock()

	// return sub.SubscriptionId
	return sub.ID, nil
}

// Publish the request to the provided topic.
// This method does not wait for a response, it is fire and forget.
func (b *Bus) Publish(topic Topic, req interface{}) error {
	msg, err := NewRequestMessage(topic.Req(), "", req, false, b.serializer)
	if err != nil {
		return err
	}

	return b.Connection.Publish(msg)
}

// Unsubscribe will cancel the subscriptions for each of the
// provided subscriptions ids.
func (b *Bus) Unsubscribe(subscriptionIDs ...string) {
	b.Connection.Unsubscribe(subscriptionIDs...)

	// delete from sub map
	b.subscriptionsLock.Lock()
	for _, subID := range subscriptionIDs {
		delete(b.subscriptions[subID])
	}
	b.subscriptionsLock.Unlock()
}
