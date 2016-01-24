package hub_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/jorgeolivero/hub"
	"github.com/jorgeolivero/hub/provider/nats"
	"github.com/pborman/uuid"
)

type Envelope struct {
	Foo string
	Bar string
}

func GetBusConnection(t *testing.T) hub.BusConnection {
	cfg := nats.DefaultConfig("HUB_TEST")
	conn, err := nats.NewConnection(cfg.ConnectionUrl(), cfg)
	if err != nil {
		t.Fatalf("Error creating NATs connection: %s", err.Error())
	}
	return conn
}

func TestCreateBus(t *testing.T) {
	bc := GetBusConnection(t)
	bus := hub.NewBus(bc, hub.JSON)
	if bus == nil {
		t.Fatalf("Error creating bus")
	}
}

func TestRequest(t *testing.T) {
	bc := GetBusConnection(t)
	bus := hub.NewBus(bc, hub.JSON)

	topic := hub.Topic(uuid.New())
	req := Envelope{uuid.New(), uuid.New()}

	// subscribe
	subID, err := bus.Subscribe(topic.Req(), func(c *hub.Context) {
		var data Envelope
		if err := c.Bind(&data); err != nil {
			t.Fatalf("error binding request: %s", err.Error())
		}
		if data.Bar != req.Bar || data.Foo != req.Foo {
			t.Fatalf("incorrect request received")
		}
		c.Respond(data)
	})
	if err != nil {
		t.Fatalf("Error subscribing: %s", err.Error())
	}
	defer func(subID string) {
		bus.Unsubscribe(subID)
	}(subID)

	// publish
	var res Envelope
	if err := bus.Request(topic, req, &res); err != nil {
		t.Fatalf("Error requesting to bus: %s", err.Error())
	}
	if res.Foo != req.Foo || res.Bar != res.Bar {
		t.Fatalf("Request failed, incorrect response")
	}
}

func TestRequestError(t *testing.T) {
	bc := GetBusConnection(t)
	bus := hub.NewBus(bc, hub.JSON)

	topic := hub.Topic(uuid.New())
	req := Envelope{uuid.New(), uuid.New()}
	errorValue := "EXPECTED_ERROR"

	// subscribe
	subID, err := bus.Subscribe(topic.Req(), func(c *hub.Context) {
		var data Envelope
		if err := c.Bind(&data); err != nil {
			t.Fatalf("error binding request: %s", err.Error())
		}
		if err := c.RespondError(fmt.Errorf(errorValue)); err != nil {
			t.Fatalf("Responding with error produced error: %s", err.Error())
		}
	})
	if err != nil {
		t.Fatalf("Error subscribing: %s", err.Error())
	}
	defer func(subID string) {
		bus.Unsubscribe(subID)
	}(subID)

	// publish
	println("requesting")
	var res Envelope
	if err := bus.Request(topic, req, &res); err == nil {
		t.Fatalf("Expected error")
	} else if err.Error() != errorValue {
		t.Fatalf("Response error incorrect: %s", err.Error())
	}
}

func TestPublish(t *testing.T) {
	bc := GetBusConnection(t)
	bus := hub.NewBus(bc, hub.JSON)

	topic := hub.Topic(uuid.New())
	req := Envelope{uuid.New(), uuid.New()}
	done := make(chan struct{})

	// subscribe
	subID, err := bus.Subscribe(topic.Req(), func(c *hub.Context) {
		var data Envelope
		if err := c.Bind(&data); err != nil {
			t.Fatalf("error binding request: %s", err.Error())
		}
		if data.Bar != req.Bar || data.Foo != req.Foo {
			t.Fatalf("incorrect request received")
		}
		done <- struct{}{}
	})
	if err != nil {
		t.Fatalf("Error subscribing: %s", err.Error())
	}
	defer func(subID string) {
		bus.Unsubscribe(subID)
	}(subID)

	// publish
	if err := bus.Publish(topic, req); err != nil {
		t.Fatalf("Error requesting to bus: %s", err.Error())
	}

	select {
	case <-done:
		break
	case <-time.After(time.Second * 1):
		t.Fatalf("Timed out")
	}
}

func TestListenRequest(t *testing.T) {
	bc := GetBusConnection(t)
	bus := hub.NewBus(bc, hub.JSON)

	topic := hub.Topic(uuid.New())
	req := Envelope{uuid.New(), uuid.New()}
	done := make(chan struct{})
	numListeners := 3

	// create 3 listeners
	i := numListeners
	for i > 0 {
		subID, err := bus.Listen(topic.Req(), func(c *hub.Context) {
			var data Envelope
			if err := c.Bind(&data); err != nil {
				t.Fatalf("error binding request: %s", err.Error())
			}
			if data.Bar != req.Bar || data.Foo != req.Foo {
				t.Fatalf("incorrect request received")
			}
			done <- struct{}{}
		})
		if err != nil {
			t.Fatalf("Error listening: %s", err.Error())
		}
		defer func(subID string) {
			bus.Unsubscribe(subID)
		}(subID)

		i--
	}

	// publish
	if err := bus.Publish(topic, req); err != nil {
		t.Fatalf("Error requesting to bus: %s", err.Error())
	}

	x := 3
	for x > 0 {
		select {
		case <-done:
			break
		case <-time.After(time.Second * 1):
			t.Fatalf("Timed out")
		}
		x--
	}
}

func TestListenResponse(t *testing.T) {
	bc := GetBusConnection(t)
	bus := hub.NewBus(bc, hub.JSON)

	topic := hub.Topic(uuid.New())
	req := Envelope{uuid.New(), uuid.New()}
	done := make(chan struct{})
	numListeners := 3

	// create 3 listeners
	i := numListeners
	for i > 0 {
		subID, err := bus.Listen(topic.ResWildcard(), func(c *hub.Context) {
			var data Envelope
			if err := c.Bind(&data); err != nil {
				t.Fatalf("error binding request: %s", err.Error())
			}
			if data.Bar != req.Bar || data.Foo != req.Foo {
				t.Fatalf("incorrect request received")
			}
			done <- struct{}{}
		})
		if err != nil {
			t.Fatalf("Error listening: %s", err.Error())
		}
		defer func(subID string) {
			bus.Unsubscribe(subID)
		}(subID)

		i--
	}

	// subscribe
	subID, err := bus.Subscribe(topic.Req(), func(c *hub.Context) {
		var data Envelope
		if err := c.Bind(&data); err != nil {
			t.Fatalf("error binding request: %s", err.Error())
		}
		if data.Bar != req.Bar || data.Foo != req.Foo {
			t.Fatalf("incorrect request received")
		}
		c.Respond(data)
	})
	if err != nil {
		t.Fatalf("Error subscribing: %s", err.Error())
	}
	defer func(subID string) {
		bus.Unsubscribe(subID)
	}(subID)

	// publish
	if err := bus.Request(topic, req, &struct{}{}); err != nil {
		t.Fatalf("Error requesting to bus: %s", err.Error())
	}

	x := 3
	for x > 0 {
		select {
		case <-done:
			break
		case <-time.After(time.Second * 1):
			t.Fatalf("Timed out")
		}
		x--
	}
}

func TestSubscribe(t *testing.T) {}

func TestUnsubscribe(t *testing.T) {}
