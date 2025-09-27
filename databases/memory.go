package databases

import (
	"context"
	"slices"
	"sync"
	"time"

	"git.wh64.net/muffin/goMuffin/cache"
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
}

type memoryCacheItem struct {
	mu     sync.RWMutex
	memory *[]*Memory
}

type MemoryCollection struct {
	Collection *mongo.Collection
	caches     *cache.CacheManager[*memoryCacheItem]
}

func (c *MemoryCollection) Save(chatID bson.ObjectID, userID, content, answer string) error {
	data := Memory{
		UserID:  userID,
		Content: content,
		Answer:  answer,
		ChatID:  chatID,
	}

	if _, err := c.Collection.InsertOne(context.TODO(), data); err != nil {
		return err
	}

	if item, ok := c.caches.Get(chatID.Hex()); ok {
		item.mu.Lock()
		defer item.mu.Unlock()
		*item.memory = append(*item.memory, &data)
	} else {
		c.caches.Set(data.ChatID.Hex(), &memoryCacheItem{memory: &[]*Memory{&data}})
	}

	return nil
}

func (c *MemoryCollection) Get(chatID bson.ObjectID) ([]*genai.Content, error) {
	var memory []*genai.Content

	if item, ok := c.caches.Get(chatID.Hex()); ok {
		item.mu.RLock()
		defer item.mu.RUnlock()

		for _, cache := range *item.memory {
			memory = append(memory,
				genai.NewContentFromText(cache.Content, genai.RoleUser),
				genai.NewContentFromText(cache.Answer, genai.RoleModel),
			)
		}
	} else {
		var data []*Memory

		cur, err := c.Collection.Find(context.TODO(), User{ChatID: chatID})
		if err != nil {
			return memory, err
		}

		defer cur.Close(context.TODO())

		if err = cur.All(context.TODO(), &data); err != nil {
			return memory, err
		}

		if len(data) == 0 {
			return memory, nil
		}

		for _, data := range data {
			memory = append(memory,
				genai.NewContentFromText(data.Content, genai.RoleUser),
				genai.NewContentFromText(data.Answer, genai.RoleModel),
			)
		}

		c.caches.Set(chatID.Hex(), &memoryCacheItem{memory: &data})
	}

	return memory, nil
}

func (c *MemoryCollection) DeleteByUserID(userID string) (*mongo.DeleteResult, error) {
	var memory []Memory
	var chatIDList []bson.ObjectID

	cur, err := c.Collection.Find(context.TODO(), Memory{UserID: userID})
	if err != nil {
		return nil, err
	}

	if err = cur.All(context.TODO(), &memory); err != nil {
		return nil, err
	}

	if len(memory) == 0 {
		return nil, nil
	}

	for _, memory := range memory {
		if slices.Contains(chatIDList, memory.ChatID) {
			continue
		}

		chatIDList = append(chatIDList, memory.ChatID)
	}

	result, err := c.Collection.DeleteMany(context.TODO(), Memory{UserID: userID})
	if err != nil {
		return nil, err
	}

	for _, id := range chatIDList {
		c.caches.Delete(id.Hex())
	}

	return result, nil
}

func (c *MemoryCollection) DeleteByChatID(chatID bson.ObjectID) (*mongo.DeleteResult, error) {
	result, err := c.Collection.DeleteMany(context.TODO(), Memory{ChatID: chatID})
	if err != nil {
		return nil, err
	}

	c.caches.Delete(chatID.Hex())
	return result, err
}
