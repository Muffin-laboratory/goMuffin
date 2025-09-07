package main

import (
	"git.wh64.net/muffin/goMuffin/chatbot"
	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/handler"
	"github.com/bwmarrin/discordgo"
)

var dg *discordgo.Session

func init() {
	dg, _ = discordgo.New("Bot " + configs.GetConfig().Bot.Token)
	go chatbot.Make(dg)

	// Handler
	go dg.AddHandler(handler.MessageCreate)
	go dg.AddHandler(handler.InteractionCreate)
}
