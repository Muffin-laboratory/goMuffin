package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"git.wh64.net/muffin/goMuffin/commands"
	"git.wh64.net/muffin/goMuffin/components"
	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/handler"
	"git.wh64.net/muffin/goMuffin/scripts"
	"github.com/bwmarrin/discordgo"
)

func main() {
	config := configs.Config

	if len(os.Args) > 1 {
		switch strings.ToLower(os.Args[1]) {
		case "dbmigrate":
			scripts.DBMigrate()
		case "deleteallcommands":
			scripts.DeleteAllCommands()
		default:
			log.Fatalln(fmt.Errorf("[goMuffin] 명령어 인자에는 dbmigrate나 deleteallcommands만 올 수 있어요"))
		}
		return
	}

	dg, err := discordgo.New("Bot " + config.Bot.Token)
	if err != nil {
		log.Println("[goMuffin] 봇의 세션을 만들수가 없어요.")
		log.Fatalln(err)
	}

	go commands.Discommand.LoadCommand(commands.HelpCommand)
	go commands.Discommand.LoadCommand(commands.DataLengthCommand)
	go commands.Discommand.LoadCommand(commands.LearnCommand)
	go commands.Discommand.LoadCommand(commands.LearnedDataListCommand)
	go commands.Discommand.LoadCommand(commands.InformationCommand)
	go commands.Discommand.LoadCommand(commands.DeleteLearnedDataCommand)

	go commands.Discommand.LoadComponent(components.DeleteLearnedDataComponent)

	go dg.AddHandler(handler.MessageCreate)
	go dg.AddHandler(handler.InteractionCreate)

	dg.Open()
	defer dg.Close()

	go func() {
		for {
			dg.UpdateCustomStatus("ㅅ살려주세요..!")
			time.Sleep(time.Minute * 10)
		}
	}()

	for _, cmd := range commands.Discommand.Commands {
		if cmd.Name == commands.HelpCommand.Name {
			// 극한의 성능 똥망 코드 탄생!
			// 무려 똑같은 걸 반복해서 돌리는!
			for _, a := range commands.Discommand.Commands {
				cmd.Options[0].Choices = append(cmd.Options[0].Choices, &discordgo.ApplicationCommandOptionChoice{
					Name:  a.Name,
					Value: a.Name,
				})
			}
		}

		go dg.ApplicationCommandCreate(dg.State.User.ID, "", cmd.ApplicationCommand)
	}

	defer databases.Client.Disconnect(context.TODO())

	log.Println("[goMuffin] 봇이 실행되고 있어요. 버전:", configs.MUFFIN_VERSION)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
