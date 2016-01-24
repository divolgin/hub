package hub

import (
	"fmt"
)

type Context struct {
	message *Message
	bus     *Bus
}

// Binds the payload to the provided data store.
func (c *Context) Bind(receiver interface{}) error {
	return c.bus.serializer.Deserialize(c.message.Payload.Data, receiver)
}

// Responds using the reply inbox held in the context.
func (c *Context) Respond(res interface{}) error {
	// preconditions
	if !c.IsReplyable() {
		return fmt.Errorf("Respond OP not allowed: %#v\n", c.message)
	}

	// create message
	data, err := c.bus.serializer.Serialize(res)
	if err != nil {
		return err
	}
	msg := NewDefaultMessage(func(m *Message) {
		m.Topic = c.message.Reply
		m.Reply = ""
		m.IsResponse = true
		m.Payload.Data = data
	})

	// publish
	return c.bus.Connection.Publish(msg)
}

// Responds using the provided topic.
func (c *Context) RespondToTopic(topic Topic, res interface{}) error {
	return nil
}

// Responds with an error using the reply inbox held in the context.
func (c *Context) RespondError(err error) error {
	// pre conditions
	if !c.IsReplyable() {
		return fmt.Errorf("RespondError OP not allowed: %#v\n", c.message)
	}

	// create message
	msg := NewDefaultMessage(func(m *Message) {
		m.Topic = c.message.Reply
		m.Reply = ""
		m.IsResponse = true
		m.Payload.Error = err.Error()
	})

	return c.bus.Connection.Publish(msg)
}

// Responds with an error using the provided topic.
func (c *Context) RespondErrorToTopic(topic Topic, err error) error {
	return nil
}

func (c *Context) IsReplyable() bool {
	if c.message.IsResponse || len(c.message.Reply) == 0 {
		return false
	}
	return true
}
