package repository

import (
	"context"
	"slices"
	"time"

	"github.com/Muffin-laboratory/goMuffin/internal/cache"
	"github.com/Muffin-laboratory/goMuffin/internal/repository/query"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Knowledge struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	Command   string        `bson:"command,omitempty"`
	Result    string        `bson:"result,omitempty"`
	UserID    string        `bson:"user_id,omitempty"`
	CreatedAt time.Time     `bson:"created_at,omitempty"`
}

type KnowledgeCollection struct {
	coll    *mongo.Collection
	caches  *cache.CacheManager[bson.ObjectID, Knowledge]
	indexes *cache.CacheManager[string, *indexItem]
}

func newKnowledgeCollection(coll *mongo.Collection) *KnowledgeCollection {
	return &KnowledgeCollection{
		coll:    coll,
		caches:  cache.New[bson.ObjectID, Knowledge](timeToExpire),
		indexes: cache.New[string, *indexItem](timeToExpire),
	}
}

func (c *KnowledgeCollection) createCache(data Knowledge) {
	c.caches.Set(data.ID, data)

	userIDIndex := knowledgeIndexBuilder().setUserID(data.UserID)
	createIndexCache(c.indexes, data.ID, userIDIndex)

	commandIndex := knowledgeIndexBuilder().setCommand(data.Command)
	createIndexCache(c.indexes, data.ID, commandIndex)

	index := knowledgeIndexBuilder().setUserID(data.UserID).setCommand(data.Command)
	createIndexCache(c.indexes, data.ID, index)
}

func (c *KnowledgeCollection) Create(ctx context.Context, userID, command, answer string) (*Knowledge, error) {
	data := Knowledge{
		UserID:    userID,
		Command:   command,
		Result:    answer,
		CreatedAt: time.Now(),
	}

	result, err := c.coll.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}

	data.ID = result.InsertedID.(bson.ObjectID)
	c.createCache(data)

	return &data, nil
}

func (c *KnowledgeCollection) Find(ctx context.Context, filter query.QueryBuilder) ([]Knowledge, error) {
	rawFilter := filter.Build()
	index := knowledgeIndexBuilder()
	for _, filter := range rawFilter {
		switch filter.Key {
		case "user_id":
			index.setUserID(filter.Value.(string))
		case "command":
			if value, ok := filter.Value.(string); ok {
				index.setCommand(value)
			}
		}
	}

	if idx, ok := c.indexes.Get(index.build()); ok {
		var knowledge []Knowledge

		idx.mu.RLock()
		for _, id := range idx.ids {
			if cache, ok := c.caches.Get(id); ok {
				knowledge = append(knowledge, cache)
			}
		}
		idx.mu.RUnlock()

		if len(knowledge) > 0 {
			return knowledge, nil
		}
	}

	var knowledge []Knowledge
	var ids []bson.ObjectID

	cur, err := c.coll.Find(ctx, rawFilter)
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	if err = cur.All(ctx, &knowledge); err != nil {
		return nil, err
	}

	for _, data := range knowledge {
		c.caches.Set(data.ID, data)
		ids = append(ids, data.ID)
	}

	c.indexes.Set(index.build(), indexItemBuilder(ids...))

	return knowledge, nil
}

func (c *KnowledgeCollection) DeleteMany(ctx context.Context, filter query.QueryBuilder) error {
	rawFilter := filter.Build()
	_, err := c.coll.DeleteMany(ctx, rawFilter)
	if err != nil {
		return err
	}

	var ids []bson.ObjectID

	index := knowledgeIndexBuilder()
	for _, filter := range rawFilter {
		switch filter.Key {
		case "user_id":
			index.setUserID(filter.Value.(string))
		case "command":
			index.setCommand(filter.Value.(string))
		}
	}

	if idx, ok := c.indexes.Get(index.build()); ok {
		ids = append(ids, idx.ids...)
		c.indexes.Delete(index.build())
	}

	for _, id := range ids {
		c.caches.Delete(id)
	}

	return nil
}

func (c *KnowledgeCollection) DeleteByID(ctx context.Context, id bson.ObjectID) error {
	_, err := c.coll.DeleteOne(ctx, query.KnowledgeQueryBuilder().SetID(id))
	if err != nil {
		return err
	}

	for key, cache := range c.indexes.All() {
		if slices.Contains(cache.ids, id) {
			cache.mu.Lock()
			cache.ids = slices.DeleteFunc(cache.ids, func(cache bson.ObjectID) bool {
				return cache == id
			})
			cache.mu.Unlock()
		}

		if len(cache.ids) == 0 {
			c.indexes.Delete(key)
		}
	}

	c.caches.Delete(id)

	return nil
}
