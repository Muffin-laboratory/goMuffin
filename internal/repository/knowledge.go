package repository

import (
	"context"
	"slices"
	"sync"
	"time"

	"github.com/Muffin-laboratory/goMuffin/internal/cache"
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
	Collection *mongo.Collection
	caches     *cache.CacheManager[*knowledgeCacheItem]
}

type knowledgeCacheItem struct {
	mu        sync.RWMutex
	knowledge *[]*Knowledge
}

func (c *KnowledgeCollection) Create(ctx context.Context, userID, command, answer string) (*mongo.InsertOneResult, error) {
	data := Knowledge{
		UserID:    userID,
		Command:   command,
		Result:    answer,
		CreatedAt: time.Now(),
	}

	result, err := c.Collection.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}

	if item, ok := c.caches.Get(userID); ok {
		item.mu.Lock()
		defer item.mu.Unlock()

		*item.knowledge = append(*item.knowledge, &data)
	} else {
		c.caches.Set(userID, &knowledgeCacheItem{knowledge: &[]*Knowledge{&data}})
	}

	return result, nil
}

func (c *KnowledgeCollection) Get(ctx context.Context, userID string) ([]*Knowledge, error) {
	if cache, ok := c.caches.Get(userID); ok {
		return *cache.knowledge, nil
	}

	var data []*Knowledge

	cur, err := c.Collection.Find(ctx, Knowledge{UserID: userID})
	if err != nil {
		return data, err
	}

	defer cur.Close(ctx)

	if err = cur.All(ctx, &data); err != nil {
		return data, err
	}

	c.caches.Set(userID, &knowledgeCacheItem{knowledge: &data})

	return data, nil
}

func (c *KnowledgeCollection) GetByCommand(ctx context.Context, command string) ([]*Knowledge, error) {
	var data []*Knowledge

	if caches := c.caches.All(); len(caches) != 0 {
		for _, cache := range caches {
			cache.mu.RLock()

			for _, knowledge := range *cache.knowledge {
				if knowledge.Command == command {
					data = append(data, knowledge)
				}
			}

			cache.mu.RUnlock()
		}

		return data, nil
	}

	cur, err := c.Collection.Find(ctx, Knowledge{
		Command: command,
	})
	if err != nil {
		return data, err
	}

	defer cur.Close(ctx)

	if err = cur.All(ctx, &data); err != nil {
		return data, err
	}

	return data, nil
}

// It doesn't support cache.
func (c *KnowledgeCollection) GetByFilter(ctx context.Context, filter any) ([]*Knowledge, error) {
	var data []*Knowledge

	cur, err := c.Collection.Find(ctx, filter)
	if err != nil {
		return data, err
	}

	defer cur.Close(ctx)

	if err = cur.All(ctx, &data); err != nil {
		return data, err
	}

	return data, nil
}

func (c *KnowledgeCollection) All(ctx context.Context) ([]*Knowledge, error) {
	var data []*Knowledge

	if caches := c.caches.All(); len(caches) != 0 {
		for _, cache := range caches {
			cache.mu.RLock()
			data = append(data, *cache.knowledge...)
			cache.mu.RUnlock()
		}

		return data, nil
	}

	cur, err := c.Collection.Find(ctx, bson.D{})
	if err != nil {
		return data, err
	}

	defer cur.Close(ctx)

	if err = cur.All(ctx, &data); err != nil {
		return data, err
	}

	dataMap := make(map[string]*[]*Knowledge)

	for _, data := range data {
		if list, ok := dataMap[data.UserID]; ok {
			*list = append(*list, data)
			continue
		}

		dataMap[data.UserID] = &[]*Knowledge{data}
	}

	for k, v := range dataMap {
		c.caches.Set(k, &knowledgeCacheItem{knowledge: v})
	}

	return data, nil
}

func (c *KnowledgeCollection) Delete(ctx context.Context, id bson.ObjectID) (*mongo.DeleteResult, error) {
	result, err := c.Collection.DeleteOne(ctx, Knowledge{ID: id})
	if err != nil {
		return nil, err
	}

	caches := c.caches.All()

	if len(caches) == 0 {
		return result, err
	}

	for _, cache := range caches {
		cache.mu.Lock()

		*cache.knowledge = slices.DeleteFunc(*cache.knowledge, func(data *Knowledge) bool {
			return data.ID == id
		})

		cache.mu.Unlock()
	}

	return result, err
}

func (c *KnowledgeCollection) DeleteByUserID(ctx context.Context, userID string) (*mongo.DeleteResult, error) {
	result, err := c.Collection.DeleteMany(ctx, Knowledge{UserID: userID})
	if err != nil {
		return nil, err
	}

	c.caches.Delete(userID)

	return result, nil
}
