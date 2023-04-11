// Package discord provides a bot service of Discord.
package discord

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	discord "github.com/bwmarrin/discordgo"
	"github.com/kechako/gopher-bot/v2/plugin"
	"github.com/kechako/gopher-bot/v2/service"
	"golang.org/x/exp/slog"
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
	return slog.New(slog.NewTextHandler(os.Stdout))
}

// discordService represents a service for Discord.
type discordService struct {
	session *discord.Session
	ch      chan *service.Event
	l       *slog.Logger
}

// New returns a new Discord service as service.Service.
func New(token string, cfg *Config) (service.Service, error) {
	if token == "" {
		return nil, errors.New("the token is empty")
	}

	session, err := discord.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new Discord session: %w", err)
	}

	s := &discordService{
		session: session,
		l:       cfg.logger(),
	}
	s.addHandlers()

	return s, nil
}

// Start implements the service.Service interface.
func (s *discordService) Start(ctx context.Context) (<-chan *service.Event, error) {
	s.l.Info("Start Discord bot service")

	s.ch = make(chan *service.Event)

	if err := s.session.Open(); err != nil {
		close(s.ch)
		return nil, err
	}

	return s.ch, nil
}

// Start implements the service.Service interface.
func (s *discordService) Close() error {
	return s.session.Close()
}

// UserID implements the service.Service interface.
func (s *discordService) UserID() string {
	return s.session.State.User.ID
}

// Post implements the service.Service interface.
func (s *discordService) Post(channelID, text string) {
	_, err := s.session.ChannelMessageSend(channelID, text)
	if err != nil {
		s.l.Error("Failed to post message to %s", channelID)
	}
}

// Mention implements the service.Service interface.
func (s *discordService) Mention(channelID, userID, text string) {
	user, err := s.session.User(userID)
	if err != nil {
		s.l.Error("Failed to get user info : %s", userID)
	}

	text = user.Mention() + text
	_, err = s.session.ChannelMessageSend(channelID, text)
	if err != nil {
		s.l.Error("Failed to post mention message to %s", channelID)
	}
}

// ProcessCommmand processes the specified command on the channel.
func (s *discordService) ProcessCommand(channelID string, command string) {
	go func() {
		s.ch <- &service.Event{
			Type: service.MessageEvent,
			Data: newCommandMessage(s, channelID, command),
		}
	}()
}

// Channel returns a channel of specified channelID.
func (s *discordService) Channel(channelID string) plugin.Channel {
	ch, err := s.session.Channel(channelID)
	if err != nil {
		s.l.Error("Failed to get channel info : %s", channelID)
		return nil
	}

	return &channel{
		id:   ch.ID,
		name: ch.Name,
	}
}

// User returns a user of specified userID.
func (s *discordService) User(userID string) plugin.User {
	u, err := s.session.User(userID)
	if err != nil {
		s.l.Error("Failed to get user info : %s", userID)
		return nil
	}

	return &user{
		id:   u.ID,
		name: u.Username,
	}
}

// EscapeHelp implements the service.Service interface.
func (s *discordService) EscapeHelp(help string) string {
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

// addHandlers adds discord event handlers.
func (s *discordService) addHandlers() {
	s.session.AddHandler(func(session *discord.Session, event *discord.Connect) {
		s.handleConnect(event)
	})

	s.session.AddHandler(func(session *discord.Session, event *discord.Disconnect) {
		s.handleDisconnect(event)
	})

	s.session.AddHandler(func(session *discord.Session, event *discord.Ready) {
		s.handleReady(event)
	})

	s.session.AddHandler(func(session *discord.Session, event *discord.MessageCreate) {
		s.handleMessageCreate(event)
	})
}

// handleConnect handles the Connect event.
func (s *discordService) handleConnect(msg *discord.Connect) {
	s.l.Info("Discord session is connected")
}

// handleDisconnect handles the Disconnect event.
func (s *discordService) handleDisconnect(msg *discord.Disconnect) {
	s.l.Info("Discord session is disconnected")
}

// handleReady handles the Ready event.
func (s *discordService) handleReady(msg *discord.Ready) {
	s.l.Info("Discord session is ready")

	s.ch <- &service.Event{
		Type: service.ConnectedEvent,
		Data: newHello(s),
	}
}

// handleMessageCreate handles the MessageCreate event.
func (s *discordService) handleMessageCreate(msg *discord.MessageCreate) {
	if msg.Author.ID == s.UserID() {
		// bot message
		return
	}

	s.ch <- &service.Event{
		Type: service.MessageEvent,
		Data: newMessage(s, msg.Message),
	}
}
