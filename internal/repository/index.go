package repository

import (
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type indexItem struct {
	mu  *sync.RWMutex
	ids []bson.ObjectID
}

func indexItemBuilder(ids ...bson.ObjectID) *indexItem {
	return &indexItem{
		mu:  &sync.RWMutex{},
		ids: ids,
	}
}

type chatIndex struct {
	userID *string
	name   *string
}

func chatIndexBuilder() *chatIndex {
	return &chatIndex{}
}

func (c *chatIndex) setUserID(userID string) *chatIndex {
	c.userID = &userID
	return c
}

func (c *chatIndex) setName(name string) *chatIndex {
	c.name = &name
	return c
}

func (c *chatIndex) build() string {
	index := "chat:"

	if c.userID != nil {
		index += fmt.Sprintf("userID:%s", *c.userID)
	}

	if c.name != nil {
		index += fmt.Sprintf("&name:%s", *c.name)
	}

	return index
}

type memoryIndex struct {
	chatID *bson.ObjectID
	userID *string
}

func memoryIndexBuilder() *memoryIndex {
	return &memoryIndex{}
}

func (m *memoryIndex) setChatID(chatID bson.ObjectID) *memoryIndex {
	m.chatID = &chatID
	return m
}

func (m *memoryIndex) setUserID(userID string) *memoryIndex {
	m.userID = &userID
	return m
}

func (m *memoryIndex) build() string {
	index := "memory:"

	if m.chatID != nil {
		index += fmt.Sprintf("chatID:%s", m.chatID.Hex())
	}

	if m.userID != nil {
		index += fmt.Sprintf("&userID:%s", *m.userID)
	}

	return index
}
