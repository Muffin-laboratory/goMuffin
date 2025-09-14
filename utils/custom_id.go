package utils

import (
	"fmt"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	DeleteLearnedData = "#muffin/deleteLearnedData$"

	PaginationEmbedPrev    = "#muffin-pages/prev$"
	PaginationEmbedPages   = "#muffin-pages/pages$"
	PaginationEmbedNext    = "#muffin-pages/next$"
	PaginationEmbedModal   = "#muffin-pages/modal$"
	PaginationEmbedSetPage = "#muffin-pages/modal/set$"

	ServiceAgree    = "#muffin/service/agree@"
	ServiceDisagree = "#muffin/service/disagree@"

	DeregisterAgree    = "#muffin/deregister/agree@"
	DeregisterDisagree = "#muffin/deregister/disagree@"

	SelectChat       = "#muffin/chat/select$"
	DeleteChat       = "#muffin/chat/delete$"
	DeleteChatCancel = "#muffin/chat/delete/cancel@"

	UserInformationDeregister = "#muffin/info/user/deregister@"
)

func MakeDeleteLearnedData(id string, number int, userID string) string {
	return fmt.Sprintf("%sid=%s&no=%d&user_id=%s", DeleteLearnedData, id, number, userID)
}

func GetDeleteLearnedDataID(customID string) (id bson.ObjectID, itemID int) {
	id, _ = bson.ObjectIDFromHex(strings.ReplaceAll(RegexpID.FindAllString(customID, 1)[0], "id=", ""))
	stringItemId := strings.ReplaceAll(RegexpItemID.FindAllString(customID, 1)[0], "no=", "")
	itemID, _ = strconv.Atoi(stringItemId)
	return
}

func GetDeleteLearnedDataUserId(customID string) string {
	return strings.ReplaceAll(RegexpUserID.FindAllString(customID, 1)[0], "user_id=", "")
}

func MakePaginationEmbedPrev(id string) string {
	return fmt.Sprintf("%s%s", PaginationEmbedPrev, id)
}

func MakePaginationEmbedPages(id string) string {
	return fmt.Sprintf("%s%s", PaginationEmbedPages, id)
}

func MakePaginationEmbedNext(id string) string {
	return fmt.Sprintf("%s%s", PaginationEmbedNext, id)
}

func MakePaginationEmbedModal(id string) string {
	return fmt.Sprintf("%s%s", PaginationEmbedModal, id)
}

func MakePaginationEmbedSetPage(id string) string {
	return fmt.Sprintf("%s%s", PaginationEmbedSetPage, id)
}

func GetPaginationEmbedID(customID string) string {
	switch {
	case strings.HasPrefix(customID, PaginationEmbedPrev):
		return customID[len(PaginationEmbedPrev):]
	case strings.HasPrefix(customID, PaginationEmbedPages):
		return customID[len(PaginationEmbedPages):]
	case strings.HasPrefix(customID, PaginationEmbedNext):
		return customID[len(PaginationEmbedNext):]
	case strings.HasPrefix(customID, PaginationEmbedModal):
		return customID[len(PaginationEmbedModal):]
	case strings.HasPrefix(customID, PaginationEmbedSetPage):
		return customID[len(PaginationEmbedSetPage):]
	default:
		return customID
	}
}

func GetPaginationEmbedUserID(id string) string {
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
	return fmt.Sprintf("%s%s", DeregisterAgree, userID)
}

func MakeDeregisterDisagree(userID string) string {
	return fmt.Sprintf("%s%s", DeregisterDisagree, userID)
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

func MakeSelectChat(id string, number int, userID string) string {
	return fmt.Sprintf("%sid=%s&no=%d&user_id=%s", SelectChat, id, number, userID)
}

func GetSelectChatID(customID string) (id bson.ObjectID, itemID int) {
	id, _ = bson.ObjectIDFromHex(strings.ReplaceAll(RegexpID.FindAllString(customID, 1)[0], "id=", ""))
	stringItemId := strings.ReplaceAll(RegexpItemID.FindAllString(customID, 1)[0], "no=", "")
	itemID, _ = strconv.Atoi(stringItemId)
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

func MakeDeleteChat(id string, number int, userID string) string {
	return fmt.Sprintf("%sid=%s&no=%d&user_id=%s", DeleteChat, id, number, userID)
}

func GetDeleteChatID(customID string) (id bson.ObjectID, itemID int) {
	id, _ = bson.ObjectIDFromHex(strings.ReplaceAll(RegexpID.FindAllString(customID, 1)[0], "id=", ""))
	stringItemId := strings.ReplaceAll(RegexpItemID.FindAllString(customID, 1)[0], "no=", "")
	itemID, _ = strconv.Atoi(stringItemId)
	return
}

func MakeDeleteChatCancel(userID string) string {
	return fmt.Sprintf("%s%s", DeleteChatCancel, userID)
}

func MakeUserInformationDeregister(userID string) string {
	return fmt.Sprintf("%s%s", UserInformationDeregister, userID)
}

func GetUserInformationDeregisterUserID(customID string) string {
	return customID[len(UserInformationDeregister):]
}
