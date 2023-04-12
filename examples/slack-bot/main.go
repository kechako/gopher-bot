package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	bot "github.com/kechako/gopher-bot/v2"
	"github.com/kechako/gopher-bot/v2/examples/plugins/echo"
	"github.com/kechako/gopher-bot/v2/service/slack"
	"golang.org/x/exp/slog"
)

func main() {
	var token string
	var appToken string
	flag.StringVar(&token, "token", "", "Slack bot token.")
	flag.StringVar(&appToken, "app-token", "", "Slack bot application token.")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout))

	service, err := slack.New(token, appToken, &slack.Config{
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

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := b.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
