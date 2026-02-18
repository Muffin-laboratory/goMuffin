package customid

import "fmt"

const (
	UserSettings             = "/muffin/user/settings"
	UserSettingsChattingMode = "/muffin/user/settings/chatting_mode"
	UserSettingsReplyUser    = "/muffin/user/settings/reply_user"
	UserSettings12Hours      = "/muffin/user/settings/12hours"
	UserSettingsPrompt       = "/muffin/user/settings/prompt"
	UserSettingsPromptSet    = "/muffin/user/settings/prompt/set"
	UserSettingsSubmit       = "/muffin/user/settings/submit"
	UserSettingsCancel       = "/muffin/user/settings/cancel"
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

func MakeUserSettingsSubmit(id string) string {
	return fmt.Sprintf("%s/%s", UserSettingsSubmit, id)
}

func MakeUserSettingsCancel(id string) string {
	return fmt.Sprintf("%s/%s", UserSettingsCancel, id)
}
