package streamer_test

import (
	"sync"
	"testing"
	"time"

	"github.com/jorgeolivero/hub"
	"github.com/jorgeolivero/hub/streamer"
	"github.com/pborman/uuid"
)

func TestProducerHeartbeat(t *testing.T) {
	// Setup
	bus := GetBus(t)

	topic := hub.Topic(uuid.New())
	producer := streamer.NewProducer(bus, topic)

	<-time.After(time.Millisecond * 50)

	// Send n heartbeats, expect them all to
	// be acknowledged
	n := 50
	for i := 0; i < n; i++ {
		err := bus.Request(producer.StreamInfo.HeartbeatTopic, struct{}{}, &struct{}{})
		if err != nil {
			t.Fatalf("Heartbeat request failed with error: %s", err.Error())
		}
	}
}

func TestNewProducerStreamInfo(t *testing.T) {
	bus := GetBus(t)
	topic := hub.Topic(uuid.New())
	producer := streamer.NewProducer(bus, topic)

	if len(producer.StreamInfo.StreamTopic) <= len(topic.String()) {
		t.Fatalf("Stream topic incorrectly generated")
	}
	if len(producer.StreamInfo.HeartbeatTopic) <= len(topic.String()) {
		t.Fatalf("Heartbeat topic incorrectly generated")
	}
}

// Sending an event on the stream should be
// received by multiple subscribers
func TestProducerSend(t *testing.T) {
	bus := GetBus(t)
	topic := hub.Topic(uuid.New())
	producer := streamer.NewProducer(bus, topic)

	n := 50
	wg := sync.WaitGroup{}

	// Mock consumers
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(idx int) {
			counter := 0
			mtx := &sync.Mutex{}
			bus.Listen(producer.StreamInfo.StreamTopic, func(c *hub.Context) {
				mtx.Lock()
				defer mtx.Unlock()

				counter++
				if counter == n {
					wg.Done()
				}
			})
		}(i)
	}

	// Give consumers time to boot
	time.Sleep(time.Millisecond * 100)

	// Publish n times
	for i := 0; i < n; i++ {
		go func() {
			if err := producer.Send(struct{}{}); err != nil {
				t.Fatalf("Producer send interface failed with error: %s", err.Error())
			}
		}()
	}

	wg.Wait()
}

func TestProducerClose(t *testing.T) {
	// Setup
	bus := GetBus(t)
	topic := hub.Topic(uuid.New())
	producer := streamer.NewProducer(bus, topic, func(si *streamer.StreamInfo) {
		si.HeartbeatInterval = time.Millisecond * 50
	})

	producer.Close()
	err := bus.Request(producer.StreamInfo.HeartbeatTopic, struct{}{}, &struct{}{})
	if err == nil {
		t.Fatalf("Expected heartbeat to fail after producer closed")
	}
}

func TestProducerTimeout(t *testing.T) {
	// Setup
	bus := GetBus(t)
	topic := hub.Topic(uuid.New())
	dur := time.Millisecond * 500
	producer := streamer.NewProducer(bus, topic, func(si *streamer.StreamInfo) {
		si.HeartbeatInterval = dur
	})

	<-time.After(dur + (time.Millisecond * 10))

	if producer.IsOpen() {
		t.Fatalf("Producer should be closed")
	}
}
