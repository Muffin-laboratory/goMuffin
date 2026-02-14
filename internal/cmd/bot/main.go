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
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
)

func main() {
	err := session.OpenGateway(context.Background())
	if err != nil {
		slog.Error("[Fatal] failed to start bot.", "error", err)
		os.Exit(1)
	}

	defer session.Close(context.Background())

	// 봇의 상태메세지 변경
	go func() {
		for {
			session.SetPresence(context.Background(), gateway.WithCustomActivity("ㅅ살려주세요..!"))
			time.Sleep(time.Minute * 10)
		}
	}()

	var globalCmds []discord.ApplicationCommandCreate
	for _, cmd := range loader.GetDiscommand().OldCommands {
		globalCmds = append(globalCmds, cmd.SlashCommandCreate)
	}

	_, err = session.Rest.SetGlobalCommands(session.ApplicationID, globalCmds)
	if err != nil {
		slog.Error("error in set global commands.", "error", err)
	}

	developerOnlyGuildCmds := loader.GetDiscommand().DevCommands()
	if len(developerOnlyGuildCmds) != 0 {
		developerOnlyGuildID := configs.GetConfig().Command.DeveloperOnlyGuildID
		_, err = session.Rest.SetGuildCommands(session.ApplicationID, developerOnlyGuildID, developerOnlyGuildCmds)
		if err != nil {
			slog.Error("error in set developer only commands.", "error", err)
		}
	}

	defer repository.GetDatabase().Disconnect()

	slog.Info("bot is running. press ctrl+C to exit program.", "version", configs.MuffinVersion)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
