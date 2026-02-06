package main

import (
	"log"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/chatbot"
	_ "github.com/Muffin-laboratory/goMuffin/internal/bot/commands/dev"
	_ "github.com/Muffin-laboratory/goMuffin/internal/bot/components"
	_ "github.com/Muffin-laboratory/goMuffin/internal/bot/modals"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/gateway"
)

var session *bot.Client

func init() {
	var err error
	session, err = disgo.New(configs.GetConfig().Bot.Token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuilds,
				gateway.IntentGuildMessages,
				gateway.IntentMessageContent,
			),
		),
	)
	if err != nil {
		log.Fatalln(err)
	}

	if err = chatbot.Make(session); err != nil {
		log.Fatalln(err)
	}
}
