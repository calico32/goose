package token

import (
	"fmt"
)

type Position struct {
	Filename string
	Offset   int // offset, starting at 0
	Line     int
	Column   int
}

func (pos *Position) IsValid() bool { return pos.Line > 0 }

func (pos Position) String() string {
	s := pos.Filename
	if pos.IsValid() {
		if s != "" {
			s += ":"
		}
		s += fmt.Sprintf("%d", pos.Line)
		if pos.Column != 0 {
			s += fmt.Sprintf(":%d", pos.Column)
		}
	} else {
		s += "<none>"
	}
	if s == "" {
		s = "-"
	}
	return s
}

type Pos int

const NoPos Pos = 0

func (p Pos) IsValid() bool {
	return p != NoPos
}
