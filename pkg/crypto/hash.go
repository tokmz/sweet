package crypto

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"math/big"
)

// MD5 计算字符串的 MD5 值
func MD5(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

// Hash256 计算字符串的 SHA256 值
func Hash256(str string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(str)))
}

// Hash384 计算字符串的 SHA384 值
func Hash384(str string) string {
	return fmt.Sprintf("%x", sha512.Sum384([]byte(str)))
}

// Hash512 计算字符串的 SHA512 值
func Hash512(str string) string {
	return fmt.Sprintf("%x", sha512.Sum512([]byte(str)))
}

// Base64Encode 对字符串进行 Base64 编码
func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// Base64Decode 对 Base64 编码的字符串进行解码
func Base64Decode(str string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", fmt.Errorf("base64 decode error: %w", err)
	}
	return string(decoded), nil
}

// Salt 生成随机盐值，包含字母和数字
// 默认长度为16位，包含大小写字母和数字
func Salt() string {
	return SaltWithLength(16)
}

// SaltWithLength 生成指定长度的随机盐值
// length: 盐值长度，建议不少于8位
func SaltWithLength(length int) string {
	// 字符集：包含大小写字母和数字
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charsetLen := big.NewInt(int64(len(charset)))
	
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		// 使用加密安全的随机数生成器
		num, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			// 如果随机数生成失败，使用当前时间戳作为备选方案
			// 这种情况极少发生，但为了程序健壮性需要处理
			result[i] = charset[i%len(charset)]
			continue
		}
		result[i] = charset[num.Int64()]
	}
	
	return string(result)
}
