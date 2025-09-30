package databases

import "context"

type CountDocuments struct {
	All           int64
	Muffin        int64
	Knowledge     int64
	UserKnowledge int64
	Chat          int64
	UserChat      int64
	Memory        int64
	UserMemory    int64
}

func (d *MuffinDatabase) Counts(userID string) (*CountDocuments, error) {
	muffinLength, err := d.Texts.CountDocuments(context.TODO(), Text{Persona: "muffin"})
	if err != nil {
		return nil, err
	}

	knowledgeLength, err := d.Knowledge.Collection.EstimatedDocumentCount(context.TODO())
	if err != nil {
		return nil, err
	}

	userKnowledgeLength, err := d.Knowledge.Collection.CountDocuments(context.TODO(), Knowledge{UserID: userID})
	if err != nil {
		return nil, err
	}

	chatLength, err := d.Chats.EstimatedDocumentCount(context.TODO())
	if err != nil {
		return nil, err
	}

	userChatLength, err := d.Chats.CountDocuments(context.TODO(), Chat{UserID: userID})
	if err != nil {
		return nil, err
	}

	memoryLength, err := d.Memory.Collection.EstimatedDocumentCount(context.TODO())
	if err != nil {
		return nil, err
	}

	userMemoryLength, err := d.Memory.Collection.CountDocuments(context.TODO(), Memory{UserID: userID})
	if err != nil {
		return nil, err
	}

	sum := muffinLength +
		knowledgeLength +
		userKnowledgeLength +
		chatLength +
		userChatLength +
		memoryLength +
		userMemoryLength

	return &CountDocuments{
		All:           sum,
		Muffin:        muffinLength,
		Knowledge:     knowledgeLength,
		UserKnowledge: userKnowledgeLength,
		Chat:          chatLength,
		UserChat:      userChatLength,
		Memory:        memoryLength,
		UserMemory:    userMemoryLength,
	}, nil
}
