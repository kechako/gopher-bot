package main

import (
	"context"
	"flag"
	"log"
	"os"

	bot "github.com/kechako/gopher-bot"
	"github.com/kechako/gopher-bot/examples/plugins/echo"
	"github.com/kechako/gopher-bot/service/discord"
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

	b, err := bot.New(service)
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
