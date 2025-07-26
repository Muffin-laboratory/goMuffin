package chatbot

import (
	"context"

	"git.wh64.net/muffin/goMuffin/databases"
	"go.mongodb.org/mongo-driver/v2/bson"
	"google.golang.org/genai"
)

func SaveMemory(data *databases.Memory) error {
	_, err := databases.GetDatabase().Memory.InsertOne(context.TODO(), *data)
	return err
}

func GetMemory(chatId bson.ObjectID) ([]*genai.Content, error) {
	var data []databases.Memory

	memory := []*genai.Content{}

	cur, err := databases.GetDatabase().Memory.Find(context.TODO(), databases.User{ChatId: chatId})
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
