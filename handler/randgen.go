package handler

import (
	"math/rand"
	"time"
)

const (
	letterBytes   = "abcdefghkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var randUnixNanoSrc = rand.NewSource(time.Now().UnixNano())

func GenerateRandomString(n int) string {
	b := make([]byte, n)
	// A randUnixNanoSrc.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, randUnixNanoSrc.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randUnixNanoSrc.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}
