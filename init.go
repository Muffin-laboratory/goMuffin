package main

import (
	"log"

	"git.wh64.net/muffin/goMuffin/chatbot"
	_ "git.wh64.net/muffin/goMuffin/commands/dev"
	_ "git.wh64.net/muffin/goMuffin/components"
	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/handler"
	_ "git.wh64.net/muffin/goMuffin/modals"
	"github.com/bwmarrin/discordgo"
)

var dg *discordgo.Session

func init() {
	dg, _ = discordgo.New("Bot " + configs.GetConfig().Bot.Token)
	err := chatbot.Make(dg)
	if err != nil {
		log.Fatalln(err)
	}

	// Handler
	go dg.AddHandler(handler.MessageCreate)
	go dg.AddHandler(handler.InteractionCreate)
}
