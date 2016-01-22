package hub

type Context struct {
	Message *Message
	Bus     *Bus
}

// Binds the payload to the provided data store.
func (c *Context) Bind(receiver interface{}) error {
	return nil
}

// Responds using the reply inbox held in the context.
func (c *Context) Respond(res interface{}) error {
	return nil
}

// Responds using the provided topic.
func (c *Context) RespondToTopic(topic Topic, res interface{}) error {
	return nil
}

// Responds with an error using the reply inbox held in the context.
func (c *Context) RespondError(err error) error {
	return nil
}

// Responds with an error using the provided topic.
func (c *Context) RespondErrorToTopic(topic Topic, err error) error {
	return nil
}
