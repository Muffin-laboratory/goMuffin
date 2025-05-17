package utils

import "github.com/bwmarrin/discordgo"

type MessageCreate struct {
	*discordgo.MessageCreate
	Session *discordgo.Session
}

type MessageSender struct {
	Embeds     []*discordgo.MessageEmbed
	Content    string
	Components []discordgo.MessageComponent
	Ephemeral  bool
	Reply      bool
	m          any
}

func NewMessageSender(m any) *MessageSender {
	return &MessageSender{m: m}
}

func (s *MessageSender) AddEmbed(embed *discordgo.MessageEmbed) *MessageSender {
	s.Embeds = append(s.Embeds, embed)
	return s
}

func (s *MessageSender) SetEmbeds(embeds []*discordgo.MessageEmbed) *MessageSender {
	s.Embeds = embeds
	return s
}

func (s *MessageSender) AddComponent(cmp discordgo.MessageComponent) *MessageSender {
	s.Components = append(s.Components, cmp)
	return s
}

func (s *MessageSender) SetComponents(components []discordgo.MessageComponent) *MessageSender {
	s.Components = components
	return s
}

func (s *MessageSender) SetContent(content string) *MessageSender {
	s.Content = content
	return s
}

func (s *MessageSender) SetEphemeral(ephemeral bool) *MessageSender {
	s.Ephemeral = ephemeral
	return s
}

func (s *MessageSender) SetReply(reply bool) *MessageSender {
	s.Reply = reply
	return s
}

func (s *MessageSender) Send() {
	switch m := s.m.(type) {
	case *MessageCreate:
		var reference *discordgo.MessageReference = nil

		if s.Reply {
			reference = m.Reference()
		}

		m.Session.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content:    s.Content,
			Embeds:     s.Embeds,
			Components: s.Components,
			Reference:  reference,
		})
		return
	case *InteractionCreate:
		var flags discordgo.MessageFlags

		if s.Ephemeral {
			flags = discordgo.MessageFlagsEphemeral
		}

		if m.Replied || m.Deferred {
			m.EditReply(&discordgo.WebhookEdit{
				Content:    &s.Content,
				Embeds:     &s.Embeds,
				Components: &s.Components,
			})
			return
		}

		m.Reply(&discordgo.InteractionResponseData{
			Content:    s.Content,
			Embeds:     s.Embeds,
			Components: s.Components,
			Flags:      flags,
		})
	}
}
