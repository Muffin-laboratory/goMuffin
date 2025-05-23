package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

type ModalData struct {
	CustomId   string                       `json:"custom_id"`
	Title      string                       `json:"title"`
	Components []discordgo.MessageComponent `json:"components"`
}

// InteractionCreate custom data of discordgo.InteractionCreate
type InteractionCreate struct {
	*discordgo.InteractionCreate
	Session *discordgo.Session
	// NOTE: It's only can ApplicationCommand
	Options  map[string]*discordgo.ApplicationCommandInteractionDataOption
	Deferred bool
	Replied  bool
}

// Reply to this interaction.
func (i *InteractionCreate) Reply(data *discordgo.InteractionResponseData) {
	i.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: data,
	})

	i.Replied = true
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

	i.Deferred = true
}

// DeferUpdate to this interaction.
func (i *InteractionCreate) DeferUpdate() {
	i.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})

	i.Deferred = true
}

// EditReply to this interaction.
func (i *InteractionCreate) EditReply(data *discordgo.WebhookEdit) {
	i.Session.InteractionResponseEdit(i.Interaction, data)

	i.Replied = true
}

// Update to this interaction.
func (i *InteractionCreate) Update(data *discordgo.InteractionResponseData) {
	i.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: data,
	})

	i.Replied = true
}

func (i *InteractionCreate) ShowModal(data *ModalData) error {
	var reqData struct {
		Type discordgo.InteractionResponseType `json:"type"`
		Data ModalData                         `json:"data"`
	}

	reqData.Type = discordgo.InteractionResponseModal
	reqData.Data = *data

	endpoint := discordgo.EndpointInteractionResponse(i.ID, i.Token)
	_, err := i.Session.RequestWithBucketID("POST", endpoint, reqData, endpoint)
	return err
}
