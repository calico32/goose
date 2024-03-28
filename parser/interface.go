// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains the exported entry points for invoking the parser.

package parser

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/calico32/goose/ast"
	"github.com/calico32/goose/lib"
	"github.com/calico32/goose/token"
)

// If src != nil, readSource converts src to a []byte if possible;
// otherwise it returns an error. If src == nil, readSource returns
// the result of reading the file specified by filename.
func readSource(specifier string, src any) ([]byte, error) {
	if src != nil {
		switch s := src.(type) {
		case string:
			return []byte(s), nil
		case []byte:
			return s, nil
		case *bytes.Buffer:
			// is io.Reader, but src is already available in []byte form
			if s != nil {
				return s.Bytes(), nil
			}
		case io.Reader:
			return io.ReadAll(s)
		}
		return nil, errors.New("invalid source")
	}

	parts := strings.Split(specifier, ":")
	var scheme string
	var path string
	if len(parts) == 1 {
		scheme = "file"
		path = parts[0]
	} else {
		scheme = parts[0]
		path = strings.Join(parts[1:], ":")
	}

	switch scheme {
	case "file":
		return os.ReadFile(path)
	case "pkg":
		// TODO: implement
		return nil, errors.New("pkg scheme not implemented")
	case "std":
		return lib.Stdlib.ReadFile(filepath.Join("std", path))
	default:
		return nil, errors.New("invalid scheme")
	}
}

// ParseFile parses the source code of a single Goose source file and returns
// the corresponding ast.File node. The source code may be provided via
// the filename of the source file, or via the src parameter.
//
// If src != nil, ParseFile parses the source from src and the filename is
// only used when recording position information. The type of the argument
// for the src parameter must be string, []byte, or io.Reader.
// If src == nil, ParseFile parses the file specified by filename.
//
// The mode parameter controls the amount of source text parsed and other
// optional parser functionality. If the SkipObjectResolution mode bit is set,
// the object resolution phase of parsing will be skipped, causing File.Scope,
// File.Unresolved, and all Ident.Obj fields to be nil.
//
// Position information is recorded in the file set fset, which must not be
// nil.
//
// If the source couldn't be read, the returned AST is nil and the error
// indicates the specific failure. If the source was read but syntax
// errors were found, the result is a partial AST (with ast.Bad* nodes
// representing the fragments of erroneous source code). Multiple errors
// are returned via a scanner.ErrorList which is sorted by source position.
func ParseFile(fset *token.FileSet, specifier string, src any, trace io.Writer) (f *ast.Module, err error) {
	if fset == nil {
		panic("parser.ParseFile: no token.FileSet provided (fset == nil)")
	}

	// get source
	text, err := readSource(specifier, src)
	if err != nil {
		return nil, err
	}

	var p Parser
	// defer func() {
	// 	if e := recover(); e != nil {
	// 		// set result values
	// 		if f == nil {
	// 			// source is not a valid Go source file - satisfy
	// 			// ParseFile API and return a valid (but) empty
	// 			// *ast.File
	// 			f = &ast.Module{
	// 				// Name:  new(ast.Ident),
	// 				// Scope: ast.NewScope(nil),
	// 			}
	// 		}

	// 		p.errors.Sort()
	// 		err = p.errors.Err()
	// 		panic(e)
	// 	}
	// }()

	// parse source
	p.Init(fset, specifier, text, trace)
	f = p.ParseFile()

	if len(p.errors) > 0 {
		p.errors.Sort()
		err = p.errors.Err()
	}

	return
}

// ParseExprFrom is a convenience function for parsing an expression.
// The arguments have the same meaning as for ParseFile, but the source must
// be a valid Goose (type or value) expression. Specifically, fset must not
// be nil.
//
// If the source couldn't be read, the returned AST is nil and the error
// indicates the specific failure. If the source was read but syntax
// errors were found, the result is a partial AST (with ast.Bad* nodes
// representing the fragments of erroneous source code). Multiple errors
// are returned via a scanner.ErrorList which is sorted by source position.
func ParseExprFrom(fset *token.FileSet, specifier string, src any, trace io.Writer) (expr ast.Expr, err error) {
	if fset == nil {
		panic("parser.ParseExprFrom: no token.FileSet provided (fset == nil)")
	}

	// get source
	text, err := readSource(specifier, src)
	if err != nil {
		return nil, err
	}

	var p Parser
	defer func() {
		if e := recover(); e != nil {
			p.errors.Sort()
			err = p.errors.Err()
			panic(e)
		}
	}()

	// parse expr
	p.Init(fset, specifier, text, trace)
	expr = p.ParseExpr()

	p.expect(token.EOF)

	return
}

// ParseExpr is a convenience function for obtaining the AST of an expression x.
// The position information recorded in the AST is undefined. The filename used
// in error messages is the empty string.
//
// If syntax errors were found, the result is a partial AST (with ast.Bad* nodes
// representing the fragments of erroneous source code). Multiple errors are
// returned via a scanner.ErrorList which is sorted by source position.
func ParseExpr(x string, trace io.Writer) (ast.Expr, error) {
	return ParseExprFrom(token.NewFileSet(), "", []byte(x), trace)
}
