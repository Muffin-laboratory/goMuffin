package builders

import (
	"github.com/disgoorg/disgo/discord"
)

func MakeErrorContainer(text string, a ...any) discord.ContainerComponent {
	return discord.NewContainer(
		discord.NewTextDisplay("### ❌ 오류"),
		discord.NewTextDisplayf(text, a...),
	)
}

func MakeDeclineContainer(text string, a ...any) discord.ContainerComponent {
	return discord.NewContainer(
		discord.NewTextDisplay("### ❌ 거부"),
		discord.NewTextDisplayf(text, a...),
	)
}

func MakeCanceledContainer(text string, a ...any) discord.ContainerComponent {
	return discord.NewContainer(
		discord.NewTextDisplay("### ❌ 취소"),
		discord.NewTextDisplayf(text, a...),
	)
}

func MakeSuccessContainer(text string, a ...any) discord.ContainerComponent {
	return discord.NewContainer(
		discord.NewTextDisplay("### ✅ 성공"),
		discord.NewTextDisplayf(text, a...),
	)
}

func MakeUserIsNotRegisteredErrContainer() discord.ContainerComponent {
	return MakeErrorContainer("해당 기능은 등록된 사용자만 쓸 수 있어요. `/가입`으로 가입해주새요.")
}

func MakeUserIsBlockedContainer(globalName, reason string) discord.ContainerComponent {
	return MakeDeclineContainer("- %s님은 서비스에서 차단되었어요.\n> 사유: %s", globalName, reason)
}
