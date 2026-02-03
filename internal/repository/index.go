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

func indexItemBuilder(ids []bson.ObjectID) *indexItem {
	return &indexItem{&sync.RWMutex{}, ids}
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
