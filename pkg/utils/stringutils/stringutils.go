package stringutils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"unicode"
)

// pascal case / camel case -> snake case
func ToSnakeCase(s string) string {
	if s == "" {
		return s
	}

	var b strings.Builder
	runes := []rune(s)
	n := len(runes)

	for i := range n {
		b.WriteRune(unicode.ToLower(runes[i]))

		nextIsUpper := i+1 < n && unicode.IsUpper(runes[i+1])
		overIsLowerOrNil := (i+2 >= n && unicode.IsLower(runes[i])) ||
			(i+2 < n && unicode.IsLower(runes[i+2]))
		if nextIsUpper && overIsLowerOrNil {
			b.WriteRune('_')
		}

	}

	return b.String()
}

func HashString(stringutils string, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(stringutils))
	return hex.EncodeToString(h.Sum(nil))
}
