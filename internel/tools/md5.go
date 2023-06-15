package tools

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(input string) string {
	// 将输入字符串转换为字节数组
	data := []byte(input)
	// 计算 MD5 值
	hash := md5.Sum(data)
	// 将 MD5 值转换为十六进制字符串
	md5String := hex.EncodeToString(hash[:])
	return md5String
}
