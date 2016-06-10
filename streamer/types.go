package streamer

import (
	"time"

	"github.com/jorgeolivero/hub"
)

type StreamInfo struct {
	StreamTopic       hub.Topic     `json:"topic"`
	HeartbeatTopic    hub.Topic     `json:"heartbeatTopic"`
	HeartbeatInterval time.Duration `json:"heartbeatInterval"`
}

func NewStreamInfo(topic hub.Topic, opts ...func(si *StreamInfo)) StreamInfo {
	si := StreamInfo{
		StreamTopic:       topic.Stream(),
		HeartbeatTopic:    topic.Heartbeat(),
		HeartbeatInterval: time.Second * 6,
	}

	for _, f := range opts {
		f(&si)
	}

	return si
}
