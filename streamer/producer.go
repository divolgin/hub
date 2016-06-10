package streamer

import (
	"fmt"
	"time"

	"github.com/jorgeolivero/hub"
)

type Producer struct {
	Bus        *hub.Bus
	StreamInfo StreamInfo

	isClosed            bool
	closeHeartbeatsChan chan struct{}
}

func NewProducer(bus *hub.Bus, topic hub.Topic, opts ...func(si *StreamInfo)) *Producer {
	// Generate stream topic and heartbeat topic
	si := NewStreamInfo(topic)
	for _, f := range opts {
		f(&si)
	}

	p := &Producer{
		Bus:                 bus,
		StreamInfo:          si,
		isClosed:            false,
		closeHeartbeatsChan: make(chan struct{}),
	}

	go p.handleHeartbeats()

	return p
}

func (p *Producer) handleHeartbeats() {
	// Subscribe to heartbeats, and respond when they come in
	heartbeats := make(chan struct{}, 100)
	subID, err := p.Bus.Subscribe(p.StreamInfo.HeartbeatTopic.Req(), func(c *hub.Context) {
		c.Respond(struct{}{}) // Ack
		heartbeats <- struct{}{}
	})
	if err != nil {
		panic(fmt.Errorf("Subscribing to heartbeats [%s] failed with error: %s", p.StreamInfo.HeartbeatTopic, err.Error()))
	}
	defer func() {
		p.Bus.Unsubscribe(subID)
	}()

	// If we do not receive heartbeats in the heartbeat interval
	// unsubscribe to heartbeats & set `isClosed` to true
	for {
		select {
		case <-heartbeats:
			continue
		case <-time.After(p.StreamInfo.HeartbeatInterval):
			p.Close()
			return
		case <-p.closeHeartbeatsChan:
			return
		}
	}
}

func (p *Producer) Send(event interface{}) error {
	return p.Bus.Publish(p.StreamInfo.StreamTopic, event)
}

func (p *Producer) IsOpen() bool {
	return !p.isClosed
}

func (p *Producer) Close() {
	if p.isClosed {
		return
	}

	p.isClosed = true
	p.closeHeartbeatsChan <- struct{}{}
}
