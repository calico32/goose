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

	LiteralStart
	Ident
	Int   // int64
	Float // float64
	StringStart
	StringMid
	StringInterpIdent
	StringInterpExprStart
	StringInterpExprEnd
	StringEnd
	LiteralEnd

	OperatorStart
	Assign
	OverloadAllowedStart
	Add
	Sub
	Mul
	Quo
	Rem
	Pow
	Gt
	Lt
	Lte
	Gte
	Inc
	Dec
	Question
	BitShl
	BitShr
	BitNot
	BitAnd
	BitOr
	BitXor
	Eq
	Neq
	Spaceship
	OverloadAllowedEnd
	AddAssign
	SubAssign
	MulAssign
	QuoAssign
	RemAssign
	PowAssign
	LogAnd
	LogOr
	LogNot
	LogNull
	LogAndAssign
	LogOrAssign
	LogNullAssign
	BitAndAssign
	BitOrAssign
	BitXorAssign
	BitShlAssign
	BitShrAssign
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
	Arrow
	Ellipsis
	HashLBracket
	Bind
	EOL
	OperatorEnd

	KeywordStart
	Null
	True
	False
	Let
	Const
	Symbol
	If
	Then
	Else
	Repeat
	While
	Forever
	Times
	For
	In
	Break
	Continue
	Func
	End
	Return
	Memo
	Import
	Export
	As
	Show
	Use
	Generator
	Yield
	To
	Step
	Struct
	Init
	Operator
	Try
	Catch
	Finally
	Throw
	Do
	Async
	Await
	Native
	KeywordEnd

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

	Assign:        "=",
	Add:           "+",
	Sub:           "-",
	Mul:           "*",
	Quo:           "/",
	Rem:           "%",
	Pow:           "**",
	AddAssign:     "+=",
	SubAssign:     "-=",
	MulAssign:     "*=",
	QuoAssign:     "/=",
	RemAssign:     "%=",
	PowAssign:     "**=",
	Inc:           "++",
	Dec:           "--",
	LogAnd:        "&&",
	LogOr:         "||",
	LogNot:        "!",
	LogNull:       "??",
	LogAndAssign:  "&&=",
	LogOrAssign:   "||=",
	LogNullAssign: "??=",
	Question:      "?",
	Eq:            "==",
	Neq:           "!=",
	Lt:            "<",
	Gt:            ">",
	Lte:           "<=",
	Gte:           ">=",
	BitNot:        "~",
	BitAnd:        "&",
	BitOr:         "|",
	BitXor:        "^",
	BitShl:        "<<",
	BitShr:        ">>",
	BitAndAssign:  "&=",
	BitOrAssign:   "|=",
	BitXorAssign:  "^=",
	BitShlAssign:  "<<=",
	BitShrAssign:  ">>=",

	LParen:       "(",
	RParen:       ")",
	LBracket:     "[",
	RBracket:     "]",
	LBrace:       "{",
	RBrace:       "}",
	Colon:        ":",
	Semi:         ";",
	Comma:        ",",
	Period:       ".",
	Arrow:        "->",
	Ellipsis:     "...",
	HashLBracket: "#[",
	Bind:         "::",
	EOL:          "<EOL>",

	Null:      "null",
	True:      "true",
	False:     "false",
	Let:       "let",
	Const:     "const",
	Symbol:    "symbol",
	If:        "if",
	Then:      "then",
	Else:      "else",
	Repeat:    "repeat",
	While:     "while",
	Forever:   "forever",
	Times:     "times",
	For:       "for",
	In:        "in",
	Break:     "break",
	Continue:  "continue",
	Func:      "fn",
	End:       "end",
	Return:    "return",
	Memo:      "memo",
	Import:    "import",
	Export:    "export",
	As:        "as",
	Show:      "show",
	Use:       "use",
	Generator: "generator",
	Yield:     "yield",
	To:        "to",
	Step:      "step",
	Struct:    "struct",
	Init:      "init",
	Operator:  "operator",
	Try:       "try",
	Catch:     "catch",
	Finally:   "finally",
	Throw:     "throw",
	Do:        "do",
	Async:     "async",
	Await:     "await",
	Native:    "native",
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
	PrecedenceUnary   = 9
	PrecedenceHighest = 10
)

func (tok Token) Precedence() int {
	switch tok {
	case Arrow:
		return 1
	case LogNull:
		return 2
	case LogOr:
		return 3
	case LogAnd:
		return 4
	case Eq, Neq, Lt, Gt, Lte, Gte:
		return 5
	case Add, Sub:
		return 6
	case Mul, Quo, Rem:
		return 7
	case Pow:
		return 8
	}
	return PrecedenceLowest
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token, KeywordEnd-(KeywordStart+1))
	for i := KeywordStart + 1; i < KeywordEnd; i++ {
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
	return LiteralStart < tok && tok < LiteralEnd
}

func (tok Token) IsOperator() bool {
	return OperatorStart < tok && tok < OperatorEnd
}

func (tok Token) IsKeyword() bool {
	return KeywordStart < tok && tok < KeywordEnd
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

func IsPostfixOperator(tok Token) bool {
	return tok == Inc ||
		tok == Dec ||
		tok == Question
}
