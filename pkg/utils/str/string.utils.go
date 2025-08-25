package str

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"unicode"
)

func ToSnakeCase(s string) string {
	if s == "" {
		return s
	}
	var result []rune
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}

	return string(result)
}

func HashString(str string, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
