package security

import "crypto/sha512"
import (
	"crypto/subtle"
	"encoding/hex"

	"github.com/go-playground/log"
)

var shaHash = sha512.New()

func Hash(passwd string) (string, string) {
	slt := RandStringBytesMaskImprSrc(512)
	return HashWithSlt(passwd, slt), slt
}

func HashWithSlt(passwd string, slt string) string {
	shaHash.Reset()
	defer shaHash.Reset()
	shaHash.Write([]byte(passwd))
	shaHash.Write([]byte(slt))
	passwdHash := shaHash.Sum([]byte{})
	return hex.EncodeToString(passwdHash)
}

func Validate(passwd string, slt string, hash string) bool {
	reHash := HashWithSlt(passwd, slt)
	calculated, err := hex.DecodeString(reHash)
	if err != nil {
		log.Error("error: Cannot hex decode hashed message")
		panic("Cannot hex decode hashed message")
	}
	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		log.Error(err)
		return false
	}
	return subtle.ConstantTimeCompare(calculated, hashBytes) == 1
}
