package types

type StdlibDoc struct {
	Name        string
	Desc        string
	Description string
}

type BuiltinDoc struct {
	// Short name of the builtin, like `len` or `int.parse`.
	Name string
	// Label for the builtin, like `len(x)` or `int.parse(str)`.
	Label string
	// Full typed signature for the builtin, like `len(x: array | string) -> int` or `int.parse(s: str) -> int`.
	Signature string
	// Desc is a short description of what the builtin does. It should be a single sentence.
	Desc string
	// Description of the builtin. This can be a longer description of the builtin, including examples and more detailed information.
	Description string
	// Examples of how to use the builtin. Each example should include a code snippet and a caption.
	Examples []CodeSnippet
}

type CodeSnippet struct {
	// The language of the code snippet, like `go` or `python`. Defaults to `goose`.
	Language string
	Content  string
	Caption  string
}
