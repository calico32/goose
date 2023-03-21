package scanner

import "github.com/calico32/goose/token"

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
