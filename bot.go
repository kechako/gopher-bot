// Package bot implements core functions for a bot service.
package bot

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/kechako/gopher-bot/v2/internal/database"
	"github.com/kechako/gopher-bot/v2/logger"
	"github.com/kechako/gopher-bot/v2/plugin"
	"github.com/kechako/gopher-bot/v2/service"
)

// Bot represents a bot.
type Bot struct {
	service service.Service
	plugins []plugin.Plugin

	db          *database.DB
	databaseDir string

	l logger.Logger

	helloOnce sync.Once
}

// New returns a new *Bot.
func New(s service.Service, opts ...Option) (*Bot, error) {
	bot := &Bot{
		service: s,
	}

	for _, opt := range opts {
		opt(bot)
	}

	err := bot.init()
	if err != nil {
		return nil, err
	}

	return bot, nil
}

func (b *Bot) init() error {
	if b.databaseDir == "" {
		dir, err := ioutil.TempDir(os.TempDir(), "gopher-bot")
		if err != nil {
			return fmt.Errorf("failed to create database dir: %w", err)
		}
		b.databaseDir = dir
	} else {
		if stat, err := os.Stat(b.databaseDir); err != nil {
			if err := os.MkdirAll(b.databaseDir, 0755); err != nil {
				return fmt.Errorf("failed to create database dir: %w", err)
			}
		} else if !stat.IsDir() {
			return fmt.Errorf("%s is not a directory", b.databaseDir)
		}
	}

	path := filepath.Join(b.databaseDir, "gopher-bot.db")

	db, err := database.Open(path)
	if err != nil {
		return err
	}
	b.db = db

	if b.l == nil {
		b.l = logger.NewNop()
	}

	return nil
}

func (b *Bot) Close() error {
	for _, p := range b.plugins {
		if c, ok := p.(io.Closer); ok {
			c.Close()
		}
	}

	return b.db.Close()
}

// AddPlugin adds a plugin to the bot.
func (b *Bot) AddPlugin(p plugin.Plugin) {
	b.plugins = append(b.plugins, p)
}

// Run runs the bot.
func (b *Bot) Run(ctx context.Context) error {
	b.l.Info("Start to run bot service.")

	// set database to context
	ctx = database.ContextWithDB(ctx, b.db)

	// set logger to context
	ctx = logger.ContextWithLogger(ctx, b.l)

	ch, err := b.service.Start(ctx)
	if err != nil {
		return err
	}

	b.loop(ctx, ch)

	return nil
}

func (b *Bot) loop(ctx context.Context, ch <-chan *service.Event) {
	defer b.service.Close()

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case event := <-ch:
			switch event.Type {
			case service.ConnectedEvent:
				if hello := event.GetHello(); hello != nil {
					b.hello(ctx, hello)
				}
			case service.MessageEvent:
				if msg := event.GetMessage(); msg != nil {
					b.doAction(ctx, msg)
				}
			}
		}
	}
}

func (b *Bot) hello(ctx context.Context, hello plugin.Hello) {
	b.helloOnce.Do(func() {
		for _, p := range b.plugins {
			b.callPluginHello(ctx, p, hello)
		}
	})
}

func (b *Bot) callPluginHello(ctx context.Context, plugin plugin.Plugin, hello plugin.Hello) {
	defer func() {
		if err := recover(); err != nil {
			b.l.Errorf("plugin.Hello (%T): %v", plugin, err)
		}
	}()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	ch := make(chan struct{})
	go func() {
		plugin.Hello(ctx, hello)
		close(ch)
	}()

	select {
	case <-ch:
	case <-ctx.Done():
		if err := ctx.Err(); err != nil {
			b.l.Errorf("plugin.Hello (%T): %v", plugin, err)
		}
	}
}

func (b *Bot) doAction(ctx context.Context, msg plugin.Message) {
	if msg.MentionTo(b.service.UserID()) && strings.Contains(msg.Text(), "help") {
		b.postHelp(ctx, msg)
		return
	}

	for _, p := range b.plugins {
		b.callPluginDoAction(ctx, p, msg)
	}
}

func (b *Bot) callPluginDoAction(ctx context.Context, plugin plugin.Plugin, msg plugin.Message) {
	defer func() {
		if err := recover(); err != nil {
			b.l.Errorf("plugin.DoAction (%T): %v", plugin, err)
		}
	}()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	ch := make(chan struct{})
	go func() {
		plugin.DoAction(ctx, msg)
		close(ch)
	}()

	select {
	case <-ch:
	case <-ctx.Done():
		if err := ctx.Err(); err != nil {
			b.l.Errorf("plugin.DoAction (%T): %v", plugin, err)
		}
	}
}

func (b *Bot) postHelp(ctx context.Context, msg plugin.Message) {
	var doc strings.Builder

	for i, p := range b.plugins {
		h := b.callPluginHelp(ctx, p)
		if h == nil {
			continue
		}
		if i > 0 {
			doc.WriteString("\n\n")
		}

		doc.WriteString(h.String())
	}

	escaped := b.service.EscapeHelp(doc.String())
	msg.Post(escaped)
}

func (b *Bot) callPluginHelp(ctx context.Context, p plugin.Plugin) *plugin.Help {
	defer func() {
		if err := recover(); err != nil {
			b.l.Errorf("plugin.Help (%T): %v", p, err)
		}
	}()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	ch := make(chan *plugin.Help, 1)
	go func() {
		ch <- p.Help(ctx)
	}()

	select {
	case help := <-ch:
		return help
	case <-ctx.Done():
		if err := ctx.Err(); err != nil {
			b.l.Errorf("plugin.Help (%T): %v", p, err)
		}
		return nil
	}
}

type Option func(bot *Bot)

func WithDatabaseDir(dir string) Option {
	return func(bot *Bot) {
		bot.databaseDir = dir
	}
}

func WithLogger(l logger.Logger) Option {
	return func(bot *Bot) {
		bot.l = l
	}
}
