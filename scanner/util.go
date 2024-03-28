package scanner

import (
	"bytes"
	"strconv"
	"unicode"
	"unicode/utf8"
)

func stripCR(b []byte, isMultilineComment bool) []byte {
	c := make([]byte, len(b))
	i := 0
	for j, ch := range b {
		// In a /*-style comment, don't strip \r from *\r/ (incl.
		// sequences of \r from *\r\r...\r/) since the resulting
		// */ would terminate the comment too early unless the \r
		// is immediately following the opening /* in which case
		// it's ok because /*/ is not closed yet (golang/go#11151).
		if ch != '\r' || isMultilineComment && i > len("/*") && c[i-1] == '*' && j+1 < len(b) && b[j+1] == '/' {
			c[i] = ch
			i++
		}
	}
	return c[:i]
}

func trailingDigits(text []byte) (int, int, bool) {
	i := bytes.LastIndexByte(text, ':') // look from right (Windows filenames may contain ':')
	if i < 0 {
		return 0, 0, false // no ":"
	}
	// i >= 0
	n, err := strconv.ParseUint(string(text[i+1:]), 10, 0)
	return i + 1, int(n), err == nil
}

func lower(ch rune) rune     { return ('a' - 'A') | ch } // returns lower-case ch iff ch is ASCII letter
func isDecimal(ch rune) bool { return '0' <= ch && ch <= '9' }
func isHex(ch rune) bool     { return '0' <= ch && ch <= '9' || 'a' <= lower(ch) && lower(ch) <= 'f' }

func isLetter(ch rune) bool {
	return 'a' <= lower(ch) && lower(ch) <= 'z' || ch == '_' || ch >= utf8.RuneSelf && unicode.IsLetter(ch)
}

func isDigit(ch rune) bool {
	return isDecimal(ch) || ch >= utf8.RuneSelf && unicode.IsDigit(ch)
}

func litname(prefix rune) string {
	switch prefix {
	case 'x':
		return "hexadecimal literal"
	case 'o', '0':
		return "octal literal"
	case 'b':
		return "binary literal"
	}
	return "decimal literal"
}

// invalidSep returns the index of the first invalid separator in x, or -1.
func invalidSep(x string) int {
	x1 := ' ' // prefix char, we only care if it's 'x'
	d := '.'  // digit, one of '_', '0' (a digit), or '.' (anything else)
	i := 0

	// a prefix counts as a digit
	if len(x) >= 2 && x[0] == '0' {
		x1 = lower(rune(x[1]))
		if x1 == 'x' || x1 == 'o' || x1 == 'b' {
			d = '0'
			i = 2
		}
	}

	// mantissa and exponent
	for ; i < len(x); i++ {
		p := d // previous digit
		d = rune(x[i])
		switch {
		case d == '_':
			if p != '0' {
				return i
			}
		case isDecimal(d) || x1 == 'x' && isHex(d):
			d = '0'
		default:
			if p == '_' {
				return i - 1
			}
			d = '.'
		}
	}
	if d == '_' {
		return len(x) - 1
	}

	return -1
}

func digitVal(ch rune) int {
	switch {
	case '0' <= ch && ch <= '9':
		return int(ch - '0')
	case 'a' <= lower(ch) && lower(ch) <= 'f':
		return int(lower(ch) - 'a' + 10)
	}
	return 16 // larger than any legal digit val
}

func hexToUint(hex string) (uint64, error) {
	return strconv.ParseUint(hex, 16, 32)
}

func processStringEscapes(s string) string {
	var buf bytes.Buffer

	for i := 0; i < len(s); i++ {
		if s[i] != '\\' {
			buf.WriteByte(s[i])
			continue
		}

		i++
		switch s[i] {
		case 'n':
			buf.WriteByte('\n')
		case 'r':
			buf.WriteByte('\r')
		case 't':
			buf.WriteByte('\t')
		case '\\':
			buf.WriteByte('\\')
		case '"':
			buf.WriteByte('"')
		case 'a':
			buf.WriteByte('\a')
		case 'b':
			buf.WriteByte('\b')
		case 'v':
			buf.WriteByte('\v')
		case 'f':
			buf.WriteByte('\f')
		case 'x':
			i++
			// read two hex digits
			if i+1 >= len(s) {
				buf.WriteString("\\x")
				buf.WriteString(s[i:])
				continue
			}
			b, err := hexToUint(s[i : i+2])
			if err != nil {
				buf.WriteString("\\x")
				buf.WriteString(s[i : i+2])
				i++
				continue
			}
			buf.WriteRune(rune(b))
			i++
		case 'u':
			i++
			// read four hex digits
			if i+3 >= len(s) {
				buf.WriteString("\\u")
				buf.WriteString(s[i:])
				continue
			}
			value, err := hexToUint(s[i : i+4])
			if err != nil {
				buf.WriteString("\\u")
				buf.WriteString(s[i : i+4])
				i += 3
				continue
			}
			buf.WriteRune(rune(value))
			i += 3
		case 'U':
			i++
			// read eight hex digits
			if i+7 >= len(s) {
				buf.WriteString("\\U")
				buf.WriteString(s[i:])
				continue
			}
			codePoint, err := hexToUint(s[i : i+8])
			if err != nil || codePoint > 0x10FFFF {
				buf.WriteString("\\U")
				buf.WriteString(s[i : i+8])
			}
			buf.WriteRune(rune(codePoint))
			i += 7
		default:
			buf.WriteRune('\\')
			buf.WriteByte(s[i])
		}

	}

	return buf.String()
}
