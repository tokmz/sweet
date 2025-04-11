package utils

import (
	"math/rand"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	// 驼峰转换正则
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

// IsEmpty 检查字符串是否为空
func IsEmpty(s string) bool {
	return len(s) == 0
}

// IsBlank 检查字符串是否为空白
func IsBlank(s string) bool {
	if IsEmpty(s) {
		return true
	}
	for _, r := range s {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

// ToSnakeCase 将驼峰命名转换为蛇形命名
func ToSnakeCase(s string) string {
	snake := matchFirstCap.ReplaceAllString(s, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// ToCamelCase 将蛇形命名转换为驼峰命名
func ToCamelCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	n := strings.Builder{}
	n.Grow(len(s))
	capNext := false
	for i, v := range s {
		if v == '_' {
			capNext = true
		} else if capNext {
			n.WriteRune(unicode.ToUpper(v))
			capNext = false
		} else if i == 0 {
			n.WriteRune(unicode.ToLower(v))
		} else {
			n.WriteRune(v)
		}
	}
	return n.String()
}

// ToPascalCase 将蛇形命名转换为帕斯卡命名
func ToPascalCase(s string) string {
	s = ToCamelCase(s)
	if s == "" {
		return s
	}

	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[size:]
}

// Truncate 将字符串截断到指定长度
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}

// RandomString 生成指定长度的随机字符串
func RandomString(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// RandomStringWithCharset 使用指定字符集生成随机字符串
func RandomStringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// RemoveExtraSpaces 移除字符串中的多余空格
func RemoveExtraSpaces(s string) string {
	// 将连续的空白字符替换为单个空格
	re := regexp.MustCompile(`\s+`)
	s = re.ReplaceAllString(s, " ")
	// 去除字符串两端的空格
	return strings.TrimSpace(s)
}

// ReverseString 反转字符串
func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// ContainsAny 检查字符串是否包含任意一个子字符串
func ContainsAny(s string, substrings ...string) bool {
	for _, substr := range substrings {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

// ContainsAll 检查字符串是否包含所有子字符串
func ContainsAll(s string, substrings ...string) bool {
	for _, substr := range substrings {
		if !strings.Contains(s, substr) {
			return false
		}
	}
	return true
}
