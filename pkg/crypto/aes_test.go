package crypto

import (
	"bytes"
	"testing"
)

func TestAESEncryptDecrypt(t *testing.T) {
	tests := []struct {
		name      string
		key       []byte
		plaintext []byte
		wantErr   bool
	}{
		{
			name:      "AES-128加密解密",
			key:       []byte("0123456789abcdef"), // 16字节
			plaintext: []byte("这是一段测试文本"),
			wantErr:   false,
		},
		{
			name:      "AES-192加密解密",
			key:       []byte("0123456789abcdef01234567"), // 24字节
			plaintext: []byte("这是另一段测试文本，包含更多内容"),
			wantErr:   false,
		},
		{
			name:      "AES-256加密解密",
			key:       []byte("0123456789abcdef0123456789abcdef"), // 32字节
			plaintext: []byte("这是使用AES-256的测试文本，应该提供最高级别的安全性"),
			wantErr:   false,
		},
		{
			name:      "无效密钥长度",
			key:       []byte("invalid-key"), // 非16/24/32字节
			plaintext: []byte("测试文本"),
			wantErr:   true,
		},
		{
			name:      "空明文",
			key:       []byte("0123456789abcdef"),
			plaintext: []byte(""),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 加密
			ciphertext, err := EncryptAES(tt.key, tt.plaintext)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptAES() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 如果期望错误，则不进行解密测试
			if tt.wantErr {
				return
			}

			// 解密
			decrypted, err := DecryptAES(tt.key, ciphertext)
			if err != nil {
				t.Errorf("DecryptAES() error = %v", err)
				return
			}

			// 验证解密后的文本与原始文本一致
			if !bytes.Equal(decrypted, tt.plaintext) {
				t.Errorf("解密结果不匹配，期望 %s，得到 %s", tt.plaintext, decrypted)
			}
		})
	}
}

func TestInvalidDecryption(t *testing.T) {
	// 测试无效密文
	t.Run("无效密文", func(t *testing.T) {
		key := []byte("0123456789abcdef")
		invalidCiphertext := []byte("这不是有效的密文")

		_, err := DecryptAES(key, invalidCiphertext)
		if err == nil {
			t.Error("期望解密无效密文时返回错误，但没有")
		}
	})

	// 测试密文被篡改
	t.Run("密文被篡改", func(t *testing.T) {
		key := []byte("0123456789abcdef")
		plaintext := []byte("原始文本")

		// 正常加密
		ciphertext, err := EncryptAES(key, plaintext)
		if err != nil {
			t.Fatalf("加密失败: %v", err)
		}

		// 篡改密文（修改最后一个字节）
		if len(ciphertext) > 0 {
			ciphertext[len(ciphertext)-1] ^= 0x01
		}

		// 尝试解密被篡改的密文
		_, err = DecryptAES(key, ciphertext)
		if err == nil {
			t.Error("期望解密被篡改的密文时返回错误，但没有")
		}
	})
}

func TestEncryptDecryptAESCBC(t *testing.T) {
	// 测试用例
	tests := []struct {
		name      string
		key       []byte
		plaintext []byte
		wantErr   bool
	}{
		{
			name:      "有效的AES-128加密解密",
			key:       []byte("0123456789abcdef"), // 16字节
			plaintext: []byte("这是一段需要加密的测试文本"),
			wantErr:   false,
		},
		{
			name:      "有效的AES-192加密解密",
			key:       []byte("0123456789abcdef01234567"), // 24字节
			plaintext: []byte("这是一段需要加密的测试文本"),
			wantErr:   false,
		},
		{
			name:      "有效的AES-256加密解密",
			key:       []byte("0123456789abcdef0123456789abcdef"), // 32字节
			plaintext: []byte("这是一段需要加密的测试文本"),
			wantErr:   false,
		},
		{
			name:      "无效的密钥长度",
			key:       []byte("invalid-key"), // 无效长度
			plaintext: []byte("这是一段需要加密的测试文本"),
			wantErr:   true,
		},
		{
			name:      "空明文",
			key:       []byte("0123456789abcdef"),
			plaintext: []byte(""),
			wantErr:   true,
		},
		{
			name:      "PKCS7填充边界测试",
			key:       []byte("0123456789abcdef"),
			plaintext: bytes.Repeat([]byte("a"), 15), // 15字节，刚好需要1字节填充
			wantErr:   false,
		},
		{
			name:      "PKCS7填充边界测试2",
			key:       []byte("0123456789abcdef"),
			plaintext: bytes.Repeat([]byte("a"), 16), // 16字节，刚好一个块
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 加密
			ciphertext, err := EncryptAESCBC(tt.key, tt.plaintext)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptAESCBC() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return // 如果期望错误，则不进行解密测试
			}

			// 解密
			decrypted, err := DecryptAESCBC(tt.key, ciphertext)
			if err != nil {
				t.Errorf("DecryptAESCBC() error = %v", err)
				return
			}

			// 验证解密后的明文与原始明文是否一致
			if !bytes.Equal(decrypted, tt.plaintext) {
				t.Errorf("解密结果与原始明文不匹配\n原始: %v\n解密: %v", tt.plaintext, decrypted)
			}
		})
	}
}

func TestPKCS7Padding(t *testing.T) {
	tests := []struct {
		name      string
		data      []byte
		blockSize int
		wantLen   int
	}{
		{
			name:      "需要完整块填充",
			data:      []byte("test"),
			blockSize: 16,
			wantLen:   16, // 4 + 12 = 16
		},
		{
			name:      "刚好一个块",
			data:      bytes.Repeat([]byte("a"), 16),
			blockSize: 16,
			wantLen:   32, // 16 + 16 = 32，需要填充一整个块
		},
		{
			name:      "需要部分填充",
			data:      bytes.Repeat([]byte("a"), 10),
			blockSize: 16,
			wantLen:   16, // 10 + 6 = 16
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			padded := pkcs7Padding(tt.data, tt.blockSize)
			if len(padded) != tt.wantLen {
				t.Errorf("pkcs7Padding() 填充后长度 = %v, 期望 %v", len(padded), tt.wantLen)
			}

			// 验证填充值是否正确
			padValue := int(padded[len(padded)-1])
			for i := len(padded) - padValue; i < len(padded); i++ {
				if int(padded[i]) != padValue {
					t.Errorf("填充值不正确，索引 %d 的值为 %d，期望 %d", i, padded[i], padValue)
				}
			}

			// 验证解除填充
			unpadded, err := pkcs7UnPadding(padded)
			if err != nil {
				t.Errorf("pkcs7UnPadding() 错误 = %v", err)
				return
			}

			if !bytes.Equal(unpadded, tt.data) {
				t.Errorf("解除填充后与原始数据不匹配\n原始: %v\n解除填充: %v", tt.data, unpadded)
			}
		})
	}
}

func TestInvalidPKCS7UnPadding(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "空数据",
			data: []byte{},
		},
		{
			name: "填充值大于数据长度",
			data: []byte{1, 2, 3, 20}, // 最后一个字节20大于数据长度4
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := pkcs7UnPadding(tt.data)
			if err == nil {
				t.Errorf("pkcs7UnPadding() 期望错误但没有返回错误")
			}
		})
	}
}
