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
	chattingMode              repository.ChattingMode
	replyUser                 bool
	createNewChatAfter12Hours bool
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
		chattingMode:              dbUser.ChattingMode,
		replyUser:                 dbUser.ReplyUser,
		createNewChatAfter12Hours: dbUser.CreateNewChatAfter12Hours,
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
				fmt.Sprintf("모드: %s", repository.ModeString(s.chattingMode)),
				utils.MakeUserSettingsChattingMode(s.ID),
			),
			discord.NewButton(
				utils.GetStyleFromBool(s.replyUser),
				fmt.Sprintf("답장 멘션: %s", utils.BoolToString(s.replyUser)),
				utils.MakeUserSettingsReplyUser(s.ID), "", 0,
			),
			discord.NewButton(
				utils.GetStyleFromBool(s.createNewChatAfter12Hours),
				fmt.Sprintf("12 시간 후 새로운 채팅: %s", utils.BoolToString(s.createNewChatAfter12Hours)),
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

	if s.chattingMode == repository.ChattingAIMode {
		newMode = repository.ChattingMuffinMode
	}

	s.chattingMode = newMode
}

func (s *UserSettings) ToggleReplyUser() {
	s.replyUser = !s.replyUser
}

func (s *UserSettings) Toggle12Hours() {
	s.createNewChatAfter12Hours = !s.createNewChatAfter12Hours
}

func (s *UserSettings) Submit(ctx context.Context) error {
	if _, err := repository.GetDatabase().Users.Update(ctx, int64(s.user.ID), &repository.UserUpdate{
		ChattingMode:              &s.chattingMode,
		ReplyUser:                 &s.replyUser,
		CreateNewChatAfter12Hours: &s.createNewChatAfter12Hours,
	}); err != nil {
		return err
	}

	delete(userSettings, s.ID)

	return nil
}
