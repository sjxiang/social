package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
)


// 返回摘要字符串
func HashifyStr(plainToken string) string {
	h := md5.New()
	h.Write([]byte(plainToken))
	return hex.EncodeToString(h.Sum(nil))
}

func ConvertInteger(source string) (int64, error) {
	integer, err := strconv.ParseInt(source, 10, 64)
	return integer, err
}
