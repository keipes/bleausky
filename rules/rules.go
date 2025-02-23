package rules

import (
	"strings"
	"unicode"
)

func PostFilter(text string) bool {
	// Check if the string starts with a capital letter
	if len(text) == 0 || !unicode.IsUpper(rune(text[0])) {
		return false
	}
	// Check if the string ends with a period, question mark, or exclamation mark
	if !strings.HasSuffix(text, ".") && !strings.HasSuffix(text, "?") && !strings.HasSuffix(text, "!") {
		return false
	}
	// Check if all runes in the string are in ASCII printable range
	for _, r := range text {
		if r > 0x7E || r < 0x20 {
			return false
		}
	}
	return true
}
