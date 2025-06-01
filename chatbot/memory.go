package chatbot

import (
	"context"

	"git.wh64.net/muffin/goMuffin/databases"
	"go.mongodb.org/mongo-driver/v2/bson"
	"google.golang.org/genai"
)

func SaveMemory(data *databases.InsertMemory) error {
	_, err := databases.Database.Memory.InsertOne(context.TODO(), *data)
	return err
}

func GetMemory(userId string) ([]*genai.Content, error) {
	var data []databases.Memory

	MAX_LENGTH := 50
	memory := []*genai.Content{}

	cur, err := databases.Database.Memory.Find(context.TODO(), bson.D{{Key: "user_id", Value: userId}})
	if err != nil {
		return memory, err
	}

	cur.All(context.TODO(), &data)

	if len(data) > MAX_LENGTH {
		data = data[MAX_LENGTH:]
	}

	for _, data := range data {
		memory = append(memory,
			genai.NewContentFromText(data.Content, genai.RoleUser),
			genai.NewContentFromText(data.Answer, genai.RoleModel),
		)
	}

	return memory, nil
}
