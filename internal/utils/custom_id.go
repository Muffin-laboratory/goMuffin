package utils

import (
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	DeleteKnowledge = "/muffin/knowledge/delete"
	SelectKnowledge = "#muffin/knowledge$"

	PaginationContainerFirst   = "/muffin-pages/first"
	PaginationContainerPrev    = "/muffin-pages/prev"
	PaginationContainerPages   = "/muffin-pages/pages"
	PaginationContainerNext    = "/muffin-pages/next"
	PaginationContainerLast    = "/muffin-pages/last"
	PaginationContainerModal   = "/muffin-pages/modal"
	PaginationContainerSetPage = "/muffin-pages/modal/set"

	ServiceAgree    = "#muffin/service/agree@"
	ServiceDisagree = "#muffin/service/disagree@"

	DeregisterAgree    = "/muffin/deregister/agree"
	DeregisterDisagree = "/muffin/deregister/disagree"

	SelectChat       = "#muffin/chat/select$"
	DeleteChat       = "/muffin/chat/delete"
	DeleteChatCancel = "/muffin/chat/delete/cancel"

	UserInformationDeregister = "/muffin/info/user/deregister"

	UserSettingsChattingMode = "#muffin/user/set/chatting_mode$"
	UserSettingsReplyUser    = "#muffin/user/set/reply_user$"
	UserSettings12Hours      = "#muffin/user/set/12hours$"
	UserSettingsSubmit       = "#muffin/user/submit$s"
)

func MakeDeleteKnowledge(id, userID string) string {
	return fmt.Sprintf("%s/%s/%s", DeleteKnowledge, id, userID)
}

func GetDeleteKnowledgeID(customID string) bson.ObjectID {
	id, _ := bson.ObjectIDFromHex(RegexpID.FindStringSubmatch(customID)[1])
	return id
}

func GetDeleteKnowledgeUserID(customID string) string {
	return strings.ReplaceAll(RegexpUserID.FindAllString(customID, 1)[0], "user_id=", "")
}

func MakeSelectKnowledge(command string) string {
	return fmt.Sprintf("%s%s", SelectKnowledge, command)
}

func GetSelectKnowledgeCommand(customID string) string {
	return customID[len(SelectKnowledge):]
}

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

func GetPaginationContainerID(customID string) string {
	switch {
	case strings.HasPrefix(customID, PaginationContainerFirst):
		return customID[len(PaginationContainerFirst):]
	case strings.HasPrefix(customID, PaginationContainerPrev):
		return customID[len(PaginationContainerPrev):]
	case strings.HasPrefix(customID, PaginationContainerPages):
		return customID[len(PaginationContainerPages):]
	case strings.HasPrefix(customID, PaginationContainerNext):
		return customID[len(PaginationContainerNext):]
	case strings.HasPrefix(customID, PaginationContainerLast):
		return customID[len(PaginationContainerLast):]
	case strings.HasPrefix(customID, PaginationContainerModal):
		return customID[len(PaginationContainerModal)+1:]
	case strings.HasPrefix(customID, PaginationContainerSetPage):
		return customID[len(PaginationContainerSetPage):]
	default:
		return customID
	}
}

func GetUserID(id string) string {
	return RegexpPaginationContainerID.FindAllStringSubmatch(id, 1)[0][1]
}

func MakeServiceAgree(userID string) string {
	return fmt.Sprintf("%s%s", ServiceAgree, userID)
}

func MakeServiceDisagree(userID string) string {
	return fmt.Sprintf("%s%s", ServiceDisagree, userID)
}

func GetServiceUserID(customID string) string {
	switch {
	case strings.HasPrefix(customID, ServiceAgree):
		return customID[len(ServiceAgree):]
	case strings.HasPrefix(customID, ServiceDisagree):
		return customID[len(ServiceDisagree):]
	default:
		return customID
	}
}

func MakeDeregisterAgree(userID string) string {
	return fmt.Sprintf("%s/%s", DeregisterAgree, userID)
}

func MakeDeregisterDisagree(userID string) string {
	return fmt.Sprintf("%s/%s", DeregisterDisagree, userID)
}

func GetDeregisterUserID(customID string) string {
	switch {
	case strings.HasPrefix(customID, DeregisterAgree):
		return customID[len(DeregisterAgree):]
	case strings.HasPrefix(customID, DeregisterDisagree):
		return customID[len(DeregisterDisagree):]
	default:
		return customID
	}
}

func MakeSelectChat(id, name, userID string) string {
	return fmt.Sprintf("%sid=%s&name=%s&user_id=%s", SelectChat, id, name, userID)
}

func GetChatID(customID string) (id bson.ObjectID, name string) {
	id, _ = bson.ObjectIDFromHex(RegexpID.FindStringSubmatch(customID)[1])
	if !RegexpName.Match([]byte(customID)) {
		return
	}
	name = RegexpName.FindStringSubmatch(customID)[1]
	return
}

func GetChatUserID(customID string) string {
	switch {
	case strings.HasPrefix(customID, SelectChat),
		strings.HasPrefix(customID, DeleteChat):
		return strings.ReplaceAll(RegexpUserID.FindAllString(customID, 1)[0], "user_id=", "")
	case strings.HasPrefix(customID, DeleteChatCancel):
		return customID[len(DeleteChatCancel):]
	default:
		return ""
	}
}

func MakeDeleteChat(id, userID string) string {
	return fmt.Sprintf("%s/%s/%s", DeleteChat, id, userID)
}

func MakeDeleteChatCancel(userID string) string {
	return fmt.Sprintf("%s/cancel/%s", DeleteChat, userID)
}

func MakeUserInformationDeregister(userID string) string {
	return fmt.Sprintf("%s/%s", UserInformationDeregister, userID)
}

func GetUserInformationDeregisterUserID(customID string) string {
	return customID[len(UserInformationDeregister):]
}

func MakeUserSettingsChattingMode(id string) string {
	return fmt.Sprintf("%s%s", UserSettingsChattingMode, id)
}

func MakeUserSettingsReplyUser(id string) string {
	return fmt.Sprintf("%s%s", UserSettingsReplyUser, id)
}

func MakeUserSettings12Hours(id string) string {
	return fmt.Sprintf("%s%s", UserSettings12Hours, id)
}

func MakeUserSettingsSubmit(id string) string {
	return fmt.Sprintf("%s%s", UserSettingsSubmit, id)
}

func GetUserSettingsID(customID string) string {
	switch {
	case strings.HasPrefix(customID, UserSettingsChattingMode):
		return customID[len(UserSettingsChattingMode):]
	case strings.HasPrefix(customID, UserSettingsReplyUser):
		return customID[len(UserSettingsReplyUser):]
	case strings.HasPrefix(customID, UserSettings12Hours):
		return customID[len(UserSettings12Hours):]
	case strings.HasPrefix(customID, UserSettingsSubmit):
		return customID[len(UserSettingsSubmit):]
	default:
		return customID
	}
}
