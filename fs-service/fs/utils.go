package fs

import (
	"crypto/sha1"
	"encoding/hex"
)

//NewHash creates unique id based on key provided
func NewHash(key string) string {
	h := sha1.New()
	h.Write([]byte(key))
	return hex.EncodeToString(h.Sum(nil))
}
