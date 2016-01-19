// Tests require a local instance of NATs running
// on the default ports.
package nats_test

import (
	"testing"
	"time"

	"github.com/jorgeolivero/hub"
	"github.com/jorgeolivero/hub/provider/nats"
	"github.com/pborman/uuid"
)

func GetNatsTestConnection(t *testing.T) *nats.Connection {
	cfg := nats.DefaultConfig("CONNECTION_TEST")
	conn, err := nats.NewConnection(cfg.ConnectionUrl(), cfg)
	if err != nil {
		t.Fatalf("Error creating NATs connection: %s", err.Error())
	}
	return conn
}

func GenerateMsg(opts ...func(m *hub.Message)) *hub.Message {
	topic := uuid.New()
	msg := &hub.Message{
		ID:         uuid.New(),
		Topic:      topic,
		Reply:      topic + ".RES",
		IsResponse: false,
		Payload:    hub.Payload{},
	}

	for _, o := range opts {
		o(msg)
	}

	return msg
}

func TestCreatingConnection(t *testing.T) {
	conn := GetNatsTestConnection(t)
	if !conn.IsOpen() {
		t.Fatalf("Unable to create an open connection")
	}
}

// TestSubscribeTopic ensures that an open subscriptions to a topic
// is able to recieve messages.
func TestSubscribeTopic(t *testing.T) {
	msg := GenerateMsg()

	// subscribe
	conn := GetNatsTestConnection(t)
	sub, err := conn.Subscribe(msg.Topic)
	if err != nil {
		t.Fatalf("Error subscribing to subject")
	}

	// publish
	if err := conn.Publish(msg); err != nil {
		t.Fatalf("Error publishing message")
	}

	// wait
	select {
	case rec := <-sub.Messages:
		if rec.ID != msg.ID {
			t.Fatalf("Received incorrect message")
		}
	case <-time.After(3 * time.Second):
		t.Fatalf("Subscription timed out")
	}
}

// TestMultiSubscribeTopic ensures that when multiple subscriptions
// are made from the same service, messages are round robbined between
// the group.
func TestMultiSubscribeTopic(t *testing.T) {
	msg := GenerateMsg()

	// subscribe
	conn := GetNatsTestConnection(t)
	sub1, err := conn.Subscribe(msg.Topic)
	if err != nil {
		t.Fatalf("Error subscribing to subject")
	}
	sub2, err := conn.Subscribe(msg.Topic)
	if err != nil {
		t.Fatalf("Error subscribing to subject")
	}

	// publish
	if err := conn.Publish(msg); err != nil {
		t.Fatalf("Error publishing message")
	}

	// ensure that only one subscriber received the message
	i := 2
	passed := false
	for i > 0 {
		select {
		case rec := <-sub1.Messages:
			if rec.ID != msg.ID {
				t.Fatalf("Wrong message received")
			}
		case rec := <-sub2.Messages:
			if rec.ID != msg.ID {
				t.Fatalf("Wrong message received")
			}
		case <-time.After(1 * time.Second):
			passed = true
		}
		i--
	}
	if !passed {
		t.Fatalf("Both subscriptions received the message")
	}
}

// TestListenTopic ensures that an open listener to a topic is able
// to receive messages.
func TestListenTopic(t *testing.T) {
	msg := GenerateMsg()

	// subscribe
	conn := GetNatsTestConnection(t)
	sub, err := conn.Listen(msg.Topic)
	if err != nil {
		t.Fatalf("Error subscribing to subject")
	}

	// publish

	if err := conn.Publish(msg); err != nil {
		t.Fatalf("Error publishing message")
	}

	// wait
	select {
	case rec := <-sub.Messages:
		if rec.ID != msg.ID {
			t.Fatalf("Received incorrect message")
		}
	case <-time.After(3 * time.Second):
		t.Fatalf("Subscription timed out")
	}
}

// TestMultiListenTopic ensures that when multiple listners are made
// from the same service to the same topic, each listener receives all
// messages sent to that topic.
func TestMultiListenTopic(t *testing.T) {
	msg := GenerateMsg()

	// subscribe
	conn := GetNatsTestConnection(t)
	sub1, err := conn.Listen(msg.Topic)
	if err != nil {
		t.Fatalf("Error subscribing to subject")
	}
	sub2, err := conn.Listen(msg.Topic)
	if err != nil {
		t.Fatalf("Error subscribing to subject")
	}

	// publish
	if err := conn.Publish(msg); err != nil {
		t.Fatalf("Error publishing message")
	}

	// ensure that only one subscriber received the message
	i := 2
	passed := true
	for i > 0 {
		select {
		case rec := <-sub1.Messages:
			if rec.ID != msg.ID {
				t.Fatalf("Wrong message received")
			}
		case rec := <-sub2.Messages:
			if rec.ID != msg.ID {
				t.Fatalf("Wrong message received")
			}
		case <-time.After(1 * time.Second):
			passed = false
		}
		i--
	}
	if !passed {
		t.Fatalf("Both subscriptions received the message")
	}
}

func TestUnsubscribe(t *testing.T) {
	msg := GenerateMsg()

	// subscribe
	conn := GetNatsTestConnection(t)
	sub, err := conn.Listen(msg.Topic)
	if err != nil {
		t.Fatalf("Error subscribing to subject")
	}

	// publish
	if err := conn.Publish(msg); err != nil {
		t.Fatalf("Error publishing message")
	}

	// wait
	select {
	case rec := <-sub.Messages:
		if rec.ID != msg.ID {
			t.Fatalf("Received incorrect message")
		}
	case <-time.After(3 * time.Second):
		t.Fatalf("Subscription timed out")
	}

	// unsubscribe
	conn.Unsubscribe(sub.ID)

	// publish
	if err := conn.Publish(msg); err != nil {
		t.Fatalf("Error publishing message")
	}

	// wait
	select {
	case _, ok := <-sub.Messages:
		if ok {
			t.Fatalf("Received message after unsubscribing")
		}
	case <-time.After(1 * time.Second):
		break
	}
}
