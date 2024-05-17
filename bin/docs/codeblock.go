package main

import (
	"crypto/md5"
	_ "embed"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"strings"
)

//go:embed tmp/highlight.js
var highlightJs []byte
var jsBinary string
var codeBlocks = NewMutexMap[string, string]()

func init() {
	// look for bun or node on path
	bunPath, err := exec.LookPath("bun")
	if err != nil {
		panic("bun not found on path")
	}
	jsBinary = bunPath
}

func codeBlock(filename, code string) template.HTML {
	// unfortunately, no good way to do textmate highlighting in go
	// call out to TS for now

	code = stripIndent(strings.TrimPrefix(strings.TrimSuffix(strings.TrimSpace(code), "</pre>"), "<pre>"))

	var highlighted string
	// check if the existing block has already been highlighted
	if existing, ok := codeBlocks.Load(code); ok {
		highlighted = existing
	} else {
		// drop the highlight.js script onto disk
		highlightPath := fmt.Sprintf("/tmp/goose-docs-highlight-%x.js", md5.Sum([]byte(highlightJs)))
		if _, err := os.Stat(highlightPath); err != nil {
			os.WriteFile(highlightPath, highlightJs, 0755)
		}

		// run the highlight binary
		cmd := exec.Command(jsBinary, highlightPath)
		cmd.Stdin = strings.NewReader(code)
		out, err := cmd.Output()
		if err != nil {
			panic(err)
		}

		// save the code block
		codeBlocks.Store(code, string(out))
		highlighted = string(out)
	}

	return template.HTML(fmt.Sprintf(`
		<div class="code-block">
			<div class="code-header">%s</div>
			%s
		</div>
	`, filename, highlighted))
}

func stripIndent(s string) string {
	lines := strings.Split(s, "\n")
	if len(lines) == 0 {
		return ""
	}
	// if the first line is empty, skip it
	if len(strings.TrimSpace(lines[0])) == 0 {
		lines = lines[1:]
	}
	indent := len(lines[0]) - len(strings.TrimLeft(lines[0], " "))
	for i, line := range lines {
		if len(line) >= indent {
			lines[i] = line[indent:]
		}
	}
	// if the last line is empty, skip it
	if len(strings.TrimSpace(lines[len(lines)-1])) == 0 {
		lines = lines[:len(lines)-1]
	}
	return strings.Join(lines, "\n")
}
