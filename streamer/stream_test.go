package streamer_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/jorgeolivero/hub"
	"github.com/jorgeolivero/hub/streamer"
	"github.com/pborman/uuid"
)

func TestStream(t *testing.T) {
	bus := GetBus(t)
	topic := hub.Topic(uuid.New())
	producer := streamer.NewProducer(bus, topic)
	consumer := streamer.NewConsumer(bus, producer.StreamInfo)

	n := 20

	// Stream n events
	go func() {
		<-time.After(time.Millisecond * 50)
		for i := 0; i < n; i++ {
			err := producer.Send(i)
			if err != nil {
				t.Fatalf("Send event to stream failed with error: %s", err.Error())
			}
		}
	}()

	for i := 0; i < n; i++ {
		c := consumer.Next()
		var count int
		if err := c.Bind(&count); err != nil {
			t.Fatalf("Binding stream count failed with error: %s", err.Error())
		}
		if i != count {
			fmt.Printf("Expected: %d Got: %d\n", i, count)
		}
	}
}

func TestStreamConsumerClose(t *testing.T) {
	bus := GetBus(t)
	topic := hub.Topic(uuid.New())

	timeout := time.Millisecond * 30
	producer := streamer.NewProducer(bus, topic, func(si *streamer.StreamInfo) {
		si.HeartbeatInterval = timeout
	})
	consumer := streamer.NewConsumer(bus, producer.StreamInfo)
	consumer.Close()

	<-time.After(timeout + (time.Millisecond * 5))

	if producer.IsOpen() {
		t.Fatalf("Producer should have timed out")
	}
}

func TestStreamProducerClose(t *testing.T) {
	bus := GetBus(t)
	topic := hub.Topic(uuid.New())

	timeout := time.Millisecond * 30
	producer := streamer.NewProducer(bus, topic, func(si *streamer.StreamInfo) {
		si.HeartbeatInterval = timeout
	})
	consumer := streamer.NewConsumer(bus, producer.StreamInfo)

	producer.Close()

	<-time.After(bus.DefaultTimeout + producer.StreamInfo.HeartbeatInterval)
	if consumer.IsOpen() {
		t.Fatalf("Expect consumer to be closed")
	}
}
