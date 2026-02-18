package customid

import "fmt"

const (
	DeleteKnowledge = "/muffin/knowledge/delete"
	SelectKnowledge = "/muffin/knowledge/select"
)

func MakeDeleteKnowledge(id, userID string) string {
	return fmt.Sprintf("%s/%s/%s", DeleteKnowledge, id, userID)
}

func MakeSelectKnowledge(command string) string {
	return fmt.Sprintf("%s/%s", SelectKnowledge, command)
}
