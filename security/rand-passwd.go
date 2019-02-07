package security

import (
	"math/rand"
	"time"
)

/**
 * Reference:
 * http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
 */

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()_+|-0[]{};':/.?><,"
const (
	letterIDxBits = 6                    // 6 bits to represent a letter index
	letterIDxMask = 1<<letterIDxBits - 1 // All 1-bits, as many as letterIDxBits
	letterIDxMax  = 63 / letterIDxBits   // # of letter indices fitting in 63 bits
)

// DefaultPasswdStrength default password length in bytes
const DefaultPasswdStrength = 32

var src = rand.NewSource(time.Now().UnixNano())

// RandStringBytesMaskImprSrc generates a string of random typable characters of length n
func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIDxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIDxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIDxMax
		}
		if idx := int(cache & letterIDxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIDxBits
		remain--
	}

	return string(b)
}
