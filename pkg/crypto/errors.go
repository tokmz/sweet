package crypto

import "errors"

// 定义加密包可能遇到的错误。
var (
	ErrEncryptionFailed  = errors.New("encryption failed")
	ErrDecryptionFailed  = errors.New("decryption failed")
	ErrInvalidKey        = errors.New("invalid encryption key")
	ErrInvalidData       = errors.New("invalid data for encryption or decryption")
	ErrInvalidKeyLength  = errors.New("crypto: 无效的AES密钥长度，必须是16、24或32字节")
	ErrInvalidCiphertext = errors.New("crypto: 无效的密文")
	// 可以根据需要添加更多错误类型
)
