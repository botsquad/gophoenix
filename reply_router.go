package gophoenix

import (
	"sync"
)

type replyRouter struct {
	mapLock sync.RWMutex
	rr      map[int64]replyCallback
	er      map[int64]replyCallback
}

type replyCallback func(payload interface{})

func newReplyRouter() *replyRouter {
	return &replyRouter{
		rr: make(map[int64]replyCallback),
		er: make(map[int64]replyCallback),
	}
}

func (rr *replyRouter) routeReply(msg *Message) {
	rr.mapLock.RLock()
	rc, ok1 := rr.rr[msg.Ref]
	ec, ok2 := rr.er[msg.Ref]
	rr.mapLock.RUnlock()

	if !ok1 || !ok2 {
		return
	}

	rr.mapLock.Lock()
	delete(rr.rr, msg.Ref)
	delete(rr.er, msg.Ref)
	rr.mapLock.Unlock()

	if msg.Payload["status"] == "ok" {
		rc(msg.Payload["response"])
	} else {
		ec(msg.Payload["response"])
	}
}

func (rr *replyRouter) subscribe(ref int64, callback replyCallback, errorCallback replyCallback) {
	rr.mapLock.Lock()
	defer rr.mapLock.Unlock()
	rr.rr[ref] = callback
	rr.er[ref] = errorCallback
}
