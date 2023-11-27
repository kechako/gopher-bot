package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"
	"os/signal"

	bot "github.com/kechako/gopher-bot/v2"
	"github.com/kechako/gopher-bot/v2/examples/plugins/echo"
	"github.com/kechako/gopher-bot/v2/service/discord"
)

func main() {
	var token string
	flag.StringVar(&token, "token", "", "Discord bot token.")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	service, err := discord.New(token, &discord.Config{
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
