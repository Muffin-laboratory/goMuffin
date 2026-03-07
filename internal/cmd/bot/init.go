package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

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
var logFile *os.File

func init() {
	var err error

	logger := newLogger()
	slog.SetDefault(logger)

	session, err = disgo.New(configs.Configs().Bot.Token,
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

func newLogger() *slog.Logger {
	loggerConfig := configs.Configs().Logger
	handlerOptions := &slog.HandlerOptions{
		Level: loggerConfig.Level,
	}

	var handler slog.Handler = slog.NewTextHandler(os.Stdout, handlerOptions)

	if loggerConfig.WriteFile {
		now := time.Now()
		filename := fmt.Sprintf("muffin_log_%d%d%d_%d%d%d.log", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())

		err := os.Mkdir("logs", 0777)
		if err != nil && !os.IsExist(err) {
			goto HandleErr
		}

		logFile, err = os.Create("logs/" + filename)
	HandleErr:
		if err != nil {
			slog.Error("error in creating log file", "error", err)
		} else {
			handler = slog.NewMultiHandler(handler, slog.NewJSONHandler(logFile, handlerOptions))
		}
	}

	logger := slog.New(handler)
	return logger
}
