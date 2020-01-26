package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"reflect"
	"unsafe"
)

var secretKey = "4234kxzjcjj3nxnxbcvsjfj"

func generateSalt() string {
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(randomBytes)
}

func BytesToString(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}

func NewHMAC() string {
	message := "This is a sceret everest engineering message"
	salt := generateSalt()

	hash := hmac.New(sha256.New, []byte(secretKey))
	hash.Write([]byte(message + salt))
	return hex.EncodeToString(hash.Sum(nil))
}
