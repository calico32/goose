package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/calico32/goose/interpreter"
	"github.com/calico32/goose/parser"
	"github.com/calico32/goose/scanner"
	"github.com/calico32/goose/token"
)

type Token struct {
	Pos token.Pos   `json:"pos"`
	Tok token.Token `json:"tok"`
	Lit string      `json:"lit"`
}

func usage() {
	fmt.Println("Usage: goose [options] <command> [<args>]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -h, --help                          Show this help")
	fmt.Println("  -v, --version                       Show version")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  run <file> [<args>...]              Run a Goose program")
	fmt.Println("  -o <file>                           Write stdout to file")
	fmt.Println("  -e <file>                           Write stderr to file")
	fmt.Println()
	fmt.Println("  scan <file>                         Tokenize a Goose program")
	fmt.Println("  -o <file>                           Write tokens to file")
	fmt.Println("  -j	                                 Output JSON")
	fmt.Println()
	fmt.Println("  parse <file>                        Parse a Goose program")
	fmt.Println("  -o <file>                           Write stringified AST to file")
	fmt.Println("  -j                                  Output JSON")
	fmt.Println()
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	var showHelp, showVersion bool
	var outputFile, errorFile string
	var jsonOutput, gobOutput bool

	var verb string
	var args []string

	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-h", "--help":
			showHelp = true
		case "-v", "--version":
			showVersion = true
		case "-o":
			if i+1 >= len(os.Args) {
				fmt.Println("Missing output file")
				os.Exit(1)
			}
			outputFile = os.Args[i+1]
			i++
		case "-e":
			if i+1 >= len(os.Args) {
				fmt.Println("Missing error file")
				os.Exit(1)
			}
			errorFile = os.Args[i+1]
			i++
		case "-j":
			jsonOutput = true
		case "-g":
			gobOutput = true
		default:
			if strings.HasPrefix(os.Args[i], "-") {
				fmt.Println("Unknown option:", os.Args[i])
				os.Exit(1)
			}
			if verb == "" {
				verb = os.Args[i]
			} else {
				args = append(args, os.Args[i])
			}
		}
	}

	if showHelp {
		usage()
		os.Exit(0)
	}

	if showVersion {
		fmt.Println("goose version 0.1.0")
		os.Exit(0)
	}

	var outw *os.File
	var errw *os.File
	var err error

	if outputFile != "" {
		outw, err = os.Create(outputFile)
		if err != nil {
			fmt.Println("Failed to open output file:", err)
			os.Exit(1)
		}
	} else {
		outw = os.Stdout
	}

	if errorFile != "" {
		errw, err = os.Create(errorFile)
		if err != nil {
			fmt.Println("Failed to open error file:", err)
			os.Exit(1)
		}

	} else {
		errw = os.Stderr
	}

	switch verb {
	case "run":
		if len(args) < 1 {
			fmt.Println("Missing file")
			os.Exit(1)
		}

		fset := token.NewFileSet()
		absPath, err := filepath.Abs(args[0])
		if err != nil {
			panic(err)
		}
		f, err := parser.ParseFile(fset, absPath, nil, nil)
		if err != nil {
			panic(err)
		}

		exitCode, err := interpreter.RunFile(f, fset, false, true, outw, errw)
		if err != nil {
			panic(err)
		}

		os.Exit(exitCode)

	case "scan":
		if len(args) < 1 {
			fmt.Println("Missing file")
			os.Exit(1)
		}

		file, err := os.Open(args[0])
		if err != nil {
			panic(err)
		}
		defer file.Close()

		info, err := file.Stat()
		if err != nil {
			panic(err)
		}

		content := make([]byte, info.Size())
		_, err = file.Read(content)
		if err != nil {
			panic(err)
		}

		fset := token.NewFileSet()
		f := fset.AddFile(args[0], fset.Base(), len(content))
		s := scanner.Scanner{}
		s.Init(f, content, func(pos token.Position, msg string) {
			fmt.Printf("%s: %s\n", pos, msg)
		})

		var tokens []Token
		for {
			pos, tok, lit := s.Scan()
			if tok == token.EOF {
				break
			}
			tokens = append(tokens, Token{pos, tok, lit})
		}

		if jsonOutput {
			enc := json.NewEncoder(outw)
			enc.SetIndent("", "")
			enc.Encode(tokens)
		} else {
			for _, t := range tokens {
				pos := f.Position(t.Pos)
				fmt.Fprintf(outw, "%-30s", pos.String()+": "+t.Tok.String())
				fmt.Fprint(outw, "\t\t| "+t.Lit+"\n")
			}
		}

	case "parse":
		if len(args) < 1 {
			fmt.Println("Missing file")
			os.Exit(1)
		}

		var writer io.Writer
		if !jsonOutput && !gobOutput {
			writer = outw
		}

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, args[0], nil, writer)
		if err != nil {
			panic(err)
		}

		if jsonOutput {
			j := json.NewEncoder(outw)
			j.SetIndent("", "  ")
			err = j.Encode(f)
			if err != nil {
				panic(err)
			}
		}

		if gobOutput {
			// ast.InitGob()

			if outw == os.Stdout {
				fmt.Println("Cannot write gob to stdout")
				os.Exit(1)
			}

			enc := gob.NewEncoder(outw)
			err = enc.Encode(f)
			if err != nil {
				panic(err)
			}
		}

		outw.Sync()

	default:
		fmt.Println("Unknown verb:", verb)
	}

}
