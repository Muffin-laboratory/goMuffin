package main

import (
	"log"

	"git.wh64.net/muffin/goMuffin/chatbot"
	_ "git.wh64.net/muffin/goMuffin/components"
	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/handler"
	_ "git.wh64.net/muffin/goMuffin/modals"
	"git.wh64.net/muffin/goMuffin/routes"
	"github.com/bwmarrin/discordgo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var dg *discordgo.Session
var server *echo.Echo

func init() {
	dg, _ = discordgo.New("Bot " + configs.GetConfig().Bot.Token)
	err := chatbot.Make(dg)
	if err != nil {
		log.Fatalln(err)
	}

	// Handler
	go dg.AddHandler(handler.MessageCreate)
	go dg.AddHandler(handler.InteractionCreate)

	server = echo.New()

	server.Use(middleware.Recover())
	server.Use(middleware.Logger())

	server.GET("/", routes.Ping)
}
