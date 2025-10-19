package builders

import (
	"github.com/bwmarrin/discordgo"
)

type MessageCreate struct {
	*discordgo.MessageCreate
	Session *discordgo.Session
}

type MessageSender struct {
	Embeds          []*discordgo.MessageEmbed
	Content         string
	Components      []discordgo.MessageComponent
	Ephemeral       bool
	Reply           bool
	ComponentsV2    bool
	AllowedMentions *discordgo.MessageAllowedMentions
	m               any
}

func NewMessageSender(m any) *MessageSender {
	return &MessageSender{m: m}
}

func (s *MessageSender) AddEmbeds(embeds ...*discordgo.MessageEmbed) *MessageSender {
	s.Embeds = append(s.Embeds, embeds...)
	return s
}

func (s *MessageSender) AddComponents(components ...discordgo.MessageComponent) *MessageSender {
	s.Components = append(s.Components, components...)
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

func (s *MessageSender) SetAllowedMentions(allowedMentions discordgo.MessageAllowedMentions) *MessageSender {
	s.AllowedMentions = &allowedMentions
	return s
}

func (s *MessageSender) SetComponentsV2(componentsV2 bool) *MessageSender {
	s.ComponentsV2 = componentsV2
	return s
}

func (s *MessageSender) Send() error {
	var flags discordgo.MessageFlags

	if s.ComponentsV2 {
		flags |= discordgo.MessageFlagsIsComponentsV2
	}

	switch m := s.m.(type) {
	case *MessageCreate:
		var reference *discordgo.MessageReference = nil

		if s.Reply {
			reference = m.Reference()
		}

		_, err := m.Session.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content:         s.Content,
			Embeds:          s.Embeds,
			Components:      s.Components,
			AllowedMentions: s.AllowedMentions,
			Flags:           flags,
			Reference:       reference,
		})
		return err
	case *InteractionCreate:
		if s.Ephemeral {
			flags |= discordgo.MessageFlagsEphemeral
		}

		if m.Replied || m.Deferred {
			return m.EditReply(&InteractionEdit{
				Content:    &s.Content,
				Embeds:     &s.Embeds,
				Components: &s.Components,
				Flags:      &flags,
			})
		}

		return m.Reply(&discordgo.InteractionResponseData{
			Content:    s.Content,
			Embeds:     s.Embeds,
			Components: s.Components,
			Flags:      flags,
		})
	}
	return nil
}
