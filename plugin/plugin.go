package plugin

import (
	"context"
)

// Plugin is the interface implemented by bot plugins.
type Plugin interface {
	// Hello is called when the service is ready.
	Hello(ctx context.Context, h Hello)
	// DoAction is called when the message is received.
	DoAction(ctx context.Context, m Message)
	// Help returns a plugin help information.
	Help(ctx context.Context) *Help
}

// Hello is the interface to get bot information.
type Hello interface {
	Bot() Bot
}

// Bot is the interface that represents bot information.
type Bot interface {
	// UserID returns the user ID of the bot.
	UserID() string
	// Post posts a new message to the channel.
	Post(channelID string, text string)
	// Mention posts a new message that mentions to the user to the channel.
	Mention(channelID, userID, text string)
}

// Bot is the interface that represents a service message.
type Message interface {
	// ChannelID returns ID of the channel that the message was posted.
	ChannelID() string
	// UserID returns ID of the user that posted the message.
	UserID() string
	// Text is a text of the message.
	Text() string
	// Post posts a new message to the channel that the message was posted.
	Post(text string)
	// Mention posts a new message that mentions to the user that posted the message,
	// to the channel that the message was posted.
	Mention(text string)
	// Mentions returns user IDs that message mentions to.
	Mentions() []string
	// MentionTo returns whether the message mentions to the userID.
	// Returns true if the message mentions to the userID, otherwise returns false.
	MentionTo(userID string) bool
}

// Help represents a help information of a plugin.
type Help struct {
	Name        string
	Description string
	Commands    []*Command
}

// Command represents a command for a plugin.
type Command struct {
	Command     string
	Description string
}
