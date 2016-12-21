// Package message handles types for messages in Hanayo, for instance for
// handling error messages.
package message

// A Message is a generic semantic UI message that can be returned as an error
// to a request.
type Message interface {
	Type() string
	Error() string
}

// a very simple implementation of a message.
type message struct {
	t    string // (type)
	data string
}

func (m message) Type() string {
	return m.t
}

func (m message) Error() string {
	return m.data
}

// Error returns a message, having "error" (negative) as its type.
func Error(s string) Message { return message{"error", s} }

// Warning returns a message, having "warning" as its type.
func Warning(s string) Message { return message{"warning", s} }

// Info returns a message, having "info" as its type.
func Info(s string) Message { return message{"info", s} }

// Success returns a message, having "success" (positive) as its type.
func Success(s string) Message { return message{"success", s} }
