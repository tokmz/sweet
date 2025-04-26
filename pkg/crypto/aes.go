package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// EncryptAES 使用AES-GCM模式加密数据
// key必须是16、24或32字节长度（分别对应AES-128、AES-192或AES-256）
// 返回加密后的密文和可能的错误
func EncryptAES(key, plaintext []byte) ([]byte, error) {
	if len(plaintext) == 0 {
		return nil, ErrInvalidData
	}

	// 验证密钥长度
	if !isValidAESKey(key) {
		return nil, ErrInvalidKeyLength
	}

	// 创建cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 使用GCM模式，它提供认证加密
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 创建随机数（nonce）
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// 加密数据
	// 密文 = nonce + 实际密文
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// DecryptAES 使用AES-GCM模式解密数据
// key必须是16、24或32字节长度（分别对应AES-128、AES-192或AES-256）
// 返回解密后的明文和可能的错误
func DecryptAES(key, ciphertext []byte) ([]byte, error) {
	if len(ciphertext) == 0 {
		return nil, ErrInvalidCiphertext
	}

	// 验证密钥长度
	if !isValidAESKey(key) {
		return nil, ErrInvalidKeyLength
	}

	// 创建cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 使用GCM模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 检查密文长度是否足够
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, ErrInvalidCiphertext
	}

	// 提取nonce和实际密文
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// 解密数据
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// isValidAESKey 检查密钥长度是否有效
// AES密钥必须是16、24或32字节长度
func isValidAESKey(key []byte) bool {
	kLen := len(key)
	return kLen == 16 || kLen == 24 || kLen == 32
}

// pkcs7Padding 对数据进行PKCS7填充
// 填充数据使得数据长度是块大小的整数倍
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// pkcs7UnPadding 移除PKCS7填充
func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, ErrInvalidData
	}

	padding := int(data[length-1])
	if padding > length {
		return nil, ErrDecryptionFailed
	}

	return data[:length-padding], nil
}

// EncryptAESCBC 使用AES-CBC模式加密数据
// key必须是16、24或32字节长度（分别对应AES-128、AES-192或AES-256）
// 返回加密后的密文和可能的错误
func EncryptAESCBC(key, plaintext []byte) ([]byte, error) {
	if len(plaintext) == 0 {
		return nil, ErrInvalidData
	}

	// 验证密钥长度
	if !isValidAESKey(key) {
		return nil, ErrInvalidKeyLength
	}

	// 创建cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 对数据进行PKCS7填充
	blockSize := block.BlockSize()
	plaintext = pkcs7Padding(plaintext, blockSize)

	// 创建随机IV
	iv := make([]byte, blockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	// 创建CBC模式的加密器
	cbc := cipher.NewCBCEncrypter(block, iv)

	// 加密数据
	ciphertext := make([]byte, len(plaintext))
	cbc.CryptBlocks(ciphertext, plaintext)

	// 将IV附加到密文前面
	// 密文 = IV + 实际密文
	result := make([]byte, len(iv)+len(ciphertext))
	copy(result, iv)
	copy(result[len(iv):], ciphertext)

	return result, nil
}

// DecryptAESCBC 使用AES-CBC模式解密数据
// key必须是16、24或32字节长度（分别对应AES-128、AES-192或AES-256）
// 返回解密后的明文和可能的错误
func DecryptAESCBC(key, ciphertext []byte) ([]byte, error) {
	if len(ciphertext) == 0 {
		return nil, ErrInvalidCiphertext
	}

	// 验证密钥长度
	if !isValidAESKey(key) {
		return nil, ErrInvalidKeyLength
	}

	// 创建cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 检查密文长度是否足够
	blockSize := block.BlockSize()
	if len(ciphertext) < blockSize {
		return nil, ErrInvalidCiphertext
	}

	// 提取IV和实际密文
	iv := ciphertext[:blockSize]
	ciphertext = ciphertext[blockSize:]

	// 检查密文长度是否是块大小的整数倍
	if len(ciphertext)%blockSize != 0 {
		return nil, ErrInvalidCiphertext
	}

	// 创建CBC模式的解密器
	cbc := cipher.NewCBCDecrypter(block, iv)

	// 解密数据
	plaintext := make([]byte, len(ciphertext))
	cbc.CryptBlocks(plaintext, ciphertext)

	// 移除PKCS7填充
	return pkcs7UnPadding(plaintext)
}
