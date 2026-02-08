package repository

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

func (d *MuffinDatabase) Counts(ctx context.Context, userID string) (*CountDocuments, error) {
	muffinLength, err := d.Texts.coll.CountDocuments(ctx, Text{Persona: "muffin"})
	if err != nil {
		return nil, err
	}

	knowledgeLength, err := d.Knowledge.coll.EstimatedDocumentCount(ctx)
	if err != nil {
		return nil, err
	}

	userKnowledgeLength, err := d.Knowledge.coll.CountDocuments(ctx, Knowledge{UserID: userID})
	if err != nil {
		return nil, err
	}

	chatLength, err := d.Chats.coll.EstimatedDocumentCount(ctx)
	if err != nil {
		return nil, err
	}

	userChatLength, err := d.Chats.coll.CountDocuments(ctx, Chat{UserID: userID})
	if err != nil {
		return nil, err
	}

	memoryLength, err := d.Memory.coll.EstimatedDocumentCount(ctx)
	if err != nil {
		return nil, err
	}

	userMemoryLength, err := d.Memory.coll.CountDocuments(ctx, Memory{UserID: userID})
	if err != nil {
		return nil, err
	}

	sum := muffinLength +
		knowledgeLength +
		chatLength +
		memoryLength

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
