package repository

import (
	"context"
	"time"

	"github.com/Muffin-laboratory/goMuffin/internal/cache"
	"github.com/Muffin-laboratory/goMuffin/internal/repository/query"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"google.golang.org/genai"
)

type Memory struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	UserID    string        `bson:"user_id,omitempty"`
	Content   string        `bson:"content,omitempty"`
	Answer    string        `bson:"answer,omitempty"`
	ChatID    bson.ObjectID `bson:"chat_id,omitempty"`
	CreatedAt time.Time     `bson:"created_at,omitempty"`
	Files     []File        `bson:"files,omitempty"`
}

func (c *Memory) ToContents() []*genai.Content {
	var parts []*genai.Part
	var memory []*genai.Content

	if len(c.Files) != 0 {
		for _, file := range c.Files {
			parts = append(parts, genai.NewPartFromURI(file.URI, file.MIMEType))
		}
	} else {
		parts = append(parts, genai.NewPartFromText(c.Content))
	}

	memory = append(memory,
		genai.NewContentFromParts(parts, genai.RoleUser),
		genai.NewContentFromText(c.Answer, genai.RoleModel),
	)

	return memory
}

type File struct {
	URI      string `bson:"uri,omitempty"`
	MIMEType string `bson:"mime_type,omitempty"`
}

type MemoryCollection struct {
	coll    *mongo.Collection
	caches  *cache.CacheManager[bson.ObjectID, Memory]
	indexes *cache.CacheManager[string, *indexItem]
}

func newMemoryCollection(collection *mongo.Collection) *MemoryCollection {
	return &MemoryCollection{
		coll:    collection,
		caches:  cache.New[bson.ObjectID, Memory](timeToExpire),
		indexes: cache.New[string, *indexItem](timeToExpire),
	}
}

func (c *MemoryCollection) createCache(memory Memory) {
	c.caches.Set(memory.ID, memory)

	index := memoryIndexBuilder().setChatID(memory.ChatID)
	if cache, ok := c.indexes.Get(index.build()); ok {
		cache.mu.Lock()
		cache.ids = append(cache.ids, memory.ID)
		cache.mu.Unlock()
	} else {
		c.indexes.Set(index.build(), indexItemBuilder(memory.ID))
	}

	index.setUserID(memory.UserID)
	if cache, ok := c.indexes.Get(index.build()); ok {
		cache.mu.Lock()
		cache.ids = append(cache.ids, memory.ID)
		cache.mu.Unlock()
	} else {
		c.indexes.Set(index.build(), indexItemBuilder(memory.ID))
	}
}

func (c *MemoryCollection) Create(ctx context.Context, chatID bson.ObjectID, userID, content, answer string, files []File) error {
	data := Memory{
		UserID:    userID,
		Content:   content,
		Answer:    answer,
		ChatID:    chatID,
		CreatedAt: time.Now(),
		Files:     files,
	}

	createdMemory, err := c.coll.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	data.ID = createdMemory.InsertedID.(bson.ObjectID)

	c.createCache(data)

	return nil
}

func (c *MemoryCollection) Find(ctx context.Context, filter query.QueryBuilder) ([]Memory, error) {
	rawFilter := filter.Build()
	index := memoryIndexBuilder()
	for _, filter := range rawFilter {
		switch filter.Key {
		case "chat_id":
			index.setChatID(filter.Value.(bson.ObjectID))
		case "user_id":
			index.setUserID(filter.Value.(string))
		}
	}

	if idx, ok := c.indexes.Get(index.build()); ok {
		var memory []Memory

		for _, id := range idx.ids {
			if cache, ok := c.caches.Get(id); ok {
				memory = append(memory, cache)
			}
		}

		if len(memory) > 0 {
			return memory, nil
		}
	}

	cur, err := c.coll.Find(ctx, rawFilter)
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	var memory []Memory
	var ids []bson.ObjectID

	if err := cur.All(ctx, &memory); err != nil {
		return nil, err
	}

	for _, memory := range memory {
		c.caches.Set(memory.ID, memory)
		ids = append(ids, memory.ID)
	}

	c.indexes.Set(index.build(), indexItemBuilder(ids...))

	return memory, nil
}

func (c *MemoryCollection) GetLastMemoryTimestamp(ctx context.Context, chatID bson.ObjectID) (int64, error) {
	memory, err := c.Find(ctx, query.MemoryQueryBuilder().SetChatID(chatID))
	if err != nil {
		return 0, err
	}

	return memory[len(memory)-1].CreatedAt.Unix(), nil
}

func (c *MemoryCollection) CountDocuments(ctx context.Context, filter query.QueryBuilder) (int64, error) {
	return c.coll.CountDocuments(ctx, filter.Build())
}

func (c *MemoryCollection) DeleteMany(ctx context.Context, filter query.QueryBuilder) error {
	rawFilter := filter.Build()
	_, err := c.coll.DeleteMany(ctx, rawFilter)
	if err != nil {
		return err
	}

	var ids []bson.ObjectID

	index := memoryIndexBuilder()
	for _, filter := range rawFilter {
		switch filter.Key {
		case "chat_id":
			index.setChatID(filter.Value.(bson.ObjectID))
		case "user_id":
			index.setUserID(filter.Value.(string))
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
