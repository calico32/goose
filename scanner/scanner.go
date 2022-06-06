package scanner

import (
	"fmt"
	"path/filepath"
	"unicode"
	"unicode/utf8"

	"github.com/wiisportsresort/goose/token"
)

type ErrorHandler func(pos token.Position, msg string)

type ScanResult struct {
	pos     token.Pos
	tok     token.Token
	literal string
	parts   []*ScanResult
}

type Scanner struct {
	file         *token.File
	dir          string
	src          []byte
	errorHandler ErrorHandler

	scanBufffer []*ScanResult

	ch         rune
	offset     int
	readOffset int
	lineOffset int
	ErrorCount int
}

const (
	BOM = 0xFEFF
	EOF = -1
)

func (s *Scanner) next() {
	if s.readOffset >= len(s.src) {
		s.offset = len(s.src)
		if s.ch == '\n' {
			s.lineOffset = s.offset
			s.file.AddLine(s.offset)
		}
		s.ch = EOF
		return
	}

	s.offset = s.readOffset
	if s.ch == '\n' {
		s.lineOffset = s.offset
		s.file.AddLine(s.offset)
	}
	nextRune := rune(s.src[s.readOffset])

	if nextRune == 0 {
		s.error(s.offset, "illegal character NUL")
	} else if nextRune >= utf8.RuneSelf {
		// not ASCII
		r, w := utf8.DecodeRune(s.src[s.readOffset:])
		s.readOffset += w
		s.ch = r
	}

	s.readOffset++
	s.ch = nextRune
}

func (s *Scanner) scanIdentifier() string {
	offs := s.offset

	// Optimize for the common case of an ASCII identifier.
	//
	// Ranging over s.src[s.rdOffset:] lets us avoid some bounds checks, and
	// avoids conversions to runes.
	//
	// In case we encounter a non-ASCII character, fall back on the slower path
	// of calling into s.next().
	for rdOffset, b := range s.src[s.readOffset:] {
		if 'a' <= b && b <= 'z' || 'A' <= b && b <= 'Z' || b == '_' || '0' <= b && b <= '9' {
			// Avoid assigning a rune for the common case of an ascii character.
			continue
		}
		s.readOffset += rdOffset
		if 0 < b && b < utf8.RuneSelf {
			// Optimization: we've encountered an ASCII character that's not a letter
			// or number. Avoid the call into s.next() and corresponding set up.
			//
			// Note that s.next() does some line accounting if s.ch is '\n', so this
			// shortcut is only possible because we know that the preceding character
			// is not '\n'.
			s.ch = rune(b)
			s.offset = s.readOffset
			s.readOffset++
			goto exit
		}
		// We know that the preceding character is valid for an identifier because
		// scanIdentifier is only called when s.ch is a letter, so calling s.next()
		// at s.rdOffset resets the scanner state.
		s.next()
		for isLetter(s.ch) || isDigit(s.ch) {
			s.next()
		}
		goto exit
	}
	s.offset = len(s.src)
	s.readOffset = len(s.src)
	s.ch = EOF

exit:
	return string(s.src[offs:s.offset])
}

func (s *Scanner) peek() byte {
	if s.readOffset >= len(s.src) {
		return 0
	}

	return s.src[s.readOffset]
}

func (s *Scanner) Init(file *token.File, src []byte, err ErrorHandler) {
	if file.Size() != len(src) {
		panic(fmt.Sprintf("file size (%d) does not match src len (%d)", file.Size(), len(src)))
	}

	s.file = file
	s.dir, _ = filepath.Split(file.Name())
	s.src = src
	s.errorHandler = err

	s.ch = ' '
	s.offset = 0
	s.readOffset = 0
	s.lineOffset = 0
	s.ErrorCount = 0

	s.next()
	if s.ch == BOM {
		// skip BOM at file beginning
		s.next()
	}
}

func (s *Scanner) error(offs int, msg string) {
	if s.errorHandler != nil {
		s.errorHandler(s.file.Position(s.file.Pos(offs)), msg)
	}
	s.ErrorCount++
}

func (s *Scanner) errorf(offs int, format string, args ...any) {
	s.error(offs, fmt.Sprintf(format, args...))
}

