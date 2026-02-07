package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
)

func main() {
	err := session.OpenGateway(context.Background())
	if err != nil {
		log.Println("[goMuffin] 봇을 시작할 수 없어요.")
		log.Fatalln(err)
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
	var developerOnlyGuildCmds []discord.ApplicationCommandCreate
	for _, cmd := range loader.GetDiscommand().Commands {
		if cmd.Flags&loader.CommandFlagsIsDeveloperOnlyCommand != 0 {
			developerOnlyGuildCmds = append(developerOnlyGuildCmds, cmd.SlashCommandCreate)
			continue
		}

		globalCmds = append(globalCmds, cmd.SlashCommandCreate)
	}

	_, err = session.Rest.SetGlobalCommands(session.ApplicationID, globalCmds)
	if err != nil {
		log.Println(err)
	}

	if len(developerOnlyGuildCmds) != 0 {
		developerOnlyGuildID := snowflake.MustParse(configs.GetConfig().Command.DeveloperOnlyGuildID)
		_, err = session.Rest.SetGuildCommands(session.ApplicationID, developerOnlyGuildID, developerOnlyGuildCmds)
		if err != nil {
			log.Println(err)
		}
	}

	defer repository.GetDatabase().Disconnect()

	log.Println("[goMuffin] 봇이 실행되고 있어요. 버전:", configs.MuffinVersion)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
