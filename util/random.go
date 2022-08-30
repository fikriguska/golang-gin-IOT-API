package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const lowercase = "abcdefghijklmnopqrstuvwxyz"

// const uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const numeric = "0123456789"

// const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomString(n int) string {
	var sb strings.Builder
	charset := lowercase + numeric
	s := len(charset)

	for i := 0; i < n; i++ {
		ch := charset[rand.Intn(s)]
		sb.WriteByte(ch)
	}
	return sb.String()
}

func RandomEmail() string {
	return fmt.Sprintf("%s@example.com", RandomString(10))
}
