package main

import (
	"log"

	"github.com/Muffin-laboratory/goMuffin/chatbot"
	_ "github.com/Muffin-laboratory/goMuffin/commands/dev"
	_ "github.com/Muffin-laboratory/goMuffin/components"
	"github.com/Muffin-laboratory/goMuffin/configs"
	"github.com/Muffin-laboratory/goMuffin/handler"
	_ "github.com/Muffin-laboratory/goMuffin/modals"
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
