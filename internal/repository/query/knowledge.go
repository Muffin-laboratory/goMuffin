package query

import "go.mongodb.org/mongo-driver/v2/bson"

type KnowledgeQuery struct {
	filter bson.D
}

func KnowledgeQueryBuilder() *KnowledgeQuery {
	return &KnowledgeQuery{}
}

func (q *KnowledgeQuery) SetID(id bson.ObjectID) *KnowledgeQuery {
	q.filter = append(q.filter, bson.E{Key: "_id", Value: id})
	return q
}

func (q *KnowledgeQuery) SetUserID(userID string) *KnowledgeQuery {
	q.filter = append(q.filter, bson.E{Key: "user_id", Value: userID})
	return q
}

func (q *KnowledgeQuery) SetCommand(command string) *KnowledgeQuery {
	q.filter = append(q.filter, bson.E{Key: "command", Value: command})
	return q
}

func (q *KnowledgeQuery) SetCommandByRegex(command string) *KnowledgeQuery {
	q.filter = append(q.filter, bson.E{Key: "command", Value: regexQuery(command)})
	return q
}

func (q *KnowledgeQuery) Build() bson.D {
	return q.filter
}
