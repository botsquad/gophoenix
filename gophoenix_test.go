package gophoenix

import (
	"fmt"
	"sync"
	"testing"
	"time"
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
	ch, _ := t.Client.Join(channel, "bot:571a8374-9a25-45c1-90a2-5dfb4ee46ad8", params)
	channel.ch = ch
}
func (t MyConnectionReceiver) NotifyDisconnect() {
	fmt.Println("NotifyDisconnect")
}

type MyChannelReceiver struct {
	wg *sync.WaitGroup
	ch *Channel
}

func (t MyChannelReceiver) OnJoin(payload interface{}) {
	fmt.Println("OnJoin", "payload", payload)

	message := make(map[string]string)
	message["type"] = "user_message"
	message["payload"] = "hello123"

	t.ch.Push("user_action", message, func(payload interface{}) {
		fmt.Println("got response!", payload)
		time.Sleep(2 * time.Second)
		t.ch.Leave(nil)
	}, func(payload interface{}) {
		fmt.Println("got error!", payload)
	})
}
func (t MyChannelReceiver) OnJoinError(payload interface{}) {
	fmt.Println("OnJoinError", "payload", payload)
	t.wg.Done()
}
func (t MyChannelReceiver) OnChannelClose(payload interface{}) {
	fmt.Println("OnChannelClose", "payload", payload)
	t.wg.Done()
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
