package main

import (
	"log/slog"
	"os"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/chatbot"
	_ "github.com/Muffin-laboratory/goMuffin/internal/bot/handler"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/handler/events"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/gateway"
)

var session *bot.Client

func init() {
	var err error

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	session, err = disgo.New(configs.GetConfig().Bot.Token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuilds,
				gateway.IntentGuildMessages,
				gateway.IntentMessageContent,
			),
		),
		bot.WithEventListenerFunc(events.OnMessageCreate),
		bot.WithLogger(logger),
		bot.WithEventListeners(loader.GetDiscommand().Router()),
	)
	if err != nil {
		slog.Error("[Fatal] failed to create session.", "error", err)
	}

	if err = chatbot.Make(session); err != nil {
		slog.Error("[Fatal] failed to create chatbot.", "error", err)
	}
}
