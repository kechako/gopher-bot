package service

import (
	"github.com/kechako/gopher-bot/plugin"
)

//go:generate stringer -type=EventType

// EventType represents a type of service event.
type EventType int

const (
	UnknownEvent EventType = iota
	ConnectedEvent
	DisconnectedEvent
	MessageEvent
)

// Event represents service events.
type Event struct {
	Type EventType
	Data interface{}
}

// GetHello returns a plugin.Hello.
func (e *Event) GetHello() plugin.Hello {
	if hello, ok := e.Data.(plugin.Hello); ok {
		return hello
	}

	return nil
}

// GetHello returns a plugin.Message.
func (e *Event) GetMessage() plugin.Message {
	if msg, ok := e.Data.(plugin.Message); ok {
		return msg
	}

	return nil
}
