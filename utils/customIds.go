package utils

import (
	"fmt"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	DeleteLearnedData       = "#muffin/deleteLearnedData@"
	DeleteLearnedDataUserId = "#muffin/deleteLearnedData@"
	DeleteLearnedDataCancel = "#muffin/deleteLearnedData/cancel@"

	PaginationEmbedPrev    = "#muffin-pages/prev$"
	PaginationEmbedPages   = "#muffin-pages/pages$"
	PaginationEmbedNext    = "#muffin-pages/next$"
	PaginationEmbedModal   = "#muffin-pages/modal$"
	PaginationEmbedSetPage = "#muffin-pages/modal/set$"
)

func MakeDeleteLearnedData(id string, number int) string {
	return fmt.Sprintf("%s%s&No.%d", DeleteLearnedData, id, number)
}

func MakeDeleteLearnedDataUserId(userId string) string {
	return fmt.Sprintf("%s%s", DeleteLearnedDataUserId, userId)
}

func MakeDeleteLearnedDataCancel(id string) string {
	return fmt.Sprintf("%s%s", DeleteLearnedDataCancel, id)
}

func GetDeleteLearnedDataId(customId string) (id bson.ObjectID, itemId int) {
	id, _ = bson.ObjectIDFromHex(strings.ReplaceAll(RegexpItemId.ReplaceAllString(customId[len(DeleteLearnedData):], ""), "&", ""))
	stringItemId := strings.ReplaceAll(RegexpItemId.FindAllString(customId, 1)[0], "No.", "")
	itemId, _ = strconv.Atoi(stringItemId)
	return
}

func GetDeleteLearnedDataUserId(customId string) string {
	if strings.HasPrefix(customId, DeleteLearnedDataCancel) {
		return customId[len(DeleteLearnedDataCancel):]
	} else {
		return customId[len(DeleteLearnedDataUserId):]
	}
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
