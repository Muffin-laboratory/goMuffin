package utils

import (
	"fmt"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	DeleteLearnedData = "#muffin/deleteLearnedData@"

	PaginationEmbedPrev    = "#muffin-pages/prev$"
	PaginationEmbedPages   = "#muffin-pages/pages$"
	PaginationEmbedNext    = "#muffin-pages/next$"
	PaginationEmbedModal   = "#muffin-pages/modal$"
	PaginationEmbedSetPage = "#muffin-pages/modal/set$"

	ServiceAgree    = "#muffin/service/agree@"
	ServiceDisagree = "#muffin/service/disagree@"

	DeregisterAgree    = "#muffin/deregister/agree@"
	DeregisterDisagree = "#muffin/deregister/disagree@"

	SelectChat = "#muffin/chat/select@"
)

func MakeDeleteLearnedData(id string, number int, userId string) string {
	return fmt.Sprintf("%sid=%s&no=%d&user_id=%s", DeleteLearnedData, id, number, userId)
}

func GetDeleteLearnedDataId(customId string) (id bson.ObjectID, itemId int) {
	id, _ = bson.ObjectIDFromHex(strings.ReplaceAll(RegexpDLDId.FindAllString(customId, 1)[0], "id=", ""))
	stringItemId := strings.ReplaceAll(RegexpDLDItemId.FindAllString(customId, 1)[0], "no=", "")
	itemId, _ = strconv.Atoi(stringItemId)
	return
}

func GetDeleteLearnedDataUserId(customId string) string {
	return strings.ReplaceAll(RegexpDLDUserId.FindAllString(customId, 1)[0], "user_id=", "")
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

func GetPaginationEmbedId(customId string) string {
	switch {
	case strings.HasPrefix(customId, PaginationEmbedPrev):
		return customId[len(PaginationEmbedPrev):]
	case strings.HasPrefix(customId, PaginationEmbedPages):
		return customId[len(PaginationEmbedPages):]
	case strings.HasPrefix(customId, PaginationEmbedNext):
		return customId[len(PaginationEmbedNext):]
	case strings.HasPrefix(customId, PaginationEmbedModal):
		return customId[len(PaginationEmbedModal):]
	case strings.HasPrefix(customId, PaginationEmbedSetPage):
		return customId[len(PaginationEmbedSetPage):]
	default:
		return customId
	}
}

func GetPaginationEmbedUserId(id string) string {
	return RegexpPaginationEmbedId.FindAllStringSubmatch(id, 1)[0][1]
}

func MakeServiceAgree(userId string) string {
	return fmt.Sprintf("%s%s", ServiceAgree, userId)
}

func MakeServiceDisagree(userId string) string {
	return fmt.Sprintf("%s%s", ServiceDisagree, userId)
}

func GetServiceUserId(customId string) string {
	switch {
	case strings.HasPrefix(customId, ServiceAgree):
		return customId[len(ServiceAgree):]
	case strings.HasPrefix(customId, ServiceDisagree):
		return customId[len(ServiceDisagree):]
	default:
		return customId
	}
}

func MakeDeregisterAgree(userId string) string {
	return fmt.Sprintf("%s%s", DeregisterAgree, userId)
}

func MakeDeregisterDisagree(userId string) string {
	return fmt.Sprintf("%s%s", DeregisterDisagree, userId)
}

func GetDeregisterUserId(customId string) string {
	switch {
	case strings.HasPrefix(customId, DeregisterAgree):
		return customId[len(DeregisterAgree):]
	case strings.HasPrefix(customId, DeregisterDisagree):
		return customId[len(DeregisterDisagree):]
	default:
		return customId
	}
}

func MakeSelectChat(id string, number int, userId string) string {
	return fmt.Sprintf("%sid=%s&no=%d&user_id=%s", DeleteLearnedData, id, number, userId)
}
