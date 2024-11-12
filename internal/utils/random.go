package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().Unix()))
}

func RandomInteger(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[r.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomUserID() int64 {
	return RandomInteger(0, 1000)
}

func RandomRole() string {
	roles := []string{"admin", "user", "guest", "moderator"}
	n := len(roles)

	return roles[rand.Intn(n)]
}

func RandomEmail() string {
	return fmt.Sprintf("%s@vip.cn", RandomString(6))
}
