package loader

import "github.com/disgoorg/disgo/discord"

func (d *Discommand) RegisterCommand(command discord.ApplicationCommandCreate) {
	d.commands = append(d.commands, command)
}

func (d *Discommand) Commands() []discord.ApplicationCommandCreate {
	commandsCopy := make([]discord.ApplicationCommandCreate, len(d.commands))
	copy(commandsCopy, d.commands)
	return commandsCopy
}

func (d *Discommand) RegisterDevCommand(command discord.ApplicationCommandCreate) {
	d.devCommands = append(d.devCommands, command)
}

func (d *Discommand) DevCommands() []discord.ApplicationCommandCreate {
	devCommandsCopy := make([]discord.ApplicationCommandCreate, len(d.devCommands))
	copy(devCommandsCopy, d.devCommands)
	return devCommandsCopy
}
