package gophoenix

import (
	"sync"
)

type messageRouter struct {
	mapLock sync.RWMutex
	tr      map[int64]*topicReceiver
	sub     chan ChannelReceiver
}

type topicReceiver struct {
	cr ChannelReceiver
	rr *replyRouter
}

func newMessageRouter() *messageRouter {
	return &messageRouter{
		tr:  make(map[int64]*topicReceiver),
		sub: make(chan ChannelReceiver),
	}
}

func (mr *messageRouter) NotifyMessage(msg *Message) {
	mr.mapLock.RLock()
	tr, ok := mr.tr[msg.JoinRef]
	mr.mapLock.RUnlock()
	if !ok {
		return
	}

	switch msg.Event {
	case ReplyEvent:
		tr.rr.routeReply(msg)
	case CloseEvent:
		tr.cr.OnChannelClose(msg.Payload)
		mr.unsubscribe(msg.JoinRef)
	default:
		tr.cr.OnMessage(msg.Event, msg.Payload)
	}
}

func (mr *messageRouter) subscribe(joinRef int64, cr ChannelReceiver, rr *replyRouter) {
	mr.mapLock.Lock()
	defer mr.mapLock.Unlock()
	mr.tr[joinRef] = &topicReceiver{cr: cr, rr: rr}
}

func (mr *messageRouter) unsubscribe(joinRef int64) {
	mr.mapLock.Lock()
	defer mr.mapLock.Unlock()
	delete(mr.tr, joinRef)
}
