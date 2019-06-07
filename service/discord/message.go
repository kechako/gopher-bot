package discord

import (
	discord "github.com/bwmarrin/discordgo"
	"github.com/kechako/gopher-bot/plugin"
)

type message struct {
	service *discordService
	msg     *discord.Message
}

var _ plugin.Message = (*message)(nil)

// newMessage returns a new *message as plugin.Message.
func newMessage(service *discordService, msg *discord.Message) plugin.Message {
	return &message{
		service: service,
		msg:     msg,
	}
}

// ChannelID implements the plugin.Message interface.
func (m *message) ChannelID() string {
	return m.msg.ChannelID
}

// UserID implements the plugin.Message interface.
func (m *message) UserID() string {
	return m.msg.Author.ID
}

// Text implements the plugin.Message interface.
func (m *message) Text() string {
	return m.msg.Content
}

// Post implements the plugin.Message interface.
func (m *message) Post(text string) {
	m.service.Post(m.ChannelID(), text)
}

// Mention implements the plugin.Message interface.
func (m *message) Mention(text string) {
	m.service.Mention(m.ChannelID(), m.UserID(), text)
}

// Mentions implements the plugin.Message interface.
func (m *message) Mentions() []string {
	mentions := make([]string, 0, len(m.msg.Mentions))

	for _, user := range m.msg.Mentions {
		mentions = append(mentions, user.ID)
	}

	return mentions
}

// MentionTo implements the plugin.Message interface.
func (m *message) MentionTo(userID string) bool {
	for _, user := range m.msg.Mentions {
		if userID == user.ID {
			return true
		}
	}

	return false
}
