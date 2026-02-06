package query

import "go.mongodb.org/mongo-driver/v2/bson"

type QueryBuilder interface {
	Build() bson.D
}

func regexQuery(v string) bson.M {
	return bson.M{"$regex": v}
}
