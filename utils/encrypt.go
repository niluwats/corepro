package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

func EncryptStr(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
