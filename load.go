package main

import (
	"git.wh64.net/muffin/goMuffin/chatbot"
	"git.wh64.net/muffin/goMuffin/commands"
	"git.wh64.net/muffin/goMuffin/components"
	"git.wh64.net/muffin/goMuffin/configs"
	"git.wh64.net/muffin/goMuffin/handler"
	"git.wh64.net/muffin/goMuffin/modals"
	"github.com/bwmarrin/discordgo"
)

var dg *discordgo.Session

func init() {
	dg, _ = discordgo.New("Bot " + configs.GetConfig().Bot.Token)
	go chatbot.Make(dg)

	// Handler
	go dg.AddHandler(handler.MessageCreate)
	go dg.AddHandler(handler.InteractionCreate)

	// General command
	go commands.GetDiscommand().LoadCommand(commands.HelpCommand)
	go commands.GetDiscommand().LoadCommand(commands.DataLengthCommand)
	go commands.GetDiscommand().LoadCommand(commands.InformationCommand)
	go commands.GetDiscommand().LoadCommand(commands.DeregisterCommand)

	// Chatting command
	go commands.GetDiscommand().LoadCommand(commands.RegisterCommand)
	go commands.GetDiscommand().LoadCommand(commands.LearnCommand)
	go commands.GetDiscommand().LoadCommand(commands.KnowledgeListCommand)
	go commands.GetDiscommand().LoadCommand(commands.DeleteKnowledgeCommand)
	go commands.GetDiscommand().LoadCommand(commands.ChatCommand)

	// Message component
	go commands.GetDiscommand().LoadComponent(components.DeleteKnowledgeComponent)
	go commands.GetDiscommand().LoadComponent(components.PaginationContainerComponent)
	go commands.GetDiscommand().LoadComponent(components.RegisterComponent)
	go commands.GetDiscommand().LoadComponent(components.DeregisterComponent)
	go commands.GetDiscommand().LoadComponent(components.SelectChatComponent)
	go commands.GetDiscommand().LoadComponent(components.DeleteChatComponent)

	// Modal component
	go commands.GetDiscommand().LoadModal(modals.PaginationContainerModal)
}
