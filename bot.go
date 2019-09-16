package bot

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/kechako/gopher-bot/internal/store"
	"github.com/kechako/gopher-bot/plugin"
	"github.com/kechako/gopher-bot/service"
)

// Bot represents a bot.
type Bot struct {
	service service.Service
	plugins []plugin.Plugin

	store    *store.Store
	storeDir string

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
	if b.storeDir == "" {
		dir, err := ioutil.TempDir(os.TempDir(), "gopher-bot")
		if err != nil {
			return fmt.Errorf("failed to create database dir: %w", err)
		}
		b.storeDir = dir
	} else {
		if stat, err := os.Stat(b.storeDir); err != nil {
			if err := os.MkdirAll(b.storeDir, 0755); err != nil {
				return fmt.Errorf("failed to create database dir: %w", err)
			}
		} else if !stat.IsDir() {
			return fmt.Errorf("%s is not a directory", b.storeDir)
		}
	}

	store, err := store.New(b.storeDir)
	if err != nil {
		return err
	}
	b.store = store

	return nil
}

func (b *Bot) Close() error {
	for _, p := range b.plugins {
		if c, ok := p.(io.Closer); ok {
			c.Close()
		}
	}

	return b.store.Close()
}

// AddPlugin adds a plugin to the bot.
func (b *Bot) AddPlugin(p plugin.Plugin) {
	b.plugins = append(b.plugins, p)
}

// Run runs the bot.
func (b *Bot) Run(ctx context.Context) error {
	// set database store to context
	ctx = store.ContextWithStore(ctx, b.store)

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
			callPluginHello(ctx, p, hello)
		}
	})
}

func callPluginHello(ctx context.Context, plugin plugin.Plugin, hello plugin.Hello) {
	defer func() {
		if err := recover(); err != nil {
			// TODO: output error log
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
			// TODO: output error log
		}
	}
}

func (b *Bot) doAction(ctx context.Context, msg plugin.Message) {
	if msg.MentionTo(b.service.UserID()) && strings.Contains(msg.Text(), "help") {
		b.postHelp(ctx, msg)
		return
	}

	for _, p := range b.plugins {
		callPluginDoAction(ctx, p, msg)
	}
}

func callPluginDoAction(ctx context.Context, plugin plugin.Plugin, msg plugin.Message) {
	defer func() {
		if err := recover(); err != nil {
			// TODO: output error log
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
			// TODO: output error log
		}
	}
}

func (b *Bot) postHelp(ctx context.Context, msg plugin.Message) {
	var doc strings.Builder

	for i, p := range b.plugins {
		h := callPluginHelp(ctx, p)
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

func callPluginHelp(ctx context.Context, p plugin.Plugin) *plugin.Help {
	defer func() {
		if err := recover(); err != nil {
			// TODO: output error log
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
			// TODO: output error log
		}
		return nil
	}
}

type Option func(bot *Bot)

func WithStoreDir(dir string) Option {
	return func(bot *Bot) {
		bot.storeDir = dir
	}
}