func (s *Scanner) scanComment() string {
	// initial '/' already consumed; s.ch == '/' || s.ch == '*'
	offset := s.offset - 1 // position of initial '/'
	next := -1             // position immediately following the comment; < 0 means invalid comment
	numCR := 0

	if s.ch == '/' {
		// double slash comment
		s.next()
		for s.ch != '\n' && s.ch >= 0 {
			// read the rest of the line
			if s.ch == '\r' {
				numCR++
			}
			s.next()
		}

		next = s.offset
		if s.ch == '\n' {
			// newline following double slash not part of comment
			next++
		}
		goto end
	}

	s.next()
	for s.ch >= 0 {
		ch := s.ch
		if ch == '\r' {
			numCR++
		}
		s.next()
		if ch == '*' && s.ch == '/' {
			s.next()
			next = s.offset
			goto end
		}
	}

	s.error(offset, "comment not terminated")

end:
	literal := s.src[offset:s.offset]

	// On Windows, a (//-comment) line may end in "\r\n".
	// Remove the final '\r' (matching the compiler). Remove any
	// other '\r' afterwards (matching the pre-existing be-
	// havior of the scanner).
	if numCR > 0 && len(literal) >= 2 && literal[1] == '/' && literal[len(literal)-1] == '\r' {
		literal = literal[:len(literal)-1]
		numCR--
	}

	if numCR > 0 {
		literal = stripCR(literal, literal[1] == '*')
	}

	return string(literal)
}

// digits accepts the sequence { digit | '_' }.
// If base <= 10, digits accepts any decimal digit but records
// the offset (relative to the source start) of a digit >= base
// in *invalid, if *invalid < 0.
func (s *Scanner) digits(base int, invalid *int) (hadDigits, hadSeparators bool) {
	digsep := 0 // bit 0 == 1 => digit, bit 1 == 1 => underscore present
	if base <= 10 {
		max := rune('0' + base)
		for isDecimal(s.ch) || s.ch == '_' {
			ds := 0b1
			if s.ch == '_' {
				ds = 0b10
			} else if s.ch >= max && *invalid < 0 {
				*invalid = s.offset // record invalid rune offset
			}
			digsep |= ds
			s.next()
		}
	} else {
		for isHex(s.ch) || s.ch == '_' {
			ds := 1
			if s.ch == '_' {
				ds = 2
			}
			digsep |= ds
			s.next()
		}
	}

	hadDigits = digsep&0b1 != 0
	hadSeparators = digsep&0b10 != 0
	return
}

