package utils

import (
	"fmt"
	"strconv"
)

func GeneratePrivateChatKey(userA, userB string) string {
	a, _ := strconv.Atoi(userA)
	b, _ := strconv.Atoi(userB)

	if a < b {
		return fmt.Sprintf("%d:%d", a, b)
	}
	return fmt.Sprintf("%d:%d", b, a)
}