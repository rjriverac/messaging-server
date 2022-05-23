package util

import (
	"database/sql"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	str := strings.Builder{}
	str.Grow(n)
	for i, cache, remain := n-1, rand.NewSource(time.Now().UnixNano()).Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.NewSource(time.Now().UnixNano()).Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			str.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return str.String()
}

func RandomUserGen() string {
	str := strings.Builder{}
	str.WriteString(RandomString(10))
	str.WriteString(" ")
	str.WriteString(RandomString(10))
	return str.String()
}

func RandomEmail() string {
	str := strings.Builder{}
	str.WriteString(RandomString(15))
	str.WriteString("@")
	str.WriteString(RandomString(5))
	str.WriteString(".")
	str.WriteString(RandomString(3))
	return str.String()
}

func RandomHashedPW() string {
	var str strings.Builder
	for i := 0; i < 4; i++ {
		num := rand.Intn(100)
		if num%5 == 0 {
			num := rand.Intn(1000)
			str.WriteString(strconv.Itoa(num))
		} else {
			str.WriteString(RandomString(3))
		}
	}
	return str.String()
}

func NullStrGen(n int) sql.NullString {
	var str sql.NullString
	str.String = RandomString(n)
	str.Valid = true
	return str
}
