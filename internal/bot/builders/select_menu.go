package builders

import "github.com/bwmarrin/discordgo"

// A SelectMenu is select menu
type SelectMenu struct {
	selectMenu *discordgo.SelectMenu
}

func selectMenuBuilder(menuType discordgo.SelectMenuType) *SelectMenu {
	return &SelectMenu{
		selectMenu: &discordgo.SelectMenu{
			MenuType: menuType,
		},
	}
}

// SetCustomID sets custom id
func (m *SelectMenu) SetCustomID(customID string) *SelectMenu {
	m.selectMenu.CustomID = customID
	return m
}

// SetPlaceholder sets placeholder
func (m *SelectMenu) SetPlaceholder(placeholder string) *SelectMenu {
	m.selectMenu.Placeholder = placeholder
	return m
}

// SetMinValues sets minium values
func (m *SelectMenu) SetMinValues(minValues int) *SelectMenu {
	m.selectMenu.MinValues = &minValues
	return m
}

// SetMaxValues sets maximum values
func (m *SelectMenu) SetMaxValues(maxValues int) *SelectMenu {
	m.selectMenu.MaxValues = maxValues
	return m
}

// SetDisabled sets disabled
func (m *SelectMenu) SetDisabled(disabled bool) *SelectMenu {
	m.selectMenu.Disabled = disabled
	return m
}

// Build returns to discordgo.SelectMenu(discordgo.MessageComponent)
func (m *SelectMenu) Build() discordgo.MessageComponent {
	return m.selectMenu
}

// A StringSelectMenu is select menu for string
type StringSelectMenu struct {
	*SelectMenu
}

// StringSelectMenuBuilder creates a new string SelectMenu
func StringSelectMenuBuilder() *StringSelectMenu {
	return &StringSelectMenu{
		SelectMenu: selectMenuBuilder(discordgo.StringSelectMenu),
	}
}

// AddOptions adds options
func (m *StringSelectMenu) AddOptions(options ...discordgo.SelectMenuOption) *StringSelectMenu {
	m.selectMenu.Options = append(m.selectMenu.Options, options...)
	return m
}

// A UserSelectMenu is select menu for user
type UserSelectMenu struct {
	*SelectMenu
}

// UserSelectMenuBuilder creates a new user SelectMenu
func UserSelectMenuBuilder() *UserSelectMenu {
	return &UserSelectMenu{
		SelectMenu: selectMenuBuilder(discordgo.UserSelectMenu),
	}
}

// AddDefaultValues adds default values
func (m *UserSelectMenu) AddDefaultValues(defaultValues ...discordgo.SelectMenuDefaultValue) *UserSelectMenu {
	m.selectMenu.DefaultValues = append(m.selectMenu.DefaultValues, defaultValues...)
	return m
}

// A RoleSelectMenu is select menu for role
type RoleSelectMenu struct {
	*SelectMenu
}

// RoleSelectMenuBuilder creates a new role SelectMenu
func RoleSelectMenuBuilder() *RoleSelectMenu {
	return &RoleSelectMenu{
		SelectMenu: selectMenuBuilder(discordgo.RoleSelectMenu),
	}
}

// AddDefaultValues adds default values
func (m *RoleSelectMenu) AddDefaultValues(defaultValues ...discordgo.SelectMenuDefaultValue) *RoleSelectMenu {
	m.selectMenu.DefaultValues = append(m.selectMenu.DefaultValues, defaultValues...)
	return m
}

// A MentionableSelectMenu is select menu for mentionable
type MentionableSelectMenu struct {
	*SelectMenu
}

// MentionableSelectMenuBuilder creates a new mentionable SelectMenu
func MentionableSelectMenuBuilder() *MentionableSelectMenu {
	return &MentionableSelectMenu{
		SelectMenu: selectMenuBuilder(discordgo.MentionableSelectMenu),
	}
}

// AddDefaultValues adds default values
func (m *MentionableSelectMenu) AddDefaultValues(defaultValues ...discordgo.SelectMenuDefaultValue) *MentionableSelectMenu {
	m.selectMenu.DefaultValues = append(m.selectMenu.DefaultValues, defaultValues...)
	return m
}

// A ChannelSelectMenu is select menu for channel
type ChannelSelectMenu struct {
	*SelectMenu
}

// ChannelSelectMenuBuilder creates a new channel SelectMenu
func ChannelSelectMenuBuilder() *ChannelSelectMenu {
	return &ChannelSelectMenu{
		SelectMenu: selectMenuBuilder(discordgo.ChannelSelectMenu),
	}
}

// AddDefaultValues adds default values
func (m *ChannelSelectMenu) AddDefaultValues(defaultValues ...discordgo.SelectMenuDefaultValue) *ChannelSelectMenu {
	m.selectMenu.DefaultValues = append(m.selectMenu.DefaultValues, defaultValues...)
	return m
}

// AddChannelTypes adds channel types
func (m *ChannelSelectMenu) AddChannelTypes(channelTypes ...discordgo.ChannelType) *ChannelSelectMenu {
	m.selectMenu.ChannelTypes = append(m.selectMenu.ChannelTypes, channelTypes...)
	return m
}
