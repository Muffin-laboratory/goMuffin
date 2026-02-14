package loader

import (
	"context"
	"time"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type CommandFlags uint8

type Command struct {
	*discord.SlashCommandCreate
	Flags            CommandFlags
	Deferred         bool
	IsDeferEphemeral bool
	Run              func(ctx context.Context, i *builders.CommandCreate) error
	Autocomplete     func(ctx context.Context, i *events.AutocompleteInteractionCreate) error
}

const (
	CommandFlagsIsRegistered CommandFlags = 1 << iota
	CommandFlagsIsBlocked
	CommandFlagsIsDeveloperOnlyCommand
)

func (d *Discommand) LoadCommand(c *Command) {
	defer commandMutex.Unlock()
	commandMutex.Lock()
	d.OldCommands[c.Name] = c
}

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

func (d *Discommand) ChatInputRun(name string, inter *events.ApplicationCommandInteractionCreate) error {
	if command, ok := d.OldCommands[name]; ok {
		var ctx context.Context
		var cancel context.CancelFunc
		if command.Deferred {
			ctx, cancel = context.WithTimeout(context.Background(), 1*time.Minute)
		} else {
			ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
		}
		defer cancel()

		i := &builders.CommandCreate{
			ApplicationCommandInteractionCreate: inter,
			Responded:                           command.Deferred,
		}

		if command.Deferred {
			if err := i.DeferCreateMessage(command.IsDeferEphemeral); err != nil {
				return err
			}
		}

		isOwner := i.User().ID == configs.GetConfig().Bot.OwnerID
		if command.Flags&CommandFlagsIsDeveloperOnlyCommand != 0 && !isOwner {
			return builders.NewMessageSender(i).
				AddComponents(builders.MakeDeclineContainer("이 명령어는 개발자 전용 명령어에요.")).
				SetComponentsV2(true).
				SetEphemeral(true).
				Send()
		}

		if command.Flags&CommandFlagsIsRegistered != 0 && !repository.GetDatabase().Users.IsUser(ctx, int64(i.User().ID)) {
			return builders.NewMessageSender(i).
				AddComponents(builders.MakeUserIsNotRegisteredErrContainer()).
				SetComponentsV2(true).
				SetEphemeral(true).
				SetReply(true).
				Send()
		}

		blocked, reason := repository.GetDatabase().Users.IsUserBlocked(ctx, int64(i.User().ID))
		if command.Flags&CommandFlagsIsBlocked != 0 && blocked {
			return builders.NewMessageSender(i).
				AddComponents(builders.MakeUserIsBlockedContainer(*i.User().GlobalName, reason)).
				SetComponentsV2(true).
				SetReply(true).
				Send()
		}

		return command.Run(ctx, i)
	}
	return nil
}

func (d *Discommand) ChatInputAutocomplete(name string, i *events.AutocompleteInteractionCreate) error {
	if command, ok := d.OldCommands[name]; ok {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		return command.Autocomplete(ctx, i)
	}

	return nil
}
