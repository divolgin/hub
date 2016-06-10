package streamer_test

import (
	"sync"
	"testing"
	"time"

	"github.com/jorgeolivero/hub"
	"github.com/jorgeolivero/hub/provider/nats"
	"github.com/jorgeolivero/hub/streamer"
)

func TestNewConsumerHeartbeats(t *testing.T) {
	bus := GetBus(t)
	si := GenerateStreamInfo(func(si *streamer.StreamInfo) {
		si.HeartbeatInterval = time.Millisecond * 500
	})

	numHeartbeats := 3
	wg := &sync.WaitGroup{}
	wg.Add(numHeartbeats)

	counter := 0
	bus.Subscribe(si.HeartbeatTopic.Req(), func(c *hub.Context) {
		c.Respond(struct{}{})

		counter++
		if counter <= numHeartbeats {
			wg.Done()
		}
	})

	streamer.NewConsumer(bus, si)

	wg.Wait()
}

func TestConsumerHeartbeatFailure(t *testing.T) {
	bus := GetBus(t, func(c *nats.Config) {
		c.DefaultTimeout = time.Millisecond * 500
	})
	si := GenerateStreamInfo(func(si *streamer.StreamInfo) {
		si.HeartbeatInterval = time.Millisecond * 100
	})

	consumer := streamer.NewConsumer(bus, si)

	<-time.After(bus.DefaultTimeout + si.HeartbeatInterval)

	if consumer.IsOpen() {
		t.Fatalf("Consumer should be closed")
	}
}

func TestConsumerStream(t *testing.T) {
	// Setup
	bus := GetBus(t)
	si := GenerateStreamInfo(func(si *streamer.StreamInfo) {
		si.HeartbeatInterval = time.Millisecond * 500
	})

	// Respond to heartbeats
	bus.Subscribe(si.HeartbeatTopic.Req(), func(c *hub.Context) {
		c.Respond(struct{}{})
	})

	// Start consumer
	consumer := streamer.NewConsumer(bus, si)

	// Publish
	num := 5
	go func() {
		for i := 0; i < num; i++ {
			<-time.After(time.Millisecond * 40)
			bus.Publish(si.StreamTopic, struct{}{})
		}
	}()

	// Recieve
	for i := 0; i < num; i++ {
		consumer.Next()
	}
}

func TestConsumerClose(t *testing.T) {
	// Setup
	bus := GetBus(t)
	si := GenerateStreamInfo(func(si *streamer.StreamInfo) {
		si.HeartbeatInterval = time.Millisecond * 50
	})

	heartbeats := make(chan struct{})
	events := make(chan struct{})

	bus.Subscribe(si.HeartbeatTopic.Req(), func(c *hub.Context) {
		c.Respond(struct{}{})
		heartbeats <- struct{}{}
	})

	// Start consumer
	consumer := streamer.NewConsumer(bus, si)

	consumer.Close()

	// Publish
	go func() {
		for {
			<-time.After(time.Millisecond * 10)
			bus.Publish(si.StreamTopic, struct{}{})
		}
	}()

	// Consumer
	go func() {
		for {
			consumer.Next()
			events <- struct{}{}
		}
	}()

	correctBehavior := 0
	for i := 0; i < 10; i++ {
		select {
		case <-events:
			continue
		case <-heartbeats:
			continue
		case <-time.After(time.Millisecond * 20):
			correctBehavior++
		}
	}

	if correctBehavior <= 8 {
		t.Fatalf("Closing consumer did not produce the correct behavior")
	}
}
