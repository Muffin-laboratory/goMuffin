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
	configs.TimeStart()
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

	defer func() {
		if err := databases.Client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	log.Println("[goMuffin] 봇이 실행되고 있어요. 버전:", configs.MUFFIN_VERSION)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}
