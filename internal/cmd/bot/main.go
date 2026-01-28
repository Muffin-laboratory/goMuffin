package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/commands"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func main() {
	err := dg.Open()
	if err != nil {
		log.Println("[goMuffin] 봇을 시작할 수 없어요.")
		log.Fatalln(err)
	}

	defer dg.Close()

	// 봇의 상태메세지 변경
	go func() {
		for {
			dg.UpdateCustomStatus("ㅅ살려주세요..!")
			time.Sleep(time.Minute * 10)
		}
	}()

	var globalCmds []*discordgo.ApplicationCommand
	var developerOnlyGuildCmds []*discordgo.ApplicationCommand
	for _, cmd := range commands.GetDiscommand().Commands {
		if cmd.Flags&commands.CommandFlagsIsDeveloperOnlyCommand != 0 {
			developerOnlyGuildCmds = append(developerOnlyGuildCmds, cmd.ApplicationCommand)
			continue
		}

		globalCmds = append(globalCmds, cmd.ApplicationCommand)
	}

	_, err = dg.ApplicationCommandBulkOverwrite(dg.State.User.ID, "", globalCmds)
	if err != nil {
		log.Println(err)
	}

	if len(developerOnlyGuildCmds) != 0 {
		_, err = dg.ApplicationCommandBulkOverwrite(dg.State.User.ID, configs.GetConfig().Command.DeveloperOnlyGuildID, developerOnlyGuildCmds)
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
