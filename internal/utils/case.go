package utils

import (
	"strings"
	"unicode"
)

func MixCase(text string) string {
	var result strings.Builder
	upper := true

	for _, r := range text {
		if unicode.IsLetter(r) {
			if upper {
				r = unicode.ToUpper(r)
			} else {
				r = unicode.ToLower(r)
			}
			upper = !upper
		}
		result.WriteRune(r)
	}

	return result.String()
}
