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

type indexBuilder interface {
	build() string
}

func indexItemBuilder(ids ...bson.ObjectID) *indexItem {
	return &indexItem{
		mu:  &sync.RWMutex{},
		ids: ids,
	}
}

type chatIndex struct {
	userID *int64
	name   *string
}

func chatIndexBuilder() *chatIndex {
	return &chatIndex{}
}

func (c *chatIndex) setUserID(userID int64) *chatIndex {
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
		index += fmt.Sprintf("userID:%d", *c.userID)
	}

	if c.name != nil {
		index += fmt.Sprintf("&name:%s", *c.name)
	}

	return index
}

type memoryIndex struct {
	chatID *bson.ObjectID
	userID *int64
}

func memoryIndexBuilder() *memoryIndex {
	return &memoryIndex{}
}

func (m *memoryIndex) setChatID(chatID bson.ObjectID) *memoryIndex {
	m.chatID = &chatID
	return m
}

func (m *memoryIndex) setUserID(userID int64) *memoryIndex {
	m.userID = &userID
	return m
}

func (m *memoryIndex) build() string {
	index := "memory:"

	if m.chatID != nil {
		index += fmt.Sprintf("chatID:%s", m.chatID.Hex())
	}

	if m.userID != nil {
		index += fmt.Sprintf("&userID:%d", *m.userID)
	}

	return index
}

type knowledgeIndex struct {
	userID  *int64
	command *string
}

func knowledgeIndexBuilder() *knowledgeIndex {
	return &knowledgeIndex{}
}

func (k *knowledgeIndex) setUserID(userID int64) *knowledgeIndex {
	k.userID = &userID
	return k
}

func (k *knowledgeIndex) setCommand(command string) *knowledgeIndex {
	k.command = &command
	return k
}

func (k *knowledgeIndex) build() string {
	index := "knowledge:"

	if k.userID != nil {
		index += fmt.Sprintf("userID:%d", *k.userID)
	}

	if k.command != nil {
		index += fmt.Sprintf("&command:%s", *k.command)
	}

	return index
}
