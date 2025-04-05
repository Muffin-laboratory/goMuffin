package utils

import "github.com/bwmarrin/discordgo"

// InteractionCreate custom data of discordgo.InteractionCreate
type InteractionCreate struct {
	*discordgo.InteractionCreate
	Session *discordgo.Session
	// NOTE: It's only can ApplicationCommand
	Options map[string]*discordgo.ApplicationCommandInteractionDataOption
}

// Reply to this interaction.
func (i *InteractionCreate) Reply(data *discordgo.InteractionResponseData) {
	i.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: data,
	})
}

// GetInteractionOptions to this interaction.
// NOTE: It's only can ApplicationCommand
func GetInteractionOptions(i *discordgo.InteractionCreate) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	optsMap := map[string]*discordgo.ApplicationCommandInteractionDataOption{}
	for _, opt := range i.ApplicationCommandData().Options {
		optsMap[opt.Name] = opt
	}
	return optsMap
}

// DeferReply to this interaction.
func (i *InteractionCreate) DeferReply(ephemeral bool) {
	var flags discordgo.MessageFlags
	if ephemeral {
		flags = discordgo.MessageFlagsEphemeral
	}

	i.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: flags,
		},
	})
}

// DeferUpdate to this interaction.
func (i *InteractionCreate) DeferUpdate() {
	i.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
}

// EditReply to this interaction.
func (i *InteractionCreate) EditReply(data *discordgo.WebhookEdit) {
	i.Session.InteractionResponseEdit(i.Interaction, data)
}
