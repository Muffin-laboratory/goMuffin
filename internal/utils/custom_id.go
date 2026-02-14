package utils

import (
	"fmt"
)

const (
	DeleteKnowledge = "/muffin/knowledge/delete"
	SelectKnowledge = "/muffin/knowledge/select"
)

func MakeDeleteKnowledge(id, userID string) string {
	return fmt.Sprintf("%s/%s/%s", DeleteKnowledge, id, userID)
}

func MakeSelectKnowledge(command string) string {
	return fmt.Sprintf("%s/%s", SelectKnowledge, command)
}

const (
	PaginationContainerFirst   = "/muffin-pages/first"
	PaginationContainerPrev    = "/muffin-pages/prev"
	PaginationContainerPages   = "/muffin-pages/pages"
	PaginationContainerNext    = "/muffin-pages/next"
	PaginationContainerLast    = "/muffin-pages/last"
	PaginationContainerModal   = "/muffin-pages/modal"
	PaginationContainerSetPage = "/muffin-pages/modal/set"
)

func MakePaginationContainerPrev(id string) string {
	return fmt.Sprintf("%s/%s", PaginationContainerPrev, id)
}

func MakePaginationContainerFirst(id string) string {
	return fmt.Sprintf("%s/%s", PaginationContainerFirst, id)
}

func MakePaginationContainerPages(id string) string {
	return fmt.Sprintf("%s/%s", PaginationContainerPages, id)
}

func MakePaginationContainerNext(id string) string {
	return fmt.Sprintf("%s/%s", PaginationContainerNext, id)
}

func MakePaginationContainerLast(id string) string {
	return fmt.Sprintf("%s/%s", PaginationContainerLast, id)
}

func MakePaginationContainerModal(id string) string {
	return fmt.Sprintf("%s/%s", PaginationContainerModal, id)
}

func MakePaginationContainerSetPage(id string) string {
	return fmt.Sprintf("%s%s", PaginationContainerSetPage, id)
}

const (
	ServiceAgree    = "/muffin/service/agree"
	ServiceDisagree = "/muffin/service/disagree"
)

func MakeServiceAgree(userID string) string {
	return fmt.Sprintf("%s/%s", ServiceAgree, userID)
}

func MakeServiceDisagree(userID string) string {
	return fmt.Sprintf("%s/%s", ServiceDisagree, userID)
}

const (
	DeregisterAgree    = "/muffin/deregister/agree"
	DeregisterDisagree = "/muffin/deregister/disagree"
)

func MakeDeregisterAgree(userID string) string {
	return fmt.Sprintf("%s/%s", DeregisterAgree, userID)
}

func MakeDeregisterDisagree(userID string) string {
	return fmt.Sprintf("%s/%s", DeregisterDisagree, userID)
}

const (
	SelectChat       = "/muffin/chat/select"
	DeleteChat       = "/muffin/chat/delete"
	DeleteChatCancel = "/muffin/chat/delete/cancel"
)

func MakeSelectChat(id, userID string) string {
	return fmt.Sprintf("%s/%s/%s", SelectChat, id, userID)
}

func MakeDeleteChat(id, userID string) string {
	return fmt.Sprintf("%s/%s/%s", DeleteChat, id, userID)
}

func MakeDeleteChatCancel(userID string) string {
	return fmt.Sprintf("%s/cancel/%s", DeleteChat, userID)
}

const UserInformationDeregister = "/muffin/info/user/deregister"

func MakeUserInformationDeregister(userID string) string {
	return fmt.Sprintf("%s/%s", UserInformationDeregister, userID)
}

const (
	UserSettings             = "/muffin/user/settings"
	UserSettingsChattingMode = "/muffin/user/settings/chatting_mode"
	UserSettingsReplyUser    = "/muffin/user/settings/reply_user"
	UserSettings12Hours      = "/muffin/user/settings/12hours"
	UserSettingsSubmit       = "/muffin/user/settings/submit"
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

func MakeUserSettingsSubmit(id string) string {
	return fmt.Sprintf("%s/%s", UserSettingsSubmit, id)
}
