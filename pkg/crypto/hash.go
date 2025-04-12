package crypto

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"

	"github.com/google/uuid"
)

type IHash interface {
	// Md5 md5加密
	Md5(str string) string
	// Md5Salt md5加密并添加盐值
	Md5Salt(str string, salt string) string
	// Md5Byte 加密[]byte的MD5
	Md5Byte(str []byte) string
	// Base64Encode base64加密
	Base64Encode(data []byte) string
	// Base64Decode base64解密
	Base64Decode(data string) []byte
	// Sha256 sha256加密
	Sha256(str string) string
	// Sha256Byte sha256加密[]byte
	Sha256Byte(str []byte) string
	// Salt 生成盐值
	Salt() string
}

type hash struct{}

func (s *hash) Md5(str string) string {
	return s.Md5Byte([]byte(str))
}

func (s *hash) Md5Salt(str string, salt string) string {
	return s.Md5Byte([]byte(str + salt))
}

func (s *hash) Md5Byte(str []byte) string {
	h := md5.Sum(str)
	return hex.EncodeToString(h[:])
}

func (s *hash) Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func (s *hash) Base64Decode(data string) []byte {
	result, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return []byte{}
	}
	return result
}

func (s *hash) Sha256(str string) string {
	return s.Sha256Byte([]byte(str))
}

func (s *hash) Sha256Byte(str []byte) string {
	h := sha256.Sum256(str)
	return hex.EncodeToString(h[:])
}

func (s *hash) Salt() string {
	// 生成UUID并截取前12位
	id := uuid.New().String()
	// 去掉UUID中的连字符'-'后截取
	return id[:8] + id[9:13]
}

var Hash IHash = &hash{}
