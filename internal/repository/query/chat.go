package query

import "go.mongodb.org/mongo-driver/v2/bson"

type ChatQuery struct {
	filter bson.D
}

func ChatQueryBuilder() *ChatQuery {
	return &ChatQuery{}
}

func (b *ChatQuery) SetUserID(userID string) *ChatQuery {
	b.filter = append(b.filter, bson.E{Key: "user_id", Value: userID})
	return b
}

func (b *ChatQuery) SetName(name string) *ChatQuery {
	b.filter = append(b.filter, bson.E{Key: "name", Value: name})
	return b
}

func (b *ChatQuery) SetNameByRegex(name string) *ChatQuery {
	b.filter = append(b.filter, bson.E{Key: "name", Value: regexQuery(name)})
	return b
}

func (b *ChatQuery) SetID(id bson.ObjectID) *ChatQuery {
	b.filter = append(b.filter, bson.E{Key: "_id", Value: id})
	return b
}

func (b *ChatQuery) Build() bson.D {
	return b.filter
}
