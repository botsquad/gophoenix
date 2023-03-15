package gophoenix

import (
	"errors"
)

// Client is the entry point for a phoenix channel connection.
type Client struct {
	t  Transport
	mr *messageRouter
	cr ConnectionReceiver
	rc refCounter
}

// NewWebsocketClient creates the default connection using a websocket as the transport.
func NewWebsocketClient(cr ConnectionReceiver) *Client {
	return &Client{
		t:  &socketTransport{},
		mr: newMessageRouter(),
		cr: cr,
		rc: &atomicRef{ref: new(int64)},
	}
}

// Connect should be called to establish the connection through the transport.
func (c *Client) Connect(url string) error {
	if c.t == nil {
		return errors.New("transport not provided")
	}

	return c.t.Connect(url, c.mr, c.cr)
}

// Close closes the connection via the transport.
func (c *Client) Close() error {
	if c.t == nil {
		return errors.New("transport not provided")
	}

	c.t.Close()

	return nil
}

// Join subscribes to a channel via the transport and returns a reference to the channel.
func (c *Client) Join(callbacks ChannelReceiver, topic string, payload interface{}) (*Channel, error) {
	if c.t == nil {
		return nil, errors.New("transport not provided")
	}

	rr := newReplyRouter()
	joinRef := c.rc.nextRef()
	ch := &Channel{
		joinRef: joinRef,
		topic:   topic,
		t:       c.t,
		rc:      &atomicRef{ref: new(int64)},
		rr:      rr,
		ln:      func() { c.mr.unsubscribe(joinRef) },
	}
	c.mr.subscribe(joinRef, callbacks, rr)
	err := ch.join(callbacks, payload)

	if err != nil {
		return nil, err
	}

	return ch, nil
}
