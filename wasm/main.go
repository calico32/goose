package main

import (
	"fmt"
	"strings"
	"syscall/js"

	"github.com/wiisportsresort/goose/interpreter"
	"github.com/wiisportsresort/goose/parser"
	"github.com/wiisportsresort/goose/scanner"
	"github.com/wiisportsresort/goose/token"
)

var done chan bool

func init() {
	done = make(chan bool)
}

func main() {
	global := js.Global()
	global.Set("Goose", map[string]any{
		"tokenize": js.FuncOf(tokenize),
		"parse":    js.FuncOf(parse),
		"execute":  js.FuncOf(execute),
	})
	<-done
}

type Token struct {
	pos token.Pos
	tok token.Token
	lit string
}

func onError(pos token.Position, msg string) {
	fmt.Printf("%s: %s\n", pos, msg)
}

func tokenize(_ js.Value, args []js.Value) (ret any) {
	defer func() {
		if r := recover(); r != nil {
			ret = fmt.Sprintf("panic: %v\n", r)
		}
	}()

	name := args[0].String()
	source := args[1].String()
	sourceBytes := []byte(source)

	fset := token.NewFileSet()
	f := fset.AddFile(name, fset.Base(), len(source))
	s := scanner.Scanner{}
	s.Init(f, sourceBytes, onError)

	var tokens []Token
	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		tokens = append(tokens, Token{pos, tok, lit})
	}

	var out strings.Builder
	for _, t := range tokens {
		pos := f.Position(t.pos)
		fmt.Fprintf(&out, "%-30s", pos.String()+": "+t.tok.String())
		fmt.Fprint(&out, "\t\t| "+t.lit+"\n")
	}

	return out.String()
}

func parse(_ js.Value, args []js.Value) (ret any) {
	defer func() {
		if r := recover(); r != nil {
			ret = map[string]any{
				"trace": fmt.Sprintf("panic: %v+\n", r),
			}
		}
	}()

	name := args[0].String()
	source := args[1].String()
	sourceBytes := []byte(source)

	var traceWriter strings.Builder

	fset := token.NewFileSet()
	_, err := parser.ParseFile(fset, name, sourceBytes, &traceWriter)
	if err != nil {
		panic(err)
	}

	return map[string]any{
		"trace": traceWriter.String(),
	}
}

func execute(_ js.Value, args []js.Value) (ret any) {
	defer func() {
		if r := recover(); r != nil {
			ret = fmt.Sprintf("panic: %v\n", r)
		}
	}()

	name := args[0].String()
	source := args[1].String()
	sourceBytes := []byte(source)

	var traceWriter strings.Builder
	var stdout strings.Builder
	var stderr strings.Builder

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, name, sourceBytes, &traceWriter)
	if err != nil {
		panic(err)
	}

	exitCode, err := interpreter.RunFile(f, fset, false, false, &stdout, &stderr)
	if err != nil {
		panic(err)
	}

	return map[string]any{
		"exitCode": exitCode,
		"stdout":   stdout.String(),
		"stderr":   stderr.String(),
		"trace":    traceWriter.String(),
	}
}
