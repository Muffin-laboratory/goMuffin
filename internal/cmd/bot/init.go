package main

import (
	"log"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/chatbot"
	_ "github.com/Muffin-laboratory/goMuffin/internal/bot/commands/dev"
	_ "github.com/Muffin-laboratory/goMuffin/internal/bot/components"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/handler"
	_ "github.com/Muffin-laboratory/goMuffin/internal/bot/modals"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
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
