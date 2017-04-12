package security

import (
	"testing"

	"crypto/subtle"
	"encoding/hex"

	"fmt"
)

func TestHashWithSlt(t *testing.T) {
	var (
		passwd, hash, slt, rehash string
	)
	passwd = RandStringBytesMaskImprSrc(DEFAULT_PASSWD_STRENGTH)
	hash, slt = Hash(passwd)

	fmt.Println("passwd:\t", passwd)
	fmt.Println("hash:\t", hash)
	fmt.Println("slt:\t", slt)

	rehash = HashWithSlt(passwd, slt)

	fmt.Println("rehash:\t", rehash)

	if rehash != hash {
		t.Error("re-hash same passwd and slt should be the same")
		t.Fail()
	}

	hashBytes, _ := hex.DecodeString(hash)
	rehashBytes, _ := hex.DecodeString(rehash)
	if subtle.ConstantTimeCompare(hashBytes, rehashBytes) != 1 {
		t.Error("hex decoded bytes should be the same")
		t.Fail()
	}
}

func TestValidate(t *testing.T) {
	var (
		passwd, hash, slt string
	)
	passwd = RandStringBytesMaskImprSrc(DEFAULT_PASSWD_STRENGTH)
	hash, slt = Hash(passwd)

	fmt.Println("passwd:\t", passwd)
	fmt.Println("slt:\t", slt)
	fmt.Println("hash:\t", hash)

	if !Validate(passwd, slt, hash) {
		t.Error("password should be valid")
		t.Fail()
	}

	slt = RandStringBytesMaskImprSrc(DEFAULT_PASSWD_STRENGTH)
	if Validate(passwd, slt, hash) {
		t.Error("tampered slt should be invalid")
		t.Fail()
	}
}
