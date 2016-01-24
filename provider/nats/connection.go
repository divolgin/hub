package nats

import (
	"fmt"
	"time"

	"github.com/apcera/nats"
	"github.com/jorgeolivero/hub"
	"github.com/pborman/uuid"
)

type Config struct {
	User, Password, Host, Port string
	Service                    string
	DefaultTimeout             time.Duration
}

func DefaultConfig(service string) *Config {
	return &Config{
		User:           "",
		Password:       "",
		Host:           "localhost",
		Port:           "4222",
		Service:        service,
		DefaultTimeout: time.Second * 5,
	}
}

func (config *Config) ConnectionUrl() string {
	if config.Host == "localhost" {
		return nats.DefaultURL
	}
	return fmt.Sprintf(
		"nats://%s:%s@%s:%s",
		config.User,
		config.Password,
		config.Host,
		config.Port)
}

type Connection struct {
	Config        *Config
	Connection    *nats.EncodedConn
	Subscriptions map[string]*Subscription
	ReconnectChan chan bool
}

type Subscription struct {
	MsgChan      chan *hub.Message
	Subscription *nats.Subscription
}

// As a matter of course, Nats connections should have a queue group name.
// However, providing an empty string for group, will allow the client to
// create singular nats connections on which to subscribe
func NewConnection(url string, config *Config) (*Connection, error) {
	// set up connetion options
	reconnectChan := make(chan bool)
	opts := nats.Options{
		Url:     url,
		Timeout: config.DefaultTimeout,
		ReconnectedCB: func(c *nats.Conn) {
			reconnectChan <- true
		},
		DisconnectedCB: func(c *nats.Conn) {
			// restart the process
			panic(fmt.Errorf("nats disconnected"))
		},
	}

	// connect
	conn, err := opts.Connect()
	if err != nil {
		return nil, err
	}

	// get encoded connection
	ec, err := nats.NewEncodedConn(conn, nats.GOB_ENCODER)
	if err != nil {
		return nil, err
	}

	// wrap
	natsConn := &Connection{
		Config:        config,
		Connection:    ec,
		Subscriptions: make(map[string]*Subscription),
		ReconnectChan: reconnectChan,
	}
	return natsConn, nil
}

func (nc *Connection) IsOpen() bool {
	return !nc.Connection.Conn.IsClosed()
}

func (nc *Connection) Listen(subject string) (*hub.Subscription, error) {
	// create chan
	msgChan := make(chan *hub.Message)

	// start subscription
	sub, err := nc.Connection.Subscribe(subject, func(msg *hub.Message) {
		msgChan <- msg
	})
	if err != nil {
		return nil, err
	}

	// store and return
	subID := uuid.New()
	nc.Subscriptions[subID] = &Subscription{
		MsgChan:      msgChan,
		Subscription: sub,
	}
	return &hub.Subscription{
		ID:       subID,
		Messages: msgChan,
	}, nil
}

func (nc *Connection) Subscribe(subject string) (*hub.Subscription, error) {
	// create chan
	msgChan := make(chan *hub.Message)

	// start subscription
	sub, err := nc.Connection.QueueSubscribe(subject, nc.Config.Service, func(msg *hub.Message) {
		// push it through the message chan
		msgChan <- msg
	})
	if err != nil {
		return nil, err
	}

	// store and return
	subID := uuid.New()
	nc.Subscriptions[subID] = &Subscription{
		MsgChan:      msgChan,
		Subscription: sub,
	}
	return &hub.Subscription{
		ID:       subID,
		Messages: msgChan,
	}, nil
}

func (nc *Connection) Unsubscribe(subscriptionIds ...string) {
	for _, sid := range subscriptionIds {
		if sub, ok := nc.Subscriptions[sid]; ok {
			close(sub.MsgChan)
			if err := sub.Subscription.Unsubscribe(); err != nil {
				panic(err)
			}
			delete(nc.Subscriptions, sid)
		}
	}
}

func (nc *Connection) Publish(msg *hub.Message) error {
	return nc.Connection.Publish(msg.Topic, msg)
}

func (nc *Connection) Request(msg *hub.Message) error {
	return nc.Publish(msg)
}

func (nc *Connection) GetTimeout() time.Duration {
	return nc.Config.DefaultTimeout
}

func (nc *Connection) GetNumActiveSubscriptions() int {
	return len(nc.Subscriptions)

}

func (nc *Connection) GetReconnectNotifyChan() chan bool {
	return nc.ReconnectChan
}

func (nc *Connection) ServiceNameIsSet() bool {
	return len(nc.Config.Service) > 0
}
