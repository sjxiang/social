package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// 返回摘要字符串	
func Sum(plainToken string) string {
	// 摘要
	hash := md5.Sum([]byte(plainToken))
	// bytes 转 string
	return hex.EncodeToString(hash[:])
}

