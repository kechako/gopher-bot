package service

import (
	"context"
)

// Service is the interface implemented by types that provides bot functions.
type Service interface {
	// Start starts the bot service.
	Start(ctx context.Context) (<-chan *Event, error)
	// Close closes a bot session.
	Close() error
	// UserID returns the user ID of the bot.
	UserID() string
	// Post posts a new message to the channel.
	Post(channelID string, text string)
	// Mention posts a new message that mentions to the user to the channel.
	Mention(channelID, userID, text string)
	// EscapeHelp escapes help document.
	EscapeHelp(help string) string
}