func (s *Scanner) scanNumber() (token.Token, string) {
	offset := s.offset
	tok := token.Illegal

	base := 10        // number base
	prefix := rune(0) // one of 0 (decimal), '0' (0-octal), 'x', 'o', or 'b'
	hadDigits := false
	hadSeparators := false
	invalid := -1

	if s.ch != '.' {
		tok = token.Int
		if s.ch == '0' {
			s.next()
			switch lower(s.ch) {
			case 'x':
				// hexadecimal
				s.next()
				base, prefix = 16, 'x'
			case 'o':
				// octal
				s.next()
				base, prefix = 8, '0'
			case 'b':
				// binary
				s.next()
				base, prefix = 2, 'b'
			default:
				base, prefix = 8, '0'
				hadDigits = true // leading '0'
			}
		}

		d, s := s.digits(base, &invalid)
		hadDigits = hadDigits || d
		hadSeparators = hadSeparators || s
	}

	// fractional part
	if s.ch == '.' {
		tok = token.Float
		if prefix == 'b' || prefix == 'o' {
			s.error(offset, "invalid radix point in "+litname(prefix))
		}
		s.next()
		d, s := s.digits(base, &invalid)
		hadDigits = hadDigits || d
		hadSeparators = hadSeparators || s
	}

	if !hadDigits {
		s.error(s.offset, litname(prefix)+" has no digits")
	}

	// exponent
	if s.ch == 'e' || s.ch == 'E' {
		tok = token.Float
		s.next()
		if s.ch == '-' || s.ch == '+' {
			s.next()
		}
		d, sep := s.digits(10, &invalid)
		hadDigits = hadDigits || d
		hadSeparators = hadSeparators || sep

		if !hadDigits {
			s.error(s.offset, "exponent has no digits")
		}
	}

	literal := string(s.src[offset:s.offset])
	if tok == token.Int && invalid >= 0 {
		s.errorf(invalid, "invalid digit %q in %s", literal[invalid-offset], litname(prefix))
	}
	if hadSeparators {
		if i := invalidSep(literal); i >= 0 {
			s.error(offset+i, "too many successive separators in "+litname(prefix))
		}
	}

	return tok, literal
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

	if isLetter(s.ch) {
		s.scanIdentifier()
		return stringPart{
			offset:  offset,
			tok:     token.StringInterpIdent,
			literal: string(s.src[offset:s.offset]),
		}
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

func (s *Scanner) skipWhitespace() {
	for s.ch == ' ' || s.ch == '\t' || s.ch == '\n' || s.ch == '\r' {
		s.next()
	}
}

func (s *Scanner) Scan() (pos token.Pos, tok token.Token, literal string) {
	if len(s.scanBufffer) != 0 {
		pos, tok, literal, parts := s.scanBufffer[0].pos, s.scanBufffer[0].tok, s.scanBufffer[0].literal, s.scanBufffer[0].parts

		s.scanBufffer = s.scanBufffer[1:]

		// insert parts at the beginning of the buffer
		s.scanBufffer = append(parts, s.scanBufffer...)

		return pos, tok, literal
	}

	s.skipWhitespace()

	pos = s.file.Pos(s.offset)

	switch ch := s.ch; {
	case isLetter(ch):
		literal = s.scanIdentifier()
		if len(literal) > 1 {
			// keywords are longer than one letter - avoid lookup otherwise
			tok = token.Lookup(literal)
		} else {
			tok = token.Ident
		}
	case isDecimal(ch) || ch == '.' && isDecimal(rune(s.peek())):
		tok, literal = s.scanNumber()
	default:
		s.next() // always make progress
		// TODO
		switch ch {
		case EOF:
			tok = token.EOF
		case '\n':
			return pos, token.EOL, "\n"
		case '"':
			parts := s.scanString()
			pos, tok, literal = s.file.Pos(parts[0].offset), parts[0].tok, parts[0].literal
			parts = parts[1:]
			for _, part := range parts {
				s.scanBufffer = append(s.scanBufffer, &ScanResult{
					pos:     s.file.Pos(part.offset),
					tok:     part.tok,
					literal: part.literal,
					parts:   part.tokens,
				})
			}
			return
		case '\'':
			// TODO
		case ':':
			tok = token.Colon
		case ';':
			tok = token.Semi
		case '.':
			tok = token.Period
		case ',':
			tok = token.Comma
		case '(':
			tok = token.LParen
		case ')':
			tok = token.RParen
		case '[':
			tok = token.LBracket
		case ']':
			tok = token.RBracket
		case '{':
			tok = token.LBrace
		case '}':
			tok = token.RBrace
		case '+':
			switch s.ch {
			case '+':
				s.next()
				tok = token.Inc
			case '=':
				s.next()
				tok = token.AddAssign
			default:
				tok = token.Add
			}
		case '-':
			switch s.ch {
			case '-':
				s.next()
				tok = token.Dec
			case '=':
				s.next()
				tok = token.SubAssign
			default:
				tok = token.Sub
			}
		case '*':
			if s.ch == '=' {
				s.next()
				tok = token.MulAssign
			} else if s.ch == '*' {
				s.next()
				if s.ch == '=' {
					s.next()
					tok = token.PowAssign
				} else {
					tok = token.Pow
				}
			} else {
				tok = token.Mul
			}
		case '/':
			switch s.ch {
			case '/':
				fallthrough
			case '*':
				comment := s.scanComment()
				tok = token.Comment
				literal = comment
			case '=':
				s.next()
				tok = token.QuoAssign
			default:
				tok = token.Quo
			}
		case '%':
			if s.ch == '=' {
				s.next()
				tok = token.RemAssign
			} else {
				tok = token.Rem
			}
		case '&':
			if s.ch == '&' {
				s.next()
				tok = token.LogAnd
			} else {
				s.error(s.offset, "invalid character")
			}
		case '|':
			if s.ch == '|' {
				s.next()
				tok = token.LogOr
			} else {
				s.error(s.offset, "invalid character")
			}

		case '=':
			if s.ch == '=' {
				s.next()
				tok = token.Eq
			} else {
				tok = token.Assign
			}
		case '!':
			if s.ch == '=' {
				s.next()
				tok = token.Neq
			} else {
				tok = token.LogNot
			}
		case '<':
			if s.ch == '=' {
				s.next()
				tok = token.Lte
			} else {
				tok = token.Lt
			}
		case '>':
			if s.ch == '=' {

				s.next()
				tok = token.Gte
			} else {
				tok = token.Gt
			}

		default:
			if ch != BOM {
				s.errorf(s.file.Offset(pos), "invalid character %#U", ch)
			}
			tok = token.Illegal
			literal = string(ch)
		}
	}

	return
}
