package customid

import "fmt"

const (
	UserSettings               = "/muffin/user/settings"
	UserSettingsChattingMode   = UserSettings + "/chatting_mode"
	UserSettingsReplyUser      = UserSettings + "/reply_user"
	UserSettings12Hours        = UserSettings + "/12hours"
	UserSettingsPrompt         = UserSettings + "/prompt"
	UserSettingsPromptSet      = UserSettings + "/prompt/set"
	UserSettingsChatPerCHannel = UserSettings + "/chat_per_channel"
	UserSettingsSubmit         = UserSettings + "/submit"
	UserSettingsCancel         = UserSettings + "/cancel"
)

func MakeUserSettingsChattingMode(id string) string {
	return fmt.Sprintf("%s/%s", UserSettingsChattingMode, id)
}

func MakeUserSettingsReplyUser(id string) string {
	return fmt.Sprintf("%s/%s", UserSettingsReplyUser, id)
}

func MakeUserSettings12Hours(id string) string {
	return fmt.Sprintf("%s/%s", UserSettings12Hours, id)
}

func MakeUserSettingsPrompt(id string) string {
	return fmt.Sprintf("%s/%s", UserSettingsPrompt, id)
}

func MakeUserSettingsChatPerChannel(id string) string {
	return fmt.Sprintf("%s/%s", UserSettingsChatPerCHannel, id)
}

func MakeUserSettingsSubmit(id string) string {
	return fmt.Sprintf("%s/%s", UserSettingsSubmit, id)
}

func MakeUserSettingsCancel(id string) string {
	return fmt.Sprintf("%s/%s", UserSettingsCancel, id)
}
