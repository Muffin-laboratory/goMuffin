package repository

import (
	"fmt"
	"math/rand"

	"git.wh64.net/muffin/goMuffin/builders"
	"git.wh64.net/muffin/goMuffin/utils"
	"github.com/bwmarrin/discordgo"
)

type UserSettings struct {
	ID                        string
	user                      *discordgo.User
	ChattingMode              ChattingMode
	ReplyUser                 bool
	CreateNewChatAfter12Hours bool
}

var userSettings = make(map[string]*UserSettings)

func NewUserSettings(user *discordgo.User) (*UserSettings, error) {
	dbUser, err := GetDatabase().Users.Get(user.ID)
	if err != nil {
		return nil, err
	}

	id := fmt.Sprintf("%s/%d", user.ID, rand.Intn(100))
	s := UserSettings{
		ID:                        id,
		user:                      user,
		ChattingMode:              dbUser.ChattingMode,
		ReplyUser:                 dbUser.ReplyUser,
		CreateNewChatAfter12Hours: dbUser.CreateNewChatAfter12Hours,
	}

	userSettings[id] = &s

	return &s, nil
}

func GetUserSettings(id string) *UserSettings {
	return userSettings[id]
}

func (s *UserSettings) MakeContainer() *builders.Container {
	return builders.ContainerBuilder().
		AddComponents(
			builders.SectionBuilder().
				SetAccessory(builders.ThumbnailBuilder(s.user.AvatarURL("512"))).
				AddText("### 채팅 설정\n- 설정을 개인화 해주세요.").
				AddText("- **모드**\n> 해당 봇이 답을 하는 방식이에요. (AI <-> 일반)").
				AddText("- **답장 멘션**\n> 해당 봇이 답할 때 답장 멘션 여부를 정해요."),
		).
		AddText("- **12 시간 후 새로운 채팅**\n> 해당 봇과 대화하고 12시간 뒤에 새로운 대화를 시작할지 여부를 정해요.").
		AddComponents(
			builders.ActionsRowBuilder(
				builders.ButtonBuilder().
					SetStyle(discordgo.PrimaryButton).
					SetLabel(fmt.Sprintf("모드: %s", ModeString(s.ChattingMode))).
					SetCustomID(utils.MakeUserSettingsChattingMode(s.ID)),
				builders.ButtonBuilder().
					SetStyle(utils.GetStyleFromBool(s.ReplyUser)).
					SetLabel(fmt.Sprintf("답장 멘션: %s", utils.BoolToString(s.ReplyUser))).
					SetCustomID(utils.MakeUserSettingsReplyUser(s.ID)),
				builders.ButtonBuilder().
					SetStyle(utils.GetStyleFromBool(s.CreateNewChatAfter12Hours)).
					SetLabel(fmt.Sprintf("12 시간 후 새로운 채팅: %s", utils.BoolToString(s.CreateNewChatAfter12Hours))).
					SetCustomID(utils.MakeUserSettings12Hours(s.ID)),
			),
			builders.ActionsRowBuilder(
				builders.ButtonBuilder().
					SetStyle(discordgo.SuccessButton).
					SetLabel("제출").
					SetCustomID(utils.MakeUserSettingsSubmit(s.ID)),
			),
		)

}

func (s *UserSettings) Submit() error {
	if _, err := GetDatabase().Users.Update(s.user.ID, &UserUpdate{
		ChattingMode:              &s.ChattingMode,
		ReplyUser:                 &s.ReplyUser,
		CreateNewChatAfter12Hours: &s.CreateNewChatAfter12Hours,
	}); err != nil {
		return err
	}

	delete(userSettings, s.ID)

	return nil
}
