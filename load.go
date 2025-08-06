package main

import (
	"git.wh64.net/muffin/goMuffin/commands"
	"git.wh64.net/muffin/goMuffin/commands/dev"
	"git.wh64.net/muffin/goMuffin/components"
	"git.wh64.net/muffin/goMuffin/modals"
)

func init() {
	// General command
	go commands.GetDiscommand().LoadCommand(commands.HelpCommand)
	go commands.GetDiscommand().LoadCommand(commands.DataLengthCommand)
	go commands.GetDiscommand().LoadCommand(commands.InformationCommand)
	go commands.GetDiscommand().LoadCommand(commands.SwitchModeCommand)
	go commands.GetDiscommand().LoadCommand(commands.DeregisterCommand)

	// Chatting command
	go commands.GetDiscommand().LoadCommand(commands.RegisterCommand)
	go commands.GetDiscommand().LoadCommand(commands.LearnCommand)
	go commands.GetDiscommand().LoadCommand(commands.LearnedDataListCommand)
	go commands.GetDiscommand().LoadCommand(commands.DeleteLearnedDataCommand)
	go commands.GetDiscommand().LoadCommand(commands.ChatCommand)

	// Developer only command
	go commands.GetDiscommand().LoadCommand(dev.BlockCommand)
	go commands.GetDiscommand().LoadCommand(dev.ReloadPromptCommand)
	go commands.GetDiscommand().LoadCommand(dev.UnblockCommand)

	// Message component
	go commands.GetDiscommand().LoadComponent(components.DeleteLearnedDataComponent)
	go commands.GetDiscommand().LoadComponent(components.PaginationEmbedComponent)
	go commands.GetDiscommand().LoadComponent(components.RegisterComponent)
	go commands.GetDiscommand().LoadComponent(components.DeregisterComponent)
	go commands.GetDiscommand().LoadComponent(components.SelectChatComponent)
	go commands.GetDiscommand().LoadComponent(components.DeleteChatComponent)

	// Modal component
	go commands.GetDiscommand().LoadModal(modals.PaginationEmbedModal)
}
