package str

import "unicode"

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
