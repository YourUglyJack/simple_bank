package util

import (
	"math/rand"
	"strings"
)

const alphabet = "qwertyuiopasdfghjklzxcvbnm"

func init() {
	// 1.20 后，会自动播种
	// rand.Seed(time.Now().UnixNano())  // Nano是纳米，Unix是秒
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min) // Int63n: [0,max-min+1), [min, max+1) -> [min, max]
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomName() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(1, 1000)
}

func RandomCurrency() string {
	currencies := []string{"CNY", "JPY", "USD"}
	return currencies[rand.Intn(len(currencies))]
}