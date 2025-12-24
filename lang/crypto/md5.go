package crypto

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5 calculates the MD5 hash of a string
func MD5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// MD5Bytes calculates the MD5 hash of bytes
func MD5Bytes(data []byte) string {
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}
