package scanner

import (
	"fmt"
	"unicode"

	"github.com/calico32/goose/token"
)

func (s *Scanner) scanString() []stringPart {
	// '"' opening already consumed
	offset := s.offset - 1
	lastOffset := offset
	parts := []stringPart{}

	section := token.StringStart

	for {
		ch := s.ch
		if ch == '\n' || ch < 0 {
			s.error(offset, "string literal not terminated")
			break
		}
		s.next()

		if ch == '\\' {
			s.scanEscape('"')
		}

		if ch == '"' {
			if section == token.StringStart {
				parts = append(parts, stringPart{
					offset:  offset,
					tok:     token.StringStart,
					literal: processStringEscapes(string(s.src[offset : s.offset-1])),
				})
				parts = append(parts, stringPart{
					offset:  s.offset - 1,
					tok:     token.StringEnd,
					literal: "\"",
				})
				break
			} else {
				parts = append(parts, stringPart{
					offset:  lastOffset,
					tok:     token.StringEnd,
					literal: processStringEscapes(string(s.src[lastOffset:s.offset])),
				})
				break
			}
		}

		if ch == '$' {
			parts = append(parts, stringPart{
				offset:  lastOffset,
				tok:     section,
				literal: processStringEscapes(string(s.src[lastOffset : s.offset-1])),
			})
			if section == token.StringStart {
				section = token.StringMid
			}
			parts = append(parts, s.scanInterpolation('"'))
			lastOffset = s.offset
		}

	}

	return parts
}

func (s *Scanner) scanEscape(quote rune) bool {
	offs := s.offset

	var n int
	var base, max uint32
	switch s.ch {
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', '$', quote:
		// typical escapes
		s.next()
		return true
	case 'o':
		// octal \o123
		s.next()
		n, base, max = 3, 8, 255
	case 'x':
		// hex \xff
		s.next()
		n, base, max = 2, 16, 255
	case 'u':
		// unicode \uffff
		s.next()
		n, base, max = 4, 16, unicode.MaxRune
	case 'U':
		// unicode \Uffffffff
		s.next()
		n, base, max = 8, 16, unicode.MaxRune
	default:
		msg := "unknown escape sequence"
		if s.ch < 0 {
			msg = "escape sequence not terminated"
		}
		s.error(offs, msg)
		return false
	}

	var x uint32
	for n > 0 {
		d := uint32(digitVal(s.ch))
		if d >= base {
			msg := fmt.Sprintf("illegal character %#U in escape sequence", s.ch)
			if s.ch < 0 {
				msg = "escape sequence not terminated"
			}
			s.error(s.offset, msg)
			return false
		}
		x = x*base + d
		s.next()
		n--
	}

	if x > max || 0xD800 <= x && x < 0xE000 {
		s.error(offs, "escape sequence is invalid Unicode code point")
		return false
	}

	return true
}

func (s *Scanner) scanInterpolation(quote rune) stringPart {
	// '$' already consumed
	offset := s.offset - 1
	isProperty := false

	if s.ch == '#' {
		s.next()
		isProperty = true
	}

	if isLetter(s.ch) {
		s.scanIdentifier()
		return stringPart{
			offset:  offset,
			tok:     token.StringInterpIdent,
			literal: string(s.src[offset:s.offset]),
		}
	}

	if isProperty {
		// $#{...} is not a valid interpolation
		s.error(s.offset, "expected identifier after $#")
		return stringPart{}
	}

	if s.ch != '{' {
		s.error(s.offset, "expected '{' after $")
		return stringPart{}
	}

	s.next()

	part := stringPart{
		offset:  offset,
		tok:     token.StringInterpExprStart,
		literal: string(s.src[offset:s.offset]),
	}

	for {
		pos, tok, lit := s.Scan()

		if tok == token.RBrace {
			part.tokens = append(part.tokens, &ScanResult{
				pos:     pos,
				tok:     token.StringInterpExprEnd,
				literal: token.RBrace.String(),
			})
			break
		}

		if tok == token.EOF {
			s.error(s.offset, "unexpected EOF in string")
			break
		}

		part.tokens = append(part.tokens, &ScanResult{
			pos:     pos,
			tok:     tok,
			literal: lit,
		})
	}

	return part
}

type stringPart struct {
	offset  int
	tok     token.Token
	literal string
	tokens  []*ScanResult
}
