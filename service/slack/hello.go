package slack

import (
	"github.com/kechako/gopher-bot/v2/plugin"
	"golang.org/x/exp/slog"
)

type hello struct {
	bot *bot
}

var _ plugin.Hello = (*hello)(nil)

// newHello returns a new *hello as plugin.Hello.
func newHello(service *slackService) plugin.Hello {
	return &hello{
		bot: &bot{
			service: service,
		},
	}
}

// Bot implements the plugin.Hello interface.
func (h *hello) Bot() plugin.Bot {
	return h.bot
}

type bot struct {
	service *slackService
}

var _ plugin.Bot = (*bot)(nil)

// Logger implements the plugin.Bot interface.
func (b *bot) Logger() *slog.Logger {
	return b.service.l
}

// UserID implements the plugin.Bot interface.
func (b *bot) UserID() string {
	return b.service.UserID()
}

// Post implements the plugin.Bot interface.
func (b *bot) Post(channelID string, text string) {
	b.service.Post(channelID, text)

}

// Mention implements the plugin.Bot interface.
func (b *bot) Mention(channelID, userID, text string) {
	b.service.Mention(channelID, userID, text)
}

// ProcessCommand implements the plugin.Bot interface.
func (b *bot) ProcessCommand(channelID string, command string) {
	b.service.ProcessCommand(channelID, command)
}

// Channel implements the plugin.Bot interface.
func (b *bot) Channel(channelID string) plugin.Channel {
	return b.service.Channel(channelID)
}

// User implements the plugin.Bot interface.
func (b *bot) User(userID string) plugin.User {
	return b.service.User(userID)
}
