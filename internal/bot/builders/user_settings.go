package builders

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/Muffin-laboratory/goMuffin/internal/bot/builders/customid"
	"github.com/Muffin-laboratory/goMuffin/internal/repository"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func getStyleFromBool(k bool) discord.ButtonStyle {
	if k {
		return discord.ButtonStyleSuccess
	}

	return discord.ButtonStyleSecondary
}

func boolToString(k bool) string {
	if k {
		return "활성화"
	}

	return "비활성화"
}

type UserSettings struct {
	ID                        string
	user                      *discord.User
	chattingMode              repository.ChattingMode
	replyUser                 bool
	createNewChatAfter12Hours bool
	prompt                    string
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
		prompt:                    dbUser.Prompt,
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
		discord.NewTextDisplay("- **사용자 지정 프롬프트 설정**\n> 머핀 봇 전체에 적용되는 사용자 지정 프롬프트를 설정해요."),
		discord.NewActionRow(
			discord.NewPrimaryButton(
				fmt.Sprintf("모드: %s", repository.ModeString(s.chattingMode)),
				customid.MakeUserSettingsChattingMode(s.ID),
			),
			discord.NewButton(
				getStyleFromBool(s.replyUser),
				fmt.Sprintf("답장 멘션: %s", boolToString(s.replyUser)),
				customid.MakeUserSettingsReplyUser(s.ID), "", 0,
			),
			discord.NewButton(
				getStyleFromBool(s.createNewChatAfter12Hours),
				fmt.Sprintf("12 시간 후 새로운 채팅: %s", boolToString(s.createNewChatAfter12Hours)),
				customid.MakeUserSettings12Hours(s.ID), "", 0,
			),
			discord.NewPrimaryButton(
				"사용자 지정 프롬프트 설정",
				customid.MakeUserSettingsPrompt(s.ID),
			),
		),
		discord.NewActionRow(
			discord.NewSuccessButton("완료", customid.MakeUserSettingsSubmit(s.ID)),
			discord.NewSecondaryButton("취소", customid.MakeUserSettingsCancel(s.ID)),
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

func (s *UserSettings) PromptModal(e *handler.ComponentEvent) error {
	return e.Modal(
		discord.NewModalCreate(
			customid.MakeUserSettingsPrompt(s.ID),
			"사용자 지정 프롬프트 설정",
			[]discord.LayoutComponent{
				discord.NewLabel(
					"사용자 지정 프롬프트",
					discord.NewParagraphTextInput(customid.UserSettingsPromptSet).
						WithPlaceholder("여기에 프롬프트를 입력...").
						WithValue(s.prompt),
				).
					WithDescription(
						"여기에 사용자 지정 프롬프트를 지정해주세요. " +
							"단, 공란이면 머핀봇의 기본 프롬프트로 설정돼요.",
					),
			},
		),
	)
}

func (s *UserSettings) SetPrompt(prompt string) *UserSettings {
	s.prompt = prompt
	return s
}

func (s *UserSettings) Submit(ctx context.Context) error {
	if _, err := repository.GetDatabase().Users.Update(ctx, int64(s.user.ID), &repository.UserUpdate{
		ChattingMode:              &s.chattingMode,
		ReplyUser:                 &s.replyUser,
		CreateNewChatAfter12Hours: &s.createNewChatAfter12Hours,
		Prompt:                    &s.prompt,
	}); err != nil {
		return err
	}

	delete(userSettings, s.ID)

	return nil
}

func (s *UserSettings) Cancel() {
	delete(userSettings, s.ID)
}
