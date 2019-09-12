package slack

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/kechako/gopher-bot/plugin"
	"github.com/kechako/gopher-bot/service"
	"github.com/kechako/gopher-bot/service/slack/internal/msgfmt"
	"github.com/nlopes/slack"
	"sync"
)

// slackService represents a service for Slack.
type slackService struct {
	client *slack.Client
	rtm    *slack.RTM

	ch chan *service.Event

	wg   sync.WaitGroup
	exit context.CancelFunc
}

// New returns a new Slack service as service.Service.
func New(token string) (service.Service, error) {
	if token == "" {
		return nil, errors.New("the token is empty")
	}

	client := slack.New(token)

	s := &slackService{
		client: client,
		rtm:    client.NewRTM(),
	}

	return s, nil
}

// Start implements the service.Service interface.
func (s *slackService) Start(ctx context.Context) (<-chan *service.Event, error) {
	_, err := s.client.AuthTestContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate slack service: %w", err)
	}

	s.ch = make(chan *service.Event)

	go s.rtm.ManageConnection()

	ctx, cancel := context.WithCancel(ctx)
	s.exit = cancel

	s.wg.Add(1)
	go s.loop(ctx)

	return s.ch, nil
}

func (s *slackService) loop(ctx context.Context) {
	defer s.wg.Done()

loop:
	for {
		select {
		case msg := <-s.rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				// TODO : logging
				s.handleHello()
			case *slack.ConnectedEvent:
				// TODO : logging
			case *slack.DisconnectedEvent:
				// TODO : logging
				if ev.Intentional {
					break loop
				}
			case *slack.MessageEvent:
				// TODO : logging
				s.handleMessage(ev)
			}
		}
	}

}

// handleHello handles the hello event.
func (s *slackService) handleHello() {
	s.ch <- &service.Event{
		Type: service.ConnectedEvent,
		Data: newHello(s),
	}
}

// handleMessage handles the message event.
func (s *slackService) handleMessage(msg *slack.MessageEvent) {
	if msg.User == s.UserID() {
		// bot message
		return
	}

	s.ch <- &service.Event{
		Type: service.MessageEvent,
		Data: newMessage(s, msg),
	}
}

// Start implements the service.Service interface.
func (s *slackService) Close() error {
	return s.rtm.Disconnect()
}

// UserID implements the service.Service interface.
func (s *slackService) UserID() string {
	return s.rtm.GetInfo().User.ID
}

// Post implements the service.Service interface.
func (s *slackService) Post(channelID, text string) {
	s.PostToThread(channelID, text, "")
}

// Post implements the service.Service interface.
func (s *slackService) PostToThread(channelID, text string, ts string) {
	s.rtm.SendMessage(
		s.rtm.NewOutgoingMessage(text, channelID, slack.RTMsgOptionTS(ts)))
}

// Mention implements the service.Service interface.
func (s *slackService) Mention(channelID, userID, text string) {
	s.MentionToThread(channelID, userID, text, "")
}

// Mention implements the service.Service interface.
func (s *slackService) MentionToThread(channelID, userID, text string, ts string) {
	s.PostToThread(channelID, msgfmt.Format(
		&msgfmt.Block{
			Type:    msgfmt.UserBlock,
			Content: userID,
		},
		msgfmt.SpaceBlock,
		&msgfmt.Block{
			Type:    msgfmt.TextBlock,
			Content: text,
		},
	), ts)
}

// ProcessCommmand processes the specified command on the channel.
func (s *slackService) ProcessCommand(channelID string, command string) {
	go func() {
		s.ch <- &service.Event{
			Type: service.MessageEvent,
			Data: newCommandMessage(s, channelID, command),
		}
	}()
}

// Channel returns a channel of specified channelID.
func (s *slackService) Channel(channelID string) plugin.Channel {
	if len(channelID) == 0 || channelID[0] != 'C' {
		// ignore group
		return nil
	}

	ch, err := s.client.GetChannelInfo(channelID)
	if err != nil {
		// TODO: logging
		return nil
	}

	return &channel{
		id:   ch.ID,
		name: ch.Name,
	}
}

// EscapeHelp implements the service.Service interface.
func (s *slackService) EscapeHelp(help string) string {
	escaped := bytes.NewBuffer(make([]byte, len(help)+8))
	escaped.Reset()

	escaped.WriteString("```\n")
	escaped.WriteString(help)
	if !strings.HasSuffix(help, "\n") {
		escaped.WriteRune('\n')
	}
	escaped.WriteString("```")

	return escaped.String()
}
