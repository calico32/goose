package scanner

import "github.com/calico32/goose/token"

type tokens map[rune]any

func (s *Scanner) lookup(table tokens) (tok token.Token, literal string, ok bool) {
	ch := s.ch

	// check if it's in the table
	if next, ok := table[ch]; ok {
		switch next := next.(type) {
		case token.Token:
			// the current sequence is a token, consume the character
			s.next()
			return next, "", true
		case tokens:
			// the current sequence is a prefix of a token, consume the character and continue
			s.next()
			return s.lookup(next)
		}
	}

	if tok, ok := table[0]; ok {
		// everything we already consumed is a token
		return tok.(token.Token), "", true
	}

	// didn't find a token
	return token.Illegal, "", false
}

var tokenTable = tokens{
	'\n': token.EOL,
	EOF:  token.EOF,

	'#': tokens{
		0:   token.Illegal,
		'[': token.HashLBracket,
	},

	':': tokens{
		0:   token.Colon,
		':': token.Bind,
	},
	';': token.Semi,
	'.': tokens{
		0: token.Period,
		'.': tokens{
			0:   token.Illegal,
			'.': token.Ellipsis,
		},
	},
	',': token.Comma,
	'(': token.LParen,
	')': token.RParen,
	'[': token.LBracket,
	']': token.RBracket,
	'{': token.LBrace,
	'}': token.RBrace,
	'+': tokens{
		0:   token.Add,
		'+': token.Inc,
		'=': token.AddAssign,
	},
	'-': tokens{
		0:   token.Sub,
		'-': token.Dec,
		'=': token.SubAssign,
		'>': token.Arrow,
	},
	'*': tokens{
		0: token.Mul,
		'*': tokens{
			0:   token.Pow,
			'=': token.PowAssign,
		},
		'=': token.MulAssign,
	},
	'/': tokens{
		0:   token.Quo,
		'=': token.QuoAssign,
	},
	'%': tokens{
		0:   token.Rem,
		'=': token.RemAssign,
	},
	'&': tokens{
		0:   token.BitAnd,
		'=': token.BitAndAssign,
		'&': tokens{
			0:   token.LogAnd,
			'=': token.LogAndAssign,
		},
	},
	'|': tokens{
		0:   token.BitOr,
		'=': token.BitOrAssign,
		'|': tokens{
			0:   token.LogOr,
			'=': token.LogOrAssign,
		},
	},
	'?': tokens{
		0: token.Question,
		'?': tokens{
			0:   token.LogNull,
			'=': token.LogNullAssign,
		},
	},
	'=': tokens{
		0:   token.Assign,
		'=': token.Eq,
	},
	'!': tokens{
		0:   token.LogNot,
		'=': token.Neq,
	},
	'<': tokens{
		0:   token.Lt,
		'=': token.Lte,
		'<': tokens{
			0:   token.BitShl,
			'=': token.BitShlAssign,
		},
	},
	'>': tokens{
		0:   token.Gt,
		'=': token.Gte,
		'>': tokens{
			0:   token.BitShr,
			'=': token.BitShrAssign,
		},
	},
	'^': tokens{
		0:   token.BitXor,
		'=': token.BitXorAssign,
	},
	'~': token.BitNot,
}
