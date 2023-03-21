package goose

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/calico32/goose/interpreter"
	"github.com/calico32/goose/parser"
	"github.com/calico32/goose/token"
)

func Run(path string) (exitCode int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Goose panicked: %v", r)
			exitCode = -1
		}
	}()

	fset := token.NewFileSet()

	abspath, err := filepath.Abs(path)
	f, err := parser.ParseFile(fset, abspath, nil, nil)
	if err != nil {
		return -1, err
	}

	exitCode, err = interpreter.RunFile(f, fset, false, false, os.Stdout, os.Stderr)
	if err != nil {
		return -1, err
	}

	return
}

func RunCode(source string) (exitCode int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Goose panicked: %v", r)
			exitCode = -1
		}
	}()

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "goose.RunCode", source, nil)
	if err != nil {
		panic(err)
	}

	exitCode, err = interpreter.RunFile(f, fset, false, false, os.Stdout, os.Stderr)
	if err != nil {
		panic(err)
	}

	return
}
