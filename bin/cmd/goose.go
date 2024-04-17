package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/calico32/goose/compiler"
	"github.com/calico32/goose/interpreter"
	std_platform "github.com/calico32/goose/lib/std/platform"
	"github.com/calico32/goose/parser"
	"github.com/calico32/goose/scanner"
	"github.com/calico32/goose/token"
	"github.com/calico32/goose/validator"
)

type Token struct {
	Pos token.Pos   `json:"pos"`
	Tok token.Token `json:"tok"`
	Lit string      `json:"lit"`
}

var outputFile = flag.String("output", "", "Write output to `file`")
var help = flag.Bool("help", false, "Show help message")
var version = flag.Bool("version", false, "Show version")
var jsonOutput = flag.Bool("json", false, "Output JSON instead of text")

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("goose %s · built at %s\n", std_platform.Version, std_platform.BuildTime)
		fmt.Printf("running on %s/%s\n", std_platform.OS, std_platform.Arch)
		os.Exit(0)
	}

	if *help || flag.NArg() == 0 {
		fmt.Println("goose compiler v0.1.0")
		fmt.Println("Usage: goose [options] <command> [arguments]")
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("  run [file] [arguments]   Run a goose program")
		fmt.Println("  validate [file] 	        Check a goose program for errors")
		fmt.Println("  build [file]             Compile a goose program")
		fmt.Println("  scan [file]              Scan a goose program")
		fmt.Println("  parse [file]             Parse a goose program")
		fmt.Println()
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	var verb string
	var args []string

	if flag.NArg() < 2 && flag.Arg(0) != "scan" && flag.Arg(0) != "parse" && flag.Arg(0) != "validate" && flag.Arg(0) != "build" {
		verb = "run"
		args = flag.Args()
	} else if flag.NArg() < 2 {
		verb = flag.Arg(0)
		args = []string{}
	} else {
		verb = flag.Arg(0)
		args = flag.Args()[1:]
	}

	var outWriter *os.File
	var err error

	if *outputFile != "" {
		outWriter, err = os.Create(*outputFile)
		if err != nil {
			fmt.Println("Failed to open output file:", err)
			os.Exit(1)
		}
	} else {
		outWriter = os.Stdout
	}

	switch verb {
	case "build":
		if len(args) < 1 {
			fmt.Println("Missing module")
			os.Exit(1)
		}

		spec := args[0]

		if !strings.Contains(spec, ":") {
			// no scheme, assume file
			absPath, err := filepath.Abs(spec)
			if err != nil {
				panic(err)
			}
			spec = "file:" + absPath
		}

		fset := token.NewFileSet()

		f, err := parser.ParseFile(fset, spec, nil, nil)
		if err != nil {
			panic(err)
		}

		c, err := compiler.New(f, fset, false, os.Stdin, outWriter, os.Stderr)
		if err != nil {
			panic(err)
		}

		module := c.Compile()

		fmt.Println(module)
	case "validate":
		if len(args) < 1 {
			fmt.Println("Missing module")
			os.Exit(1)
		}

		spec := args[0]

		if !strings.Contains(spec, ":") {
			// no scheme, assume file
			absPath, err := filepath.Abs(spec)
			if err != nil {
				panic(err)
			}
			spec = "file:" + absPath
		}

		fset := token.NewFileSet()

		f, err := parser.ParseFile(fset, spec, nil, nil)
		if err != nil {
			panic(err)
		}

		v, err := validator.New(f, fset, false, os.Stdin, outWriter, os.Stderr)
		if err != nil {
			panic(err)
		}

		exitCode, err := v.Check()
		if err != nil {
			panic(err)
		}

		for _, d := range v.Diagnostics() {
			fmt.Printf("%s: %s\n", d.Severity.String(), d.Message)

			pos := strings.TrimPrefix(fset.Position(d.Node.Pos()).String(), "file:")
			fmt.Printf("\tat %s\n", pos)
		}

		os.Exit(exitCode)

	case "run":
		if len(args) < 1 {
			fmt.Println("Missing module")
			os.Exit(1)
		}

		spec := args[0]

		if !strings.Contains(spec, ":") {
			// no scheme, assume file
			absPath, err := filepath.Abs(spec)
			if err != nil {
				panic(err)
			}
			spec = "file:" + absPath
		}

		fset := token.NewFileSet()

		f, err := parser.ParseFile(fset, spec, nil, nil)
		if err != nil {
			panic(err)
		}

		i, err := interpreter.New(f, fset, false, os.Stdin, outWriter, os.Stderr)
		if err != nil {
			panic(err)
		}

		exitCode, err := i.Run()
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

		if *jsonOutput {
			enc := json.NewEncoder(outWriter)
			enc.SetIndent("", "")
			enc.Encode(tokens)
		} else {
			for _, t := range tokens {
				pos := f.Position(t.Pos)
				fmt.Fprintf(outWriter, "%-30s", pos.String()+": "+t.Tok.String())
				fmt.Fprint(outWriter, "\t\t| "+t.Lit+"\n")
			}
		}

	case "parse":
		if len(args) < 1 {
			fmt.Println("Missing file")
			os.Exit(1)
		}

		var writer io.Writer
		if !*jsonOutput {
			writer = outWriter
		}

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, args[0], nil, writer)
		if err != nil {
			panic(err)
		}

		if *jsonOutput {
			j := json.NewEncoder(outWriter)
			j.SetIndent("", "  ")
			err = j.Encode(f)
			if err != nil {
				panic(err)
			}
		}

		outWriter.Sync()

	default:
		fmt.Println("Unknown verb:", verb)
	}

}
