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

type InteractionEdit struct {
	Content         *string                           `json:"content,omitempty"`
	Components      *[]discordgo.MessageComponent     `json:"components,omitempty"`
	Embeds          *[]*discordgo.MessageEmbed        `json:"embeds,omitempty"`
	Flags           *discordgo.MessageFlags           `json:"flags,omitempty"`
	Attachments     *[]*discordgo.MessageAttachment   `json:"attachments,omitempty"`
	AllowedMentions *discordgo.MessageAllowedMentions `json:"allowed_mentions,omitempty"`
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
func (i *InteractionCreate) EditReply(data *InteractionEdit) error {
	endpoint := discordgo.EndpointWebhookMessage(i.AppID, i.Token, "@original")

	_, err := i.Session.RequestWithBucketID("PATCH", endpoint, *data, discordgo.EndpointWebhookToken("", ""))

	i.Replied = true
	return err
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
