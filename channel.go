package gophoenix

import (
	"fmt"
	"strconv"
)

// Channel represents a subscription to a topic. It is returned from the Client after joining a topic.
type Channel struct {
	topic   string
	t       Transport
	rc      refCounter
	rr      *replyRouter
	joinRef int64
}

type refCounter interface {
	nextRef() int64
}

// Leave notifies the channel to unsubscribe from messages on the topic.
func (ch *Channel) Leave(payload interface{}) error {
	ref := ch.rc.nextRef()
	ch.rr.subscribe(ref, func(p interface{}) {
		fmt.Println("leave OK")
	}, func(p interface{}) {
		fmt.Println("leave fail")
	})
	return ch.sendMessage(ref, LeaveEvent, payload)
}

// Push sends a message on the topic.
func (ch *Channel) Push(event string, payload interface{}, replyHandler replyCallback, errorHandler replyCallback) error {
	ref := ch.rc.nextRef()
	ch.rr.subscribe(ref, replyHandler, errorHandler)
	return ch.sendMessage(ref, event, payload)
}

// PushNoReply sends a message on the topic but does not provide a callback to receive replies.
func (ch *Channel) PushNoReply(event string, payload interface{}) error {
	ref := ch.rc.nextRef()
	return ch.sendMessage(ref, event, payload)
}

func (ch *Channel) join(callbacks ChannelReceiver, payload interface{}) error {
	ref := ch.rc.nextRef()
	ch.rr.subscribe(ref, callbacks.OnJoin, callbacks.OnJoinError)
	return ch.sendMessage(ref, JoinEvent, payload)
}

func (ch *Channel) sendMessage(ref int64, event string, payload interface{}) error {
	var m [5]interface{}

	m[0] = strconv.Itoa(int(ch.joinRef))
	m[1] = strconv.Itoa(int(ref))
	m[2] = ch.topic
	m[3] = event
	m[4] = payload

	return ch.t.Push(m)
}
