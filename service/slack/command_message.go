package slack

import (
	"github.com/kechako/gopher-bot/v2/plugin"
)

type commandMessage struct {
	service   *slackService
	channelID string
	command   string
}

var _ plugin.Message = (*message)(nil)

// newMessage returns a new *message as plugin.Message.
func newCommandMessage(service *slackService, channelID string, command string) plugin.Message {
	return &commandMessage{
		service:   service,
		channelID: channelID,
		command:   command,
	}
}

// ChannelID implements the plugin.Message interface.
func (m *commandMessage) ChannelID() string {
	return m.channelID
}

// UserID implements the plugin.Message interface.
func (m *commandMessage) UserID() string {
	return m.service.UserID()
}

// Text implements the plugin.Message interface.
func (m *commandMessage) Text() string {
	return m.command
}

// Post implements the plugin.Message interface.
func (m *commandMessage) Post(text string) {
	m.service.Post(m.ChannelID(), text)
}

// Mention implements the plugin.Message interface.
func (m *commandMessage) Mention(text string) {
	m.service.Mention(m.ChannelID(), m.UserID(), text)
}

// Mentions implements the plugin.Message interface.
func (m *commandMessage) Mentions() []string {
	return nil
}

// MentionTo implements the plugin.Message interface.
func (m *commandMessage) MentionTo(userID string) bool {
	return false
}

// PostHelp implements the plugin.Message interface.
func (m *commandMessage) PostHelp(help *plugin.Help) {
	msg := m.service.EscapeHelp(help.String())

	m.Post(msg)
}
