package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"git.wh64.net/muffin/goMuffin/commands"
	"git.wh64.net/muffin/goMuffin/components"
	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/databases"
	"git.wh64.net/muffin/goMuffin/handler"
	"github.com/bwmarrin/discordgo"
)

func main() {
	config := configs.Config

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

	for _, cmd := range commands.Discommand.Commands {
		if cmd.Name == "도움말" {
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

	defer func() {
		dg.Close()
		databases.Client.Disconnect(context.TODO())
	}()

	log.Println("[goMuffin] 봇이 실행되고 있어요. 버전:", configs.MUFFIN_VERSION)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
