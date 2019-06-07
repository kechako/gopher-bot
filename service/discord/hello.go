package discord

import (
	"github.com/kechako/gopher-bot/plugin"
)

type hello struct {
	bot *bot
}

var _ plugin.Hello = (*hello)(nil)

// newHello returns a new *hello as plugin.Hello.
func newHello(service *discordService) plugin.Hello {
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
	service *discordService
}

var _ plugin.Bot = (*bot)(nil)

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
