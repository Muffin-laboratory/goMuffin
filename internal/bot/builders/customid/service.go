package customid

import "fmt"

const (
	ServiceAgree    = "/muffin/service/agree"
	ServiceDisagree = "/muffin/service/disagree"
)

func MakeServiceAgree(userID string) string {
	return fmt.Sprintf("%s/%s", ServiceAgree, userID)
}

func MakeServiceDisagree(userID string) string {
	return fmt.Sprintf("%s/%s", ServiceDisagree, userID)
}
