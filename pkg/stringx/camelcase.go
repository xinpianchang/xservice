package stringx

// LowerCamelCase to transform string to LowerCamelCase
func LowerCamelCase(s string) string {
	return camelCaseInitCase(s, true)
}

// CamelCase to transform string to CamelCase
func CamelCase(s string) string {
	return camelCaseInitCase(s, false)
}

// CamelCase copy & modify from protobuf
func camelCaseInitCase(s string, lowerFirst bool) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, len(s))
	i := 0
	if s[0] == '_' {
		// Need a capital letter; drop the '_'.
		if lowerFirst {
			t = append(t, 'x')
		} else {
			t = append(t, 'X')
		}
		i++
	}
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && isASCIILower(s[i+1]) {
			continue // Skip the underscore in s.
		}
		if isASCIIDigit(c) {
			t = append(t, c)
			continue
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		if isASCIILower(c) && !(i == 0 && lowerFirst) {
			c ^= ' ' // Make it a capital letter.
		}
		t = append(t, c) // Guaranteed not lower case.
		// Accept lower case sequence that follows.
		for i+1 < len(s) && isASCIILower(s[i+1]) {
			i++
			t = append(t, s[i])
		}
	}
	return string(t)
}

func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// Is c an ASCII digit?
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}
