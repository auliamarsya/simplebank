package utils

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.New(rand.NewSource(time.Now().UnixMicro()))
}

func RandomInt(min, max int64) int64  {
	return min + rand.Int63n(max - min + 1)
}

func RandomString(n int) string {
	var sb strings.Builder
	lengthAlphabet := len(alphabet)

	for i := 0; i < n; i++ {
		getAlphabet := alphabet[rand.Intn(lengthAlphabet)]
		sb.WriteByte(getAlphabet)
	}

	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(0, 10000)
}

func RandomCurrency() string {
	currencies := []string{"IDR", "USD", "EUR"}
	return currencies[rand.Intn(len(currencies))]
}

