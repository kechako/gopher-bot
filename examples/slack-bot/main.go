package main

import (
	"context"
	"flag"
	"log"
	"os"

	bot "github.com/kechako/gopher-bot/v2"
	"github.com/kechako/gopher-bot/v2/examples/plugins/echo"
	"github.com/kechako/gopher-bot/v2/service/slack"
	"github.com/kechako/sigctx"
	"golang.org/x/exp/slog"
)

func main() {
	var token string
	flag.StringVar(&token, "token", "", "Slack bot token.")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout))

	service, err := slack.New(token, &slack.Config{
		Logger: logger,
	})
	if err != nil {
		log.Fatal(err)
	}

	b, err := bot.New(service, bot.WithLogger(logger))
	if err != nil {
		log.Fatal(err)
	}
	defer b.Close()

	b.AddPlugin(echo.New())

	ctx, cancel := sigctx.WithCancelBySignal(context.Background(), os.Interrupt)
	defer cancel()

	if err := b.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
