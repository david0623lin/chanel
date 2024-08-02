package lib

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

// 加解密相關工具

// sha256 加密
func (tools *Tools) Sha256(str string) string {
	hash := sha256.Sum256([]byte(str))
	return hex.EncodeToString(hash[:])
}

// md5 加密
func (tools *Tools) Md5(str string) string {
	hash := md5.Sum([]byte(str))
	// 哈希轉十六進制字串
	hashString := hex.EncodeToString(hash[:])

	return hashString
}
