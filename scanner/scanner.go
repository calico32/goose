package scanner

import (
	"fmt"
	"path/filepath"
	"unicode/utf8"

	"github.com/calico32/goose/token"
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
	anyCh      []interface{}
	anyPeek    []interface{}
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
	} else {
		s.readOffset++
		s.ch = nextRune
	}

	s.anyCh = []interface{}{s.ch}
	s.anyPeek = []interface{}{s.peek()}
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
	s.dir, _ = filepath.Split(file.Specifier())
	s.src = src
	s.errorHandler = err

	s.ch = ' '
	s.anyCh = []interface{}{s.ch}
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

func (s *Scanner) scanComment() string {
	// initial '/' already consumed; s.ch == '/' || s.ch == '*'
	offset := s.offset - 1 // position of initial '/'
	next := -1             // position immediately following the comment; < 0 means invalid comment

	if s.ch == '/' {
		// double slash comment
		s.next()

		for s.ch != '\n' && s.ch != EOF {
			// fmt.Printf("'%c' ", s.ch)
			s.next()
		}

		// fmt.Println()

		next = s.offset
		if s.ch == '\n' {
			// newline following double slash not part of comment
			next++
		}
		goto end
	}

	s.next()
	for s.ch != EOF {
		ch := s.ch
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

	return string(literal)
}

func (s *Scanner) scanProperty() (tok token.Token, lit string) {
	offset := s.offset
	s.next() // consume '#'

	if !isLetter(s.ch) {
		s.error(offset, "expected identifier after '#'")
		tok = token.Illegal
		lit = string(s.src[offset])
		return
	}

	s.scanIdentifier()

	tok = token.Ident
	lit = string(s.src[offset:s.offset])
	return
}

func (s *Scanner) scanSymbol() (tok token.Token, lit string) {
	offset := s.offset
	s.next() // consume '@'

	if !isLetter(s.ch) {
		s.error(offset, "expected identifier after '@'")
		tok = token.Illegal
		lit = string(s.src[offset])
		return
	}

	s.scanIdentifier()

	tok = token.Ident
	lit = string(s.src[offset:s.offset])
	return
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
		if s.ch == '"' {
			s.next()
			parts := s.scanString()
			if len(parts) == 0 {
				// string parsing has failed (error has already been reported)
				return
			}
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
		}

		if ch == '/' && (s.peek() == '/' || s.peek() == '*') {
			s.next()
			tok = token.Comment
			literal = s.scanComment()
			return
		}

		if ch == '#' && s.peek() != '[' {
			tok, literal = s.scanProperty()
			return
		}

		if ch == '@' {
			tok, literal = s.scanSymbol()
			return
		}

		if lookupTok, lookupLit, ok := s.lookup(tokenTable); ok {
			tok = lookupTok
			literal = lookupLit
			break
		}

		s.next()
		if ch != BOM {
			s.errorf(s.file.Offset(pos), "invalid character %#U", ch)
		}
		tok = token.Illegal
		literal = string(ch)
	}

	return
}
