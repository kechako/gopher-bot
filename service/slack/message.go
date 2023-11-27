package slack

import (
	"log/slog"
	"strings"

	"github.com/kechako/gopher-bot/v2/plugin"
	"github.com/kechako/gopher-bot/v2/service/slack/internal/msgfmt"
	"github.com/slack-go/slack/slackevents"
)

type message struct {
	service  *slackService
	msg      *slackevents.MessageEvent
	blocks   []*msgfmt.Block
	mentions []string
	text     string
}

var _ plugin.Message = (*message)(nil)

// newMessage returns a new *message as plugin.Message.
func newMessage(service *slackService, msg *slackevents.MessageEvent) plugin.Message {
	m := &message{
		service: service,
		msg:     msg,
		blocks:  msgfmt.Parse(msg.Text),
	}
	m.init()

	return m
}

func (m *message) init() {
	var s strings.Builder
	for _, b := range m.blocks {
		switch b.Type {
		case msgfmt.ChannelBlock, msgfmt.UserBlock:
			if b.Label == "" {
				user, err := m.service.client.GetUserInfo(b.Content)
				if err != nil {
					m.service.l.Error("failed to get user info", slog.String("content", b.Content))
				} else {
					b.Label = user.Name
				}
			}
		}
		s.WriteString(b.String())
	}
	m.text = s.String()

	for _, b := range m.blocks {
		if b.Type == msgfmt.UserBlock {
			m.mentions = append(m.mentions, b.Content)
		}
	}
}

// ChannelID implements the plugin.Message interface.
func (m *message) ChannelID() string {
	return m.msg.Channel
}

// UserID implements the plugin.Message interface.
func (m *message) UserID() string {
	return m.msg.User
}

// Text implements the plugin.Message interface.
func (m *message) Text() string {
	return m.text
}

// Post implements the plugin.Message interface.
func (m *message) Post(text string) {
	m.service.PostToThread(m.ChannelID(), text, m.msg.ThreadTimeStamp)
}

// Mention implements the plugin.Message interface.
func (m *message) Mention(text string) {
	m.service.MentionToThread(m.ChannelID(), m.UserID(), text, m.msg.ThreadTimeStamp)
}

// Mentions implements the plugin.Message interface.
func (m *message) Mentions() []string {
	if len(m.mentions) == 0 {
		return nil
	}

	mentions := make([]string, len(m.mentions))

	copy(mentions, m.mentions)

	return mentions
}

// MentionTo implements the plugin.Message interface.
func (m *message) MentionTo(userID string) bool {
	for _, mentionTo := range m.mentions {
		if userID == mentionTo {
			return true
		}
	}

	return false
}

// PostHelp implements the plugin.Message interface.
func (m *message) PostHelp(help *plugin.Help) {
	msg := m.service.EscapeHelp(help.String())

	m.Post(msg)
}
