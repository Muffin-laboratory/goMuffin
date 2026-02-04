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

type Chat struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	Name      string        `bson:"name,omitempty"`
	UserID    string        `bson:"user_id,omitempty"`
	CreatedAt time.Time     `bson:"created_at,omitempty"`
}

type ChatCollection struct {
	coll    *mongo.Collection
	caches  *cache.CacheManager[bson.ObjectID, Chat]
	indexes *cache.CacheManager[string, *indexItem]
}

func newChatCollection(collection *mongo.Collection) *ChatCollection {
	return &ChatCollection{
		coll:    collection,
		caches:  cache.New[bson.ObjectID, Chat](timeToExpire),
		indexes: cache.New[string, *indexItem](timeToExpire),
	}
}

func (c *ChatCollection) createCache(chat Chat) {
	c.caches.Set(chat.ID, chat)

	index := chatIndexBuilder().setUserID(chat.UserID)
	if cache, ok := c.indexes.Get(index.build()); ok {
		cache.mu.Lock()
		cache.ids = append(cache.ids, chat.ID)
		cache.mu.Unlock()
	} else {
		c.indexes.Set(index.build(), indexItemBuilder(chat.ID))
	}

	index.setName(chat.Name)
	if cache, ok := c.indexes.Get(index.build()); ok {
		cache.mu.Lock()
		cache.ids = append(cache.ids, chat.ID)
		cache.mu.Unlock()
	} else {
		c.indexes.Set(index.build(), indexItemBuilder(chat.ID))
	}

}

func (c *ChatCollection) Create(ctx context.Context, userID, name string) (*mongo.InsertOneResult, error) {
	data := Chat{UserID: userID, Name: name, CreatedAt: time.Now()}
	createdChat, err := c.coll.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}

	chatID := createdChat.InsertedID.(bson.ObjectID)
	data.ID = chatID
	if _, err := GetDatabase().Users.Update(ctx, userID, &UserUpdate{
		ChatID: &chatID,
	}); err != nil {
		return nil, err
	}

	c.createCache(data)

	return createdChat, nil
}

func (c *ChatCollection) Find(ctx context.Context, filter query.QueryBuilder) ([]Chat, error) {
	rawFilter := filter.Build()
	index := chatIndexBuilder()
	for _, filter := range rawFilter {
		switch filter.Key {
		case "user_id":
			index.setUserID(filter.Value.(string))
		case "name":
			if value, ok := filter.Value.(string); ok {
				index.setName(value)
			}
		}
	}

	if idx, ok := c.indexes.Get(index.build()); ok {
		var chats []Chat

		for _, id := range idx.ids {
			if cache, ok := c.caches.Get(id); ok {
				chats = append(chats, cache)
			}
		}

		if len(chats) > 0 {
			return chats, nil
		}
	}

	cur, err := c.coll.Find(ctx, rawFilter)
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	var chats []Chat
	var ids []bson.ObjectID

	if err := cur.All(ctx, &chats); err != nil {
		return nil, err
	}

	for _, chat := range chats {
		c.caches.Set(chat.ID, chat)
		ids = append(ids, chat.ID)
	}

	c.indexes.Set(index.build(), indexItemBuilder(ids...))

	return chats, nil
}

func (c *ChatCollection) FindOne(ctx context.Context, filter query.QueryBuilder) (*Chat, error) {
	rawFilter := filter.Build()
	for _, filter := range rawFilter {
		switch filter.Key {
		case "_id":
			if cache, ok := c.caches.Get(filter.Value.(bson.ObjectID)); ok {
				return &cache, nil
			}
		}
	}

	var chat Chat

	if err := c.coll.FindOne(ctx, rawFilter).Decode(&chat); err != nil {
		return nil, err
	}

	c.createCache(chat)

	return &chat, nil
}

func (c *ChatCollection) FindByID(ctx context.Context, id bson.ObjectID) (*Chat, error) {
	return c.FindOne(ctx, query.ChatQueryBuilder().SetID(id))
}

func (c *ChatCollection) DeleteByID(ctx context.Context, id bson.ObjectID) error {
	_, err := c.coll.DeleteOne(ctx, query.ChatQueryBuilder().SetID(id))
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
