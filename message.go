package gophoenix

// Message is a message sent or received via the Transport from the channel.
type Message struct {
	Topic   string
	Event   string
	Payload map[string]interface{}
	Ref     int64
	JoinRef int64
}
