package token

import (
	"fmt"
	"sync"
)

type Lockable interface {
	Mutex() *sync.Mutex
}

type File struct {
	set  *FileSet
	name string
	base int
	size int

	mutex sync.Mutex
	lines []int
}

func (s *File) Mutex() *sync.Mutex { return &s.mutex }

func un(mutex *sync.Mutex) {
	mutex.Unlock()
}

func lock(mutex *sync.Mutex) *sync.Mutex {
	mutex.Lock()
	return mutex
}

// Name returns the file name of file f as registered with AddFile.
func (f *File) Name() string {
	return f.name
}

// Base returns the base offset of file f as registered with AddFile.
func (f *File) Base() int {
	return f.base
}

// Size returns the size of file f as registered with AddFile.
func (f *File) Size() int {
	return f.size
}

// Pos returns the Pos value for the given file offset;
// the offset must be <= f.Size().
// f.Pos(f.Offset(p)) == p.
func (f *File) Pos(offset int) Pos {
	if offset > f.size {
		panic(fmt.Sprintf("invalid file offset %d (should be <= %d)", offset, f.size))
	}
	return Pos(f.base + offset)
}

// Offset returns the offset for the given file position p;
// p must be a valid Pos value in that file.
// f.Offset(f.Pos(offset)) == offset.
func (f *File) Offset(p Pos) int {
	if int(p) < f.base || int(p) > f.base+f.size {
		panic(fmt.Sprintf("invalid Pos value %d (should be in [%d, %d])", p, f.base, f.base+f.size))
	}
	return int(p) - f.base
}

// SetLinesForContent sets the line offsets for the given file content.
// It ignores position-altering //line comments.
func (f *File) SetLinesForContent(content []byte) {
	var lines []int
	line := 0
	for offset, b := range content {
		if line >= 0 {
			lines = append(lines, line)
		}
		line = -1
		if b == '\n' {
			line = offset + 1
		}
	}

	defer un(lock(&f.mutex))
	f.lines = lines
}

func (f *File) AddLine(offset int) {
	defer un(lock(&f.mutex))

	numLines := len(f.lines)

	if offset > f.size {
		panic(fmt.Sprintf("line offset %d > file size %d", offset, f.size))
	}
	if numLines != 0 && offset <= f.lines[numLines-1] {
		panic(fmt.Sprintf("line offset %d <= previous line offset %d", offset, f.lines[numLines-1]))
	}

	f.lines = append(f.lines, offset)
}

func (f *File) LineCount() int {
	defer un(lock(&f.mutex))
	return len(f.lines)
}

func (f *File) LineStart(line int) Pos {
	if line < 1 {
		panic(fmt.Sprintf("invalid line number %d (should be >= 1)", line))
	}
	defer un(lock(&f.mutex))
	if line > len(f.lines) {
		panic(fmt.Sprintf("invalid line number %d (should be < %d)", line, len(f.lines)))
	}
	return Pos(f.base + f.lines[line-1])
}

// Line returns the line number for the given file position p;
// p must be a Pos value in that file or NoPos.
func (f *File) Line(p Pos) int {
	return f.Position(p).Line
}

func (f *File) position(p Pos) (pos Position) {
	offset := int(p) - f.base
	pos.Offset = offset
	pos.Filename, pos.Line, pos.Column = f.unpack(offset)
	return
}

// PositionFor returns the Position value for the given file position p.
// If adjusted is set, the position may be adjusted by position-altering
// //line comments; otherwise those comments are ignored.
// p must be a Pos value in f or NoPos.
func (f *File) Position(p Pos) (pos Position) {
	if p != NoPos {
		if int(p) < f.base || int(p) > f.base+f.size {
			panic(fmt.Sprintf("invalid Pos value %d (should be in [%d, %d])", p, f.base, f.base+f.size))
		}
		pos = f.position(p)
	}
	return
}

// unpack returns the filename and line and column number for a file offset.
// If adjusted is set, unpack will return the filename and line information
// possibly adjusted by //line comments; otherwise those comments are ignored.
func (f *File) unpack(offset int) (filename string, line, column int) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	filename = f.name
	if i := searchInts(f.lines, offset); i >= 0 {
		line, column = i+1, offset-f.lines[i]+1
	}

	return
}

func searchInts(a []int, x int) int {
	// This function body is a manually inlined version of:
	// return sort.Search(len(a), func(i int) bool { return a[i] > x }) - 1

	i, j := 0, len(a)
	for i < j {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h
		// i ≤ h < j
		if a[h] <= x {
			i = h + 1
		} else {
			j = h
		}
	}
	return i - 1
}
