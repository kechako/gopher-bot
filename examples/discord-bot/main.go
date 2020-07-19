package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/kechako/gopher-bot/v2"
	"github.com/kechako/gopher-bot/v2/examples/plugins/echo"
	"github.com/kechako/gopher-bot/v2/service/discord"
	"github.com/kechako/logger"
	"github.com/kechako/sigctx"
)

func main() {
	var token string
	flag.StringVar(&token, "token", "", "Discord bot token.")
	flag.Parse()

	service, err := discord.New(token)
	if err != nil {
		log.Fatal(err)
	}

	b, err := bot.New(service, bot.WithLogger(logger.New()))
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
