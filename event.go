package gophoenix

const (
	// MessageEvent represents a regular message on a topic.
	MessageEvent string = "phx_message"
	// JoinEvent represents a successful join on a channel.
	JoinEvent string = "phx_join"
	// CloseEvent represents the closing of a channel.
	CloseEvent string = "phx_close"
	// ErrorEvent represents an error.
	ErrorEvent string = "phx_error"
	// ReplyEvent represents a reply to a message sent on a topic.
	ReplyEvent string = "phx_reply"
	// LeaveEvent represents leaving a channel and unsubscribing from a topic.
	LeaveEvent string = "phx_leave"
)
