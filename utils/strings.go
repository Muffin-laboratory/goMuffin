package utils

import "fmt"

func AddPrefix(prefix string, arr []string) (newArr []string) {
	for _, item := range arr {
		newArr = append(newArr, fmt.Sprintf("%s%s", prefix, item))
	}
	return
}
