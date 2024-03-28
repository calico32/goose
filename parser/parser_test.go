package parser_test

import (
	"github.com/calico32/goose/parser"
	"github.com/calico32/goose/scanner"
	"github.com/calico32/goose/token"
	. "github.com/onsi/ginkgo/v2"
)

func prepareParser(src string) *parser.Parser {
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("test", fset.Base(), len(src))
	s.Init(file, []byte(src), func(pos token.Position, msg string) {
		Fail(msg)
	})

	var p parser.Parser
	p.Init(fset, "test", []byte(src), nil)
	return &p
}
