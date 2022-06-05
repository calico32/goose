package main

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"github.com/wiisportsresort/goose/interpreter"
	"github.com/wiisportsresort/goose/parser"
	"github.com/wiisportsresort/goose/scanner"
	"github.com/wiisportsresort/goose/token"
)

func onError(pos token.Position, msg string) {
	fmt.Printf("%s: %s\n", pos, msg)
}

type Token struct {
	pos token.Pos
	tok token.Token
	lit string
}

// ptrToString returns a string from WebAssembly compatible numeric
// types representing its pointer and length.
func ptrToString(ptr uint32, size uint32) (ret string) {
	// Here, we want to get a string represented by the ptr and size. If we
	// wanted a []byte, we'd use reflect.SliceHeader instead.
	strHdr := (*reflect.StringHeader)(unsafe.Pointer(&ret))
	strHdr.Data = uintptr(ptr)
	strHdr.Len = uintptr(size)
	return
}

// stringToPtr returns a pointer and size pair for the given string
// in a way that is compatible with WebAssembly numeric types.
func stringToPtr(s string) (uint32, uint32) {
	if len(s) == 0 {
		return 0, 0
	}
	buf := []byte(s)
	ptr := &buf[0]
	unsafePtr := uintptr(unsafe.Pointer(ptr))
	return uint32(unsafePtr), uint32(len(buf))
}

func ptrSize(ptr, size uint32) uint32 {
	if ptr == 0 && size == 0 {
		return 0
	}

	buf := make([]byte, 8)
	binary.LittleEndian.PutUint32(buf[0:4], ptr)
	binary.LittleEndian.PutUint32(buf[4:8], size)

	return uint32(uintptr(unsafe.Pointer(&buf[0])))
}

// log a message to the console using _log.
func log(message string) {
	ptr, size := stringToPtr(message)
	_log(ptr, size)
}

func _log(ptr uint32, size uint32)

func main() {}

//export tokenize
func tokenize(name, source string) (ret uint32) {
	fset := token.NewFileSet()
	f := fset.AddFile(name, fset.Base(), len(source))
	s := scanner.Scanner{}
	s.Init(f, []byte(source), onError)

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

	return ptrSize(stringToPtr(out.String()))
}

//export parse
func parse(name, source string) (ret uint32) {
	var traceWriter strings.Builder

	fset := token.NewFileSet()
	_, err := parser.ParseFile(fset, name, source, &traceWriter)
	if err != nil {
		return ptrSize(stringToPtr(err.Error()))
	}

	return ptrSize(stringToPtr(traceWriter.String()))
}

//export execute
func execute(name, source string) (ret uint32) {
	var traceWriter strings.Builder
	var stdout strings.Builder
	var stderr strings.Builder

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, name, source, &traceWriter)
	if err != nil {
		return ptrSize(stringToPtr(err.Error()))
	}

	exitCode, err := interpreter.RunFile(f, fset, false, false, &stdout, &stderr)
	if err != nil {
		return ptrSize(stringToPtr(err.Error()))
	}

	var j strings.Builder

	stdoutPtr, stdoutSize := stringToPtr(stdout.String())
	stderrPtr, stderrSize := stringToPtr(stderr.String())
	tracePtr, traceSize := stringToPtr(traceWriter.String())

	j.WriteString("{\"exitCode\":")
	j.WriteString(fmt.Sprintf("%d", exitCode))
	j.WriteString(",\"stdout\":[")
	j.WriteString(fmt.Sprintf("%d, %d", stdoutPtr, stdoutSize))
	j.WriteString("],\"stderr\":[")
	j.WriteString(fmt.Sprintf("%d, %d", stderrPtr, stderrSize))
	j.WriteString("],\"trace\":[")
	j.WriteString(fmt.Sprintf("%d, %d", tracePtr, traceSize))
	j.WriteString("]}")

	return ptrSize(stringToPtr(j.String()))
}
