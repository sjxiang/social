package utils

import (
	"crypto/md5"
	"encoding/hex"
)


// 返回摘要字符串
func HashifyStr(plainToken string) string {
	h := md5.New()
	h.Write([]byte(plainToken))
	return hex.EncodeToString(h.Sum(nil))
}
