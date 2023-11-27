// Package slack provides a bot service of Slack.
package slack

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/kechako/gopher-bot/v2/plugin"
	"github.com/kechako/gopher-bot/v2/service"
	"github.com/kechako/gopher-bot/v2/service/slack/internal/msgfmt"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type Config struct {
	Logger *slog.Logger
}

func (cfg *Config) logger() *slog.Logger {
	var l *slog.Logger
	if cfg != nil {
		l = cfg.Logger
	}
	if l != nil {
		return l
	}
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}

// slackService represents a service for Slack.
type slackService struct {
	client *slack.Client
	socket *socketmode.Client
	botID  string
	userID string
	teamID string

	l *slog.Logger

	ch chan *service.Event

	wg   sync.WaitGroup
	exit context.CancelFunc
}

// New returns a new Slack service as service.Service.
func New(token, appToken string, cfg *Config) (service.Service, error) {
	if token == "" {
		return nil, errors.New("the token is empty")
	}

	client := slack.New(token, slack.OptionAppLevelToken(appToken))

	s := &slackService{
		client: client,
		socket: socketmode.New(client),
		l:      cfg.logger(),
	}

	return s, nil
}

// Start implements the service.Service interface.
func (s *slackService) Start(ctx context.Context) (<-chan *service.Event, error) {
	s.l.Info("Start Slack bot service")

	auth, err := s.client.AuthTestContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate slack service: %w", err)
	}

	s.botID = auth.BotID
	s.userID = auth.UserID
	s.teamID = auth.TeamID
	s.l.Info("authenticated", slog.String("bot_id", s.botID), slog.String("user_id", s.userID), slog.String("team_id", s.teamID))

	s.ch = make(chan *service.Event)

	ctx, cancel := context.WithCancel(ctx)
	s.exit = cancel

	s.wg.Add(1)
	go s.loop(ctx)

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		err := s.socket.RunContext(ctx)
		if err != nil && !errors.Is(err, context.Canceled) {
			s.l.Error("failed to run socketmode", slog.Any("err", err))
		}
	}()

	return s.ch, nil
}

func (s *slackService) loop(ctx context.Context) {
	defer s.wg.Done()

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case event := <-s.socket.Events:
			switch event.Type {
			case socketmode.EventTypeHello:
				s.l.Info("bot is ready")
				s.handleHello()
			case socketmode.EventTypeConnected:
				s.l.Info("bot has connected")
			case socketmode.EventTypeDisconnect:
				s.l.Info("bot has disconnected")
			case socketmode.EventTypeIncomingError:
				s.l.Info("bot received invomming error")
			case socketmode.EventTypeEventsAPI:
				eventsAPIEvent, ok := event.Data.(slackevents.EventsAPIEvent)
				if !ok {
					continue
				}

				s.socket.Ack(*event.Request)

				innerEvent := eventsAPIEvent.InnerEvent

				switch slackevents.EventsAPIType(innerEvent.Type) {
				case slackevents.Message:
					ev, ok := innerEvent.Data.(*slackevents.MessageEvent)
					if !ok {
						continue
					}
					s.handleMessage(ev)
				}
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
func (s *slackService) handleMessage(msg *slackevents.MessageEvent) {
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
	s.exit()
	s.wg.Wait()
	return nil
}

// UserID implements the service.Service interface.
func (s *slackService) UserID() string {
	return s.userID
}

// Post implements the service.Service interface.
func (s *slackService) Post(channelID, text string) {
	s.PostToThread(channelID, text, "")
}

// Post implements the service.Service interface.
func (s *slackService) PostToThread(channelID, text string, ts string) {
	s.client.PostMessage(channelID, slack.MsgOptionText(text, false), slack.MsgOptionTS(ts))
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
	if len(channelID) == 0 {
		// ignore group
		return nil
	}

	ch, err := s.client.GetConversationInfo(&slack.GetConversationInfoInput{
		ChannelID: channelID,
	})
	if err != nil {
		s.l.Error("Failed to get channel info", slog.String("channel_id", channelID), slog.Any("err", err))
		return nil
	}

	return &channel{
		id:   ch.ID,
		name: ch.Name,
	}
}

// User returns a user of specified userID.
func (s *slackService) User(userID string) plugin.User {
	if len(userID) == 0 || userID[0] != 'U' {
		// ignore group
		return nil
	}

	u, err := s.client.GetUserInfo(userID)
	if err != nil {
		s.l.Error("Failed to get user info", slog.String("user_id", userID), slog.Any("err", err))
		return nil
	}

	return &user{
		id:          u.ID,
		name:        u.Name,
		fullName:    u.RealName,
		displayName: u.Profile.DisplayName,
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
