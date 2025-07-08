package chatbot

import (
	"context"

	"git.wh64.net/muffin/goMuffin/databases"
	"google.golang.org/genai"
)

func SaveMemory(data *databases.Memory) error {
	_, err := databases.Database.Memory.InsertOne(context.TODO(), *data)
	return err
}

func GetMemory(userId string) ([]*genai.Content, error) {
	var data []databases.Memory

	memory := []*genai.Content{}

	cur, err := databases.Database.Memory.Find(context.TODO(), databases.User{UserId: userId})
	if err != nil {
		return memory, err
	}

	cur.All(context.TODO(), &data)

	for _, data := range data {
		memory = append(memory,
			genai.NewContentFromText(data.Content, genai.RoleUser),
			genai.NewContentFromText(data.Answer, genai.RoleModel),
		)
	}

	return memory, nil
}
