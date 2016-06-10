package streamer

import (
	"fmt"
	"sync"
	"time"

	"github.com/jorgeolivero/hub"
)

const (
	RING_BUFFER_LIMIT = 50
)

type Consumer struct {
	Bus        *hub.Bus
	StreamInfo StreamInfo

	closeLock             *sync.Mutex
	isClosed              bool
	closeSubscriptionChan chan struct{}
	closeRingChan         chan struct{}
	closeHeartbeatsChan   chan struct{}

	stream chan *hub.Context
	output chan *hub.Context
}

func NewConsumer(bus *hub.Bus, info StreamInfo) *Consumer {
	c := &Consumer{
		Bus:                   bus,
		StreamInfo:            info,
		closeLock:             &sync.Mutex{},
		isClosed:              false,
		closeSubscriptionChan: make(chan struct{}, 1),
		closeRingChan:         make(chan struct{}, 1),
		closeHeartbeatsChan:   make(chan struct{}, 1),
		stream:                make(chan *hub.Context),
		output:                make(chan *hub.Context),
	}

	go c.startRing()
	go c.subscribe()
	go c.handleHeartbeats()

	return c
}

func (c *Consumer) startRing() {
	buffer := make(chan *hub.Context, RING_BUFFER_LIMIT)
	go func() {
		for {
			c.output <- <-buffer
		}
	}()
	for {
		select {
		case buffer <- <-c.stream:
			if len(buffer) == cap(buffer) {
				// Potential deadlock here
				<-buffer
			}
		case <-c.closeRingChan:
			break
		}
	}
	fmt.Printf("Stream ring for topic [%s] closed\n", c.StreamInfo.StreamTopic)
}

func (c *Consumer) subscribe() {
	subID, err := c.Bus.Listen(c.StreamInfo.StreamTopic, func(cc *hub.Context) {
		c.stream <- cc
	})
	if err != nil {
		panic(fmt.Errorf("Subscribing to stream [%s] failed with error: %s", c.StreamInfo.StreamTopic, err.Error()))
	}

	<-c.closeSubscriptionChan

	c.Bus.Unsubscribe(subID)
	fmt.Printf("Stream subscription for topic [%s] closed\n", c.StreamInfo.StreamTopic)
}

func (c *Consumer) handleHeartbeats() {
	for {
		select {
		case <-time.After((c.StreamInfo.HeartbeatInterval / 3)):
			err := c.Bus.Request(c.StreamInfo.HeartbeatTopic, struct{}{}, &struct{}{})
			if err != nil {
				fmt.Printf("Heartbeat request failed for stream [%s] with error: %s\n", c.StreamInfo.StreamTopic, err.Error())
				c.Close()
			}
		case <-c.closeHeartbeatsChan:
			return
		}
	}
}

func (c *Consumer) Next() *hub.Context {
	return <-c.output
}

func (c *Consumer) IsOpen() bool {
	return !c.isClosed
}

func (c *Consumer) Close() {
	c.closeLock.Lock()
	defer c.closeLock.Unlock()

	if c.isClosed {
		fmt.Printf("Attempt to close closed stream consumer for topic [%s] short circuited\n", c.StreamInfo.StreamTopic)
		return
	}
	c.isClosed = true

	c.closeRingChan <- struct{}{}
	c.closeSubscriptionChan <- struct{}{}
	c.closeHeartbeatsChan <- struct{}{}

	return
}
