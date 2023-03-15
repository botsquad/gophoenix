package gophoenix

import (
	"fmt"
	"sync"
	"testing"
)

type MyConnectionReceiver struct {
	Client *Client
	wg     *sync.WaitGroup
}

func (t MyConnectionReceiver) NotifyConnect() {
	fmt.Println("NotifyConnect")
	channel := &MyChannelReceiver{wg: t.wg}
	params := make(map[string]string)
	params["user_id"] = "arjan"
	t.Client.Join(channel, "bot:571a8374-9a25-45c1-90a2-5dfb4ee46ad8", params)
}
func (t MyConnectionReceiver) NotifyDisconnect() {
	fmt.Println("NotifyDisconnect")
}

type MyChannelReceiver struct {
	wg *sync.WaitGroup
}

func (t MyChannelReceiver) OnJoin(payload interface{}) {
	fmt.Println("OnJoin", "payload", payload)
	t.wg.Done()
}
func (t MyChannelReceiver) OnJoinError(payload interface{}) {
	fmt.Println("OnJoinError", "payload", payload)
}
func (t MyChannelReceiver) OnChannelClose(payload interface{}) {
	fmt.Println("OnChannelClose", "payload", payload)
}
func (t MyChannelReceiver) OnMessage(event string, payload interface{}) {
	fmt.Println("OnMessage", event, payload)
}

func TestSocket(t *testing.T) {
	var wg sync.WaitGroup

	receiver := &MyConnectionReceiver{wg: &wg}
	client := NewWebsocketClient(receiver)
	//	defer client.Close()

	receiver.Client = client

	wg.Add(1)

	client.Connect("wss://arjan.ngrok.io/socket/websocket?vsn=2.0.0")
	fmt.Println("done")
	wg.Wait()
	fmt.Println("done2")
	client.Close()
}
