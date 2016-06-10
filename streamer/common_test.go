package streamer_test

import (
	"testing"
	"time"

	"github.com/jorgeolivero/hub"
	"github.com/jorgeolivero/hub/provider/nats"
	"github.com/jorgeolivero/hub/streamer"
	"github.com/pborman/uuid"
)

func GetBus(t *testing.T, opts ...func(c *nats.Config)) *hub.Bus {
	cfg := &nats.Config{
		User:           "",
		Password:       "",
		Host:           "192.168.99.100",
		Port:           "4222",
		Service:        "STREAMER_TEST",
		DefaultTimeout: time.Second * 5,
	}

	for _, f := range opts {
		f(cfg)
	}

	conn, err := nats.NewConnection(cfg.ConnectionUrl(), cfg)
	if err != nil {
		t.Fatalf("Error creating NATs connection: %s", err.Error())
	}
	return hub.NewBus(conn, hub.JSON)
}

func GenerateStreamInfo(opts ...func(si *streamer.StreamInfo)) streamer.StreamInfo {
	si := streamer.StreamInfo{
		StreamTopic:       hub.Topic(uuid.New()).Stream(),
		HeartbeatTopic:    hub.Topic(uuid.New()).Heartbeat(),
		HeartbeatInterval: (time.Second * 2),
	}
	for _, f := range opts {
		f(&si)
	}
	return si
}
