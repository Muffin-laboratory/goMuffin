package customid

import "fmt"

const (
	DeregisterAgree           = "/muffin/deregister/agree"
	DeregisterDisagree        = "/muffin/deregister/disagree"
	UserInformationDeregister = "/muffin/info/user/deregister"
)

func MakeDeregisterAgree(userID string) string {
	return fmt.Sprintf("%s/%s", DeregisterAgree, userID)
}

func MakeDeregisterDisagree(userID string) string {
	return fmt.Sprintf("%s/%s", DeregisterDisagree, userID)
}

func MakeUserInformationDeregister(userID string) string {
	return fmt.Sprintf("%s/%s", UserInformationDeregister, userID)
}
