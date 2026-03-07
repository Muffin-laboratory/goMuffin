package customid

import "fmt"

const (
	CreateChat          = "/muffin/chat/create"
	CreateChatSetPrompt = "/muffin/chat/create/set/prompt"
	SelectChat          = "/muffin/chat/select"
	DeleteChat          = "/muffin/chat/delete"
	DeleteChatCancel    = "/muffin/chat/delete/cancel"
)

func MakeCreateChat(name string) string {
	return fmt.Sprintf("%s/%s", CreateChat, name)
}

func MakeSelectChat(id, userID string) string {
	return fmt.Sprintf("%s/%s/%s", SelectChat, id, userID)
}

func MakeDeleteChat(id, userID string) string {
	return fmt.Sprintf("%s/%s/%s", DeleteChat, id, userID)
}

func MakeDeleteChatCancel(userID string) string {
	return fmt.Sprintf("%s/cancel/%s", DeleteChat, userID)
}
