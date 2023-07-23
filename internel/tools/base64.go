package tools

import "encoding/base64"

// Base64Encode 将输入的字符串进行 Base64 编码
func Base64Encode(input []byte) string {
	encoded := base64.StdEncoding.EncodeToString(input)
	return encoded
}

// Base64Decode 将输入的 Base64 编码字符串进行解码
func Base64Decode(input string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}
