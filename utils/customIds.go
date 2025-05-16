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

	PaginationEmbedPrev  = "#muffin-pages/prev$"
	PaginationEmbedPages = "#muffin-pages/pages$"
	PaginationEmbedNext  = "#muffin-pages/next$"
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

func MakePaginationEmbedPages(id string, total, current int) string {
	return fmt.Sprintf("%s%s/%d/%d", PaginationEmbedPages, id, total, current)
}

func MakePaginationEmbedNext(id string) string {
	return fmt.Sprintf("%s%s", PaginationEmbedNext, id)
}

func GetPaginationEmbedId(customId string) string {
	if strings.HasPrefix(customId, PaginationEmbedPrev) {
		return customId[len(PaginationEmbedPrev):]
	} else if strings.HasPrefix(customId, PaginationEmbedPages) {
		return customId[len(PaginationEmbedPages):]
	} else {
		return customId[len(PaginationEmbedNext):]
	}
}

func GetPaginationEmbedUserId(id string) string {
	return RegexpPaginationEmbedId.FindAllStringSubmatch(id, 1)[0][1]
}
