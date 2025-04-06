package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

type Aes struct{}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func pkcs7Unpad(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("pkcs7: Data is empty")
	}
	padding := int(data[length-1])
	if padding > length {
		return nil, errors.New("pkcs7: Invalid padding")
	}
	return data[:length-padding], nil
}

func (a *Aes) Encrypt(plainText []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(iv) != aes.BlockSize {
		return nil, errors.New("invalid IV length, must be equal to block size")
	}

	plainText = pkcs7Pad(plainText, aes.BlockSize)
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	copy(cipherText[:aes.BlockSize], iv)

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], plainText)

	return cipherText, nil
}

func (a *Aes) Decrypt(cipherText []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(iv) != aes.BlockSize {
		return nil, errors.New("invalid IV length, must be equal to block size")
	}

	if len(cipherText) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)

	plainText, err := pkcs7Unpad(cipherText)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}
