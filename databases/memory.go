package databases

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"google.golang.org/genai"
)

type Memory struct {
	ID      bson.ObjectID `bson:"_id,omitempty"`
	UserID  string        `bson:"user_id,omitempty"`
	Content string        `bson:"content,omitempty"`
	Answer  string        `bson:"answer,omitempty"`
	ChatID  bson.ObjectID `bson:"chat_id,omitempty"`
}

type MemoryCollection struct {
	*mongo.Collection
}

func (c *MemoryCollection) Save(data *Memory) error {
	_, err := c.InsertOne(context.TODO(), *data)
	return err
}

func (c *MemoryCollection) Get(chatId bson.ObjectID) ([]*genai.Content, error) {
	var data []Memory

	memory := []*genai.Content{}

	cur, err := c.Find(context.TODO(), User{ChatID: chatId})
	if err != nil {
		return memory, err
	}

	err = cur.All(context.TODO(), &data)
	if err != nil {
		return memory, err
	}

	for _, data := range data {
		memory = append(memory,
			genai.NewContentFromText(data.Content, genai.RoleUser),
			genai.NewContentFromText(data.Answer, genai.RoleModel),
		)
	}

	return memory, nil
}
