package token

import (
	"strconv"
	"unicode"
)

type Token int

const (
	Illegal Token = iota
	EOF
	Comment

	literal_start
	Ident
	Int   // int64
	Float // float64
	StringStart
	StringMid
	StringInterpIdent
	StringInterpExprStart
	StringInterpExprEnd
	StringEnd
	literal_end

	operator_start
	Add
	Sub
	Mul
	Quo
	Rem
	Pow

	Assign
	AddAssign
	SubAssign
	MulAssign
	QuoAssign
	RemAssign
	PowAssign
	Inc
	Dec

	LogAnd
	LogOr
	LogNot

	Eq
	Neq
	Lt
	Gt
	Lte
	Gte
	LParen
	RParen
	LBracket
	RBracket
	LBrace
	RBrace
	Colon
	Semi
	Comma
	Period
	EOL
	operator_end

	keyword_start
	If
	Else
	Repeat
	Forever
	Times
	While
	For
	In
	Break
	Continue
	Null
	End
	Func
	Memo
	Const
	Let
	Return
	Include
	keyword_end

	None
)

var tokens = [...]string{
	Illegal: "<ILLEGAL>",
	EOF:     "<EOF>",
	Comment: "<COMMENT>",

	Ident:                 "<IDENT>",
	Int:                   "<INT>",
	Float:                 "<FLOAT>",
	StringStart:           "<S_START>",
	StringMid:             "<S_MID>",
	StringInterpIdent:     "<S_IDENT>",
	StringInterpExprStart: "<S_EXPRS>",
	StringInterpExprEnd:   "<S_EXPRE>",
	StringEnd:             "<S_END>",

	Add: "+",
	Sub: "-",
	Mul: "*",
	Quo: "/",
	Rem: "%",
	Pow: "**",

	Assign:    "=",
	AddAssign: "+=",
	SubAssign: "-=",
	MulAssign: "*=",
	QuoAssign: "/=",
	RemAssign: "%=",
	PowAssign: "**=",
	Inc:       "++",
	Dec:       "--",

	LogAnd: "&&",
	LogOr:  "||",
	LogNot: "!",

	Eq:  "==",
	Neq: "!=",
	Lt:  "<",
	Gt:  ">",
	Lte: "<=",
	Gte: ">=",

	LParen:   "(",
	RParen:   ")",
	LBracket: "[",
	RBracket: "]",
	LBrace:   "{",
	RBrace:   "}",
	Colon:    ":",
	Semi:     ";",
	Comma:    ",",
	Period:   ".",
	EOL:      "\n",

	If:       "if",
	Else:     "else",
	Repeat:   "repeat",
	Forever:  "forever",
	Times:    "times",
	While:    "while",
	For:      "for",
	In:       "in",
	Break:    "break",
	Continue: "continue",
	Null:     "null",
	End:      "end",
	Func:     "fn",
	Memo:     "memo",
	Const:    "const",
	Let:      "let",
	Return:   "return",
	Include:  "include",
}

func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

const (
	PrecedenceLowest  = 0
	PrecedenceUnary   = 6
	PrecedenceHighest = 7
)

func (tok Token) Precedence() int {
	switch tok {
	case LogOr:
		return 1
	case LogAnd:
		return 2
	case Eq, Neq, Lt, Gt, Lte, Gte:
		return 3
	case Add, Sub:
		return 4
	case Mul, Quo, Rem, Pow:
		return 5
	}
	return PrecedenceLowest
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token, keyword_end-(keyword_start+1))
	for i := keyword_start + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
}

func Lookup(ident string) Token {
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}
	return Ident
}

func (tok Token) IsLiteral() bool {
	return literal_start < tok && tok < literal_end
}

func (tok Token) IsOperator() bool {
	return operator_start < tok && tok < operator_end
}

func (tok Token) IsKeyword() bool {
	return keyword_start < tok && tok < keyword_end
}

func IsKeyword(tok string) bool {
	_, is_keyword := keywords[tok]
	return is_keyword
}

func IsIdentifier(name string) bool {
	if name == "" || IsKeyword(name) {
		return false
	}

	for i, c := range name {
		if !unicode.IsLetter(c) && c != '_' && (i == 0 || !unicode.IsDigit(c)) {
			return false
		}
	}

	return true
}
