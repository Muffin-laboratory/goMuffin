package loader

import (
	"context"
	"time"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/bwmarrin/discordgo"
)

type CommandFlags uint8

type Command struct {
	*discordgo.ApplicationCommand
	DeferOptions *discordgo.InteractionResponseData
	Flags        CommandFlags
	Deferred     bool
	Run          run
	Autocomplete run
}

const (
	CommandFlagsIsRegistered CommandFlags = 1 << iota
	CommandFlagsIsBlocked
	CommandFlagsIsDeveloperOnlyCommand
)

func (d *Discommand) LoadCommand(c *Command) {
	defer commandMutex.Unlock()
	commandMutex.Lock()
	d.Commands[c.Name] = c
}

func (d *Discommand) ChatInputRun(name string, s *discordgo.Session, inter *discordgo.InteractionCreate) error {
	i := &builders.InteractionCreate{
		InteractionCreate: inter,
		Session:           s,
		Options:           builders.MakeCommandInteractionOptionsMap(inter.ApplicationCommandData().Options),
	}

	i.InteractionCreate.User = builders.GetInteractionUser(inter)

	if command, ok := d.Commands[name]; ok {
		var ctx context.Context
		var cancel context.CancelFunc
		if command.Deferred {
			ctx, cancel = context.WithTimeout(context.Background(), 15*time.Minute)
		} else {
			ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
		}
		i.Ctx = ctx

		defer cancel()

		if command.Deferred {
			if err := i.DeferReply(command.DeferOptions); err != nil {
				return err
			}
		}

		if command.Flags&CommandFlagsIsDeveloperOnlyCommand != 0 && i.User.ID != configs.GetConfig().Bot.OwnerID {
			return builders.NewMessageSender(i).
				AddComponents(builders.MakeDeclineContainer("이 명령어는 개발자 전용 명령어에요.")).
				SetComponentsV2(true).
				SetEphemeral(true).
				Send()
		}

		if command.Flags&CommandFlagsIsRegistered != 0 && !repository.GetDatabase().Users.IsUser(ctx, i.User.ID) {
			return builders.NewMessageSender(i).
				AddComponents(builders.MakeUserIsNotRegisteredErrContainer()).
				SetComponentsV2(true).
				SetEphemeral(true).
				SetReply(true).
				Send()
		}

		blocked, reason := repository.GetDatabase().Users.IsUserBlocked(ctx, i.User.ID)
		if command.Flags&CommandFlagsIsBlocked != 0 && blocked {
			user, _ := s.User(i.User.ID)
			return builders.NewMessageSender(i).
				AddComponents(builders.MakeUserIsBlockedContainer(user.GlobalName, reason)).
				SetComponentsV2(true).
				SetReply(true).
				Send()
		}

		return command.Run(i)

	}
	return nil
}

func (d *Discommand) ChatInputAutocomplete(name string, s *discordgo.Session, inter *discordgo.InteractionCreate) error {
	i := &builders.InteractionCreate{
		InteractionCreate: inter,
		Session:           s,
	}

	i.InteractionCreate.User = builders.GetInteractionUser(inter)

	if command, ok := d.Commands[name]; ok {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		i.Ctx = ctx

		defer cancel()

		return command.Autocomplete(i)
	}

	return nil
}
