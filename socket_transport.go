package gophoenix

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type socketTransport struct {
	socket *websocket.Conn
	cr     ConnectionReceiver
	mr     MessageReceiver
	close  chan struct{}
	done   chan struct{}
}

func (st *socketTransport) Connect(url string, mr MessageReceiver, cr ConnectionReceiver) error {
	st.mr = mr
	st.cr = cr

	socket, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		return err
	}

	st.socket = socket
	go st.pingLoop()
	go st.listen()
	st.cr.NotifyConnect()

	return err
}

func (st *socketTransport) Push(data interface{}) error {
	return st.socket.WriteJSON(data)
}

func (st *socketTransport) Close() {
	st.close <- struct{}{}
	<-st.done
}

func (st *socketTransport) pingLoop() {
	for {
		select {
		case <-st.close:
			return
		default:
		}

		time.Sleep(10 * time.Second)

		st.socket.WriteControl(websocket.PingMessage, []byte(""), time.Now().Add(time.Second))
		fmt.Println("ping")
	}
}

func (st *socketTransport) listen() {
	defer st.stop()
	for {
		select {
		case <-st.close:
			return
		default:
		}

		msgType, data, err := st.socket.ReadMessage()

		if err != nil {
			fmt.Println(err)
			continue
		}

		if msgType == websocket.PingMessage {
			fmt.Println("got ping")
			continue
		}

		if msgType != websocket.TextMessage {
			continue
		}

		var arr []interface{}
		err = json.Unmarshal([]byte(data), &arr)

		if err != nil {
			fmt.Printf("JSON error")
			fmt.Println(err)
			continue
		}

		if len(arr) != 5 {
			fmt.Println("Protocol error")
			continue
		}

		if arr[0] == nil || arr[2] == nil || arr[3] == nil {
			fmt.Printf("data: %s\n", data)

			fmt.Println("Protocol error, missing joinref, topic or event")
			continue
		}
		joinRef, _ := strconv.Atoi(arr[0].(string))

		var ref int64 = 0
		if arr[1] != nil {
			rr, _ := strconv.Atoi(arr[1].(string))
			ref = int64(rr)
		}

		msg := Message{}
		msg.Ref = ref
		msg.JoinRef = int64(joinRef)
		msg.Topic = arr[2].(string)
		msg.Event = arr[3].(string)
		msg.Payload = arr[4]

		st.mr.NotifyMessage(&msg)
	}
}

func (st *socketTransport) stop() {
	st.socket.Close()
	st.cr.NotifyDisconnect()
	func() { st.done <- struct{}{} }()
}
