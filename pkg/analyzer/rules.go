package analyzer

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Finding represents a rule violation found in a log message.
type Finding struct {
	Rule       string
	Message    string
	Suggestion string // suggested replacement (empty if not applicable)
}

// checkLowercase checks that the first rune of the message is lowercase or a digit.
func checkLowercase(msg string) *Finding {
	if msg == "" {
		return nil
	}
	r, _ := utf8.DecodeRuneInString(msg)
	if r == utf8.RuneError {
		return nil
	}
	if unicode.IsUpper(r) {
		lower := string(unicode.ToLower(r)) + msg[utf8.RuneLen(r):]
		return &Finding{
			Rule:       "lowercase",
			Message:    fmt.Sprintf("log message should start with a lowercase letter, got %q", string(r)),
			Suggestion: lower,
		}
	}
	return nil
}

// checkEnglishOnly checks that all runes are ASCII printable or common whitespace.
func checkEnglishOnly(msg string) *Finding {
	for i, r := range msg {
		if r == '\t' || r == '\n' || r == '\r' {
			continue
		}
		if r < 0x20 || r > 0x7E {
			return &Finding{
				Rule:    "english-only",
				Message: fmt.Sprintf("log message contains non-ASCII character %q at position %d", string(r), i),
			}
		}
	}
	return nil
}

// checkSpecialChars checks for any disallowed special characters and emoji.
//
// By default (AllowedSpecialChars == ""), only ASCII letters, digits and spaces
// are allowed (plus common whitespace \t\n\r).
func checkSpecialChars(msg string, allowedSpecial string) *Finding {
	// Explicit ellipsis checks (common in logs)
	if strings.ContainsRune(msg, '…') {
		return &Finding{
			Rule:    "no-special-chars",
			Message: "log message contains ellipsis '…'",
		}
	}
	if strings.Contains(msg, "...") {
		return &Finding{
			Rule:    "no-special-chars",
			Message: "log message contains ellipsis \"...\"",
		}
	}

	for i, r := range msg {
		// Allow ASCII letters and digits
		if r <= 0x7F {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
				continue
			}
			// Allow spaces and common whitespace
			if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
				continue
			}
			// Allow explicitly configured specials
			if allowedSpecial != "" && strings.ContainsRune(allowedSpecial, r) {
				continue
			}
			return &Finding{
				Rule:    "no-special-chars",
				Message: fmt.Sprintf("log message contains disallowed character %q at position %d", string(r), i),
			}
		}

		// Non-ASCII: if it's emoji, report it (english-only will also fire separately)
		if isEmoji(r) {
			return &Finding{
				Rule:    "no-special-chars",
				Message: fmt.Sprintf("log message contains emoji %q at position %d", string(r), i),
			}
		}
		// Other non-ASCII characters are handled by english-only.
	}

	return nil
}

// isEmoji returns true if the rune falls in common emoji Unicode ranges.
func isEmoji(r rune) bool {
	switch {
	case r >= 0x1F600 && r <= 0x1F64F: // Emoticons
		return true
	case r >= 0x1F300 && r <= 0x1F5FF: // Misc Symbols and Pictographs
		return true
	case r >= 0x1F680 && r <= 0x1F6FF: // Transport and Map
		return true
	case r >= 0x1F1E0 && r <= 0x1F1FF: // Flags
		return true
	case r >= 0x2600 && r <= 0x26FF: // Misc symbols
		return true
	case r >= 0x2700 && r <= 0x27BF: // Dingbats
		return true
	case r >= 0xFE00 && r <= 0xFE0F: // Variation Selectors
		return true
	case r >= 0x1F900 && r <= 0x1F9FF: // Supplemental Symbols
		return true
	case r >= 0x1FA00 && r <= 0x1FA6F: // Chess Symbols
		return true
	case r >= 0x1FA70 && r <= 0x1FAFF: // Symbols Extended-A
		return true
	case r == 0x200D: // Zero Width Joiner
		return true
	}
	return false
}

// checkSensitiveText checks if the literal message itself contains sensitive keywords.
func checkSensitiveText(msg string, keywords map[string]bool) *Finding {
	lower := strings.ToLower(msg)
	for kw := range keywords {
		if kw == "" {
			continue
		}
		if strings.Contains(lower, kw) {
			return &Finding{
				Rule:    "no-sensitive-data",
				Message: fmt.Sprintf("log message may contain sensitive data: message contains keyword %q", kw),
			}
		}
	}
	return nil
}

// checkSensitiveData checks if any of the given identifier parts match sensitive keywords.
func checkSensitiveData(parts []string, keywords map[string]bool) *Finding {
	for _, part := range parts {
		lower := strings.ToLower(part)
		if keywords[lower] {
			return &Finding{
				Rule:    "no-sensitive-data",
				Message: fmt.Sprintf("log message may contain sensitive data: identifier %q matches keyword %q", part, lower),
			}
		}
		for kw := range keywords {
			if len(kw) >= 4 && strings.Contains(lower, kw) {
				return &Finding{
					Rule:    "no-sensitive-data",
					Message: fmt.Sprintf("log message may contain sensitive data: identifier %q contains keyword %q", part, kw),
				}
			}
		}
	}
	return nil
}
