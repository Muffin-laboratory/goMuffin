package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/disgoorg/disgo/gateway"
)

func main() {
	defer func() {
		err := logFile.Close()
		if err != nil {
			slog.Error("error while closing log file", "error", err)
		}

		err = repository.GetDatabase().Disconnect()
		if err != nil {
			slog.Error("error while closing database", "error", err)
		}
	}()

	err := session.OpenGateway(context.Background())
	if err != nil {
		slog.Error("[Fatal] failed to start bot.", "error", err)
		os.Exit(1)
	}

	defer session.Close(context.Background())

	// 봇의 상태메세지 변경
	go func() {
		for {
			err := session.SetPresence(context.Background(), gateway.WithCustomActivity("ㅅ살려주세요..!"))
			if err != nil {
				slog.Error("failed to change presence", "error", err)
			}
			time.Sleep(time.Minute * 10)
		}
	}()

	globalCmds := loader.GetDiscommand().Commands()
	_, err = session.Rest.SetGlobalCommands(session.ApplicationID, globalCmds)
	if err != nil {
		slog.Error("error while setting global commands.", "error", err)
	}

	developerOnlyGuildCmds := loader.GetDiscommand().DevCommands()
	if len(developerOnlyGuildCmds) != 0 {
		developerOnlyGuildID := configs.Configs().Command.DeveloperOnlyGuildID
		_, err = session.Rest.SetGuildCommands(session.ApplicationID, developerOnlyGuildID, developerOnlyGuildCmds)
		if err != nil {
			slog.Error("error while setting developer only commands.", "error", err)
		}
	}

	slog.Info("bot is running. press ctrl+C to exit program.", "version", configs.MuffinVersion)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
