package discord

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	discord "github.com/bwmarrin/discordgo"
	"github.com/kechako/gopher-bot/plugin"
	"github.com/kechako/gopher-bot/service"
)

// discordService represents a service for Discord.
type discordService struct {
	session *discord.Session
	ch      chan *service.Event
}

// New returns a new Discord service as service.Service.
func New(token string) (service.Service, error) {
	if token == "" {
		return nil, errors.New("the token is empty")
	}

	session, err := discord.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new Discord session: %w", err)
	}

	s := &discordService{
		session: session,
	}
	s.addHandlers()

	return s, nil
}

// Start implements the service.Service interface.
func (s *discordService) Start(ctx context.Context) (<-chan *service.Event, error) {
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
		// TODO: output error log
	}
}

// Mention implements the service.Service interface.
func (s *discordService) Mention(channelID, userID, text string) {
	user, err := s.session.User(userID)
	if err != nil {
		// TODO: output error log
	}

	text = user.Mention() + text
	_, err = s.session.ChannelMessageSend(channelID, text)
	if err != nil {
		// TODO: output error log
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
		// TODO : output log
		return nil
	}

	return &channel{
		id:   ch.ID,
		name: ch.Name,
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
	s.session.AddHandler(func(session *discord.Session, event *discord.Ready) {
		s.handleReady(event)
	})

	s.session.AddHandler(func(session *discord.Session, event *discord.MessageCreate) {
		s.handleMessageCreate(event)
	})
}

// handleReady handles the Ready event.
func (s *discordService) handleReady(msg *discord.Ready) {
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
