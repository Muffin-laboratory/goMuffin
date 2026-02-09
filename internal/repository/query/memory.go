package query

import "go.mongodb.org/mongo-driver/v2/bson"

type MemoryQuery struct {
	filter bson.D
}

func MemoryQueryBuilder() *MemoryQuery {
	return &MemoryQuery{}
}

func (q *MemoryQuery) SetChatID(chatID bson.ObjectID) *MemoryQuery {
	q.filter = append(q.filter, bson.E{Key: "chat_id", Value: chatID})
	return q
}

func (q *MemoryQuery) SetUserID(userID int64) *MemoryQuery {
	q.filter = append(q.filter, bson.E{Key: "user_id", Value: userID})
	return q
}

func (q *MemoryQuery) Build() bson.D {
	return q.filter
}
