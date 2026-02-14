package builders

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/Muffin-laboratory/goMuffin/internal/utils"
	"github.com/disgoorg/disgo/discord"
)

type UserSettings struct {
	ID                        string
	user                      *discord.User
	ChattingMode              repository.ChattingMode
	ReplyUser                 bool
	CreateNewChatAfter12Hours bool
}

var userSettings = make(map[string]*UserSettings)

func NewUserSettings(ctx context.Context, user discord.User) (*UserSettings, error) {
	dbUser, err := repository.GetDatabase().Users.FindByID(ctx, int64(user.ID))
	if err != nil {
		return nil, err
	}

	id := fmt.Sprintf("%s:%d", user.ID, rand.Intn(100))
	s := &UserSettings{
		ID:                        id,
		user:                      &user,
		ChattingMode:              dbUser.ChattingMode,
		ReplyUser:                 dbUser.ReplyUser,
		CreateNewChatAfter12Hours: dbUser.CreateNewChatAfter12Hours,
	}

	userSettings[id] = s

	return s, nil
}

func GetUserSettings(id string) *UserSettings {
	return userSettings[id]
}

func (s *UserSettings) MakeContainer() discord.ContainerComponent {
	return discord.NewContainer(
		discord.NewSection(
			discord.NewTextDisplay("### 채팅 설정\n- 설정을 개인화 해주세요."),
			discord.NewTextDisplay("- **모드**\n> 해당 봇이 답을 하는 방식이에요. (AI <-> 일반)"),
			discord.NewTextDisplay("- **답장 멘션**\n> 해당 봇이 답할 때 답장 멘션 여부를 정해요."),
		).
			WithAccessory(
				discord.NewThumbnail(*s.user.AvatarURL()),
			),
		discord.NewTextDisplay("- **12 시간 후 새로운 채팅**\n> 해당 봇과 대화하고 12시간 뒤에 새로운 대화를 시작할지 여부를 정해요."),
		discord.NewActionRow(
			discord.NewPrimaryButton(
				fmt.Sprintf("모드: %s", repository.ModeString(s.ChattingMode)),
				utils.MakeUserSettingsChattingMode(s.ID),
			),
			discord.NewButton(
				utils.GetStyleFromBool(s.ReplyUser),
				fmt.Sprintf("답장 멘션: %s", utils.BoolToString(s.ReplyUser)),
				utils.MakeUserSettingsReplyUser(s.ID), "", 0,
			),
			discord.NewButton(
				utils.GetStyleFromBool(s.CreateNewChatAfter12Hours),
				fmt.Sprintf("12 시간 후 새로운 채팅: %s", utils.BoolToString(s.CreateNewChatAfter12Hours)),
				utils.MakeUserSettings12Hours(s.ID), "", 0,
			),
		),
		discord.NewActionRow(
			discord.NewSuccessButton("완료", utils.MakeUserSettingsSubmit(s.ID)),
		),
	)
}

func (s *UserSettings) ToggleChattingMode() {
	newMode := repository.ChattingAIMode

	if s.ChattingMode == repository.ChattingAIMode {
		newMode = repository.ChattingMuffinMode
	}

	s.ChattingMode = newMode
}

func (s *UserSettings) ToggleReplyUser() {
	s.ReplyUser = !s.ReplyUser
}

func (s *UserSettings) Toggle12Hours() {
	s.CreateNewChatAfter12Hours = !s.CreateNewChatAfter12Hours
}

func (s *UserSettings) Submit(ctx context.Context) error {
	if _, err := repository.GetDatabase().Users.Update(ctx, int64(s.user.ID), &repository.UserUpdate{
		ChattingMode:              &s.ChattingMode,
		ReplyUser:                 &s.ReplyUser,
		CreateNewChatAfter12Hours: &s.CreateNewChatAfter12Hours,
	}); err != nil {
		return err
	}

	delete(userSettings, s.ID)

	return nil
}
