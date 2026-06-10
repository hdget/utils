package text

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	nonNumericRegex      = regexp.MustCompile(`[^0-9 ]+`)       // 非数字
	nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`) // 非英文字符和数字
	nonChineseRegex      = regexp.MustCompile(`[^\p{Han}]+`)    // 非汉字
	rxCameling           = regexp.MustCompile(`[\p{L}\p{N}]+`)
)

// OnlyNumeric 去除数字以外的所有字符
func OnlyNumeric(s string) string {
	return nonNumericRegex.ReplaceAllString(s, "")
}

// OnlyAlphaNumeric 去除字符数字以外的所有字符
func OnlyAlphaNumeric(s string) string {
	return nonAlphanumericRegex.ReplaceAllString(s, "")
}

// OnlyChinese 去除中文以外的所有字符
func OnlyChinese(s string) string {
	return nonChineseRegex.ReplaceAllString(s, "")
}

// CleanString 处理字符串, args[0]为是否转换为小写
func CleanString(s string, args ...bool) string {
	// 1. 去除前后空格
	cleanString := strings.TrimSpace(s)

	// 2. 是否转换小写
	toLower := false
	if len(args) > 0 {
		toLower = args[0]
	}

	if toLower {
		cleanString = strings.ToLower(cleanString)
	}

	// 去除不可见字符
	return removeInvisibleCharacter(cleanString)
}

// Capitalize 字符串首字母大写
func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	b := []byte(s)
	if b[0] >= 'a' && b[0] <= 'z' {
		b[0] -= 32
	}
	return string(b)
}

// IsCapitalized 字符串首字母是否大写
func IsCapitalized(s string) bool {
	if len(s) == 0 {
		return false // 空字符串返回false
	}
	firstByte := s[0] // 直接取第一个字节（ASCII字符占1字节）
	return firstByte >= 'A' && firstByte <= 'Z'
}

// CamelToSnake 字符串驼峰转snake风格
func CamelToSnake(s string) string {
	if s == "" {
		return s
	}

	// 预分配足够缓冲区（ASCII每个字符1字节）
	buf := make([]byte, 0, len(s)+5) // +5为额外下划线预留

	var prevLower bool // 记录前一个字符是否是小写
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' { // 仅处理ASCII大写字母
			if i > 0 && prevLower {
				buf = append(buf, '_')
			}
			buf = append(buf, c+32) // 快速转小写（ASCII码+32）
			prevLower = false
		} else {
			buf = append(buf, c)
			prevLower = c >= 'a' && c <= 'z'
		}
	}
	return string(buf)
}

// ToCamelCase converts from underscore separated form to camel case form.
func ToCamelCase(s string) string {
	byteSrc := []byte(s)
	chunks := rxCameling.FindAll(byteSrc, -1)
	for idx, val := range chunks {
		chunks[idx] = cases.Title(language.English).Bytes(val)
	}
	return string(bytes.Join(chunks, nil))
}

// ToSnakeCase converts from camel case form to underscore separated form.
func ToSnakeCase(s string) string {
	s = ToCamelCase(s)
	runes := []rune(s)
	length := len(runes)
	var out []rune
	for i := 0; i < length; i++ {
		out = append(out, unicode.ToLower(runes[i]))
		if i+1 < length && (unicode.IsUpper(runes[i+1]) && unicode.IsLower(runes[i])) {
			out = append(out, '_')
		}
	}

	return string(out)
}

// removeInvisibleCharacter 去除掉不能显示的字符
func removeInvisibleCharacter(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}, s)
}

func Truncate(v any, max int) string {
	if max <= 0 {
		return ""
	}

	var data []byte

	switch x := v.(type) {
	case string:
		data = []byte(x)
	case []byte:
		data = x
	default:
		data = []byte(fmt.Sprint(v))
	}

	if len(data) <= max {
		return string(data)
	}

	const suffix = "..."
	suffixLen := len(suffix)

	if max <= suffixLen {
		return suffix[:max]
	}

	cut := max - suffixLen

	// UTF-8: 避免截断到半个字符
	if utf8.Valid(data) {
		for cut > 0 && (data[cut]&0xC0) == 0x80 {
			cut--
		}
	}

	return string(data[:cut]) + suffix
}
