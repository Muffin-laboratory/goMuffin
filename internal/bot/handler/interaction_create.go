package handler

import (
	"log/slog"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders"
	"github.com/Muffin-laboratory/goMuffin/internal/bot/loader"
	"github.com/Muffin-laboratory/goMuffin/internal/configs"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/rest"
)

func OnApplicationCommandInteractionCreate(i *events.ApplicationCommandInteractionCreate) {
	if err := loader.GetDiscommand().ChatInputRun(i.Data.CommandName(), i); err != nil {
		slog.Error("error in responding chat input command.", "user_id", i.User().ID, "error", err)
		i.CreateMessage(
			discord.NewMessageCreateBuilder().
				SetComponents(getErrContainer(i.Client().Rest)).
				SetIsComponentsV2(true).
				SetEphemeral(true).
				Build(),
		)
	}
}

func OnComponentInteractionCreate(i *events.ComponentInteractionCreate) {
	if err := loader.GetDiscommand().ComponentRun(i); err != nil {
		slog.Error("error in responding component.", "user_id", i.User().ID, "error", err)
		i.CreateMessage(
			discord.NewMessageCreateBuilder().
				SetComponents(getErrContainer(i.Client().Rest)).
				SetIsComponentsV2(true).
				SetEphemeral(true).
				Build(),
		)
	}
}

func OnModalSubmitInteractionCreate(i *events.ModalSubmitInteractionCreate) {
	if err := loader.GetDiscommand().ModalRun(i); err != nil {
		slog.Error("error in responding modal submit.", "user_id", i.User().ID, "error", err)
		i.CreateMessage(
			discord.NewMessageCreateBuilder().
				SetComponents(getErrContainer(i.Client().Rest)).
				SetIsComponentsV2(true).
				SetEphemeral(true).
				Build(),
		)
	}
}

func OnAutocompleteInteractionCreate(i *events.AutocompleteInteractionCreate) {
	if err := loader.GetDiscommand().ChatInputAutocomplete(i.Data.CommandName, i); err != nil {
		slog.Error("error in responding autocomplete.", "user_id", i.User().ID, "error", err)
		i.Respond(discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().
				SetComponents(getErrContainer(i.Client().Rest)).
				SetIsComponentsV2(true).
				SetEphemeral(true).
				Build(),
		)
	}
}

func getErrContainer(client rest.Rest) discord.ContainerComponent {
	owner, _ := client.GetUser(configs.GetConfig().Bot.OwnerID)
	return builders.MakeErrorContainer("오류가 발생하였어요. 만약 계속 발생한다면, %s으로 연락해주세요.", utils.InlineCode(owner.Username))
}
