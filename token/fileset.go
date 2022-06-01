package token

import (
	"fmt"
	"sync"
)

type FileSet struct {
	mutex sync.Mutex // protects the file set
	base  int
	files []*File
	last  *File
}

func (s *FileSet) Mutex() *sync.Mutex { return &s.mutex }

func NewFileSet() *FileSet {
	return &FileSet{
		base: 1, // 0 == NoPos
	}
}

func (s *FileSet) Base() int {
	defer un(lock(s))
	return s.base
}

func (s *FileSet) AddFile(filename string, base, size int) *File {
	defer un(lock(s))
	if base < 0 {
		base = s.base
	}
	if base < s.base {
		panic(fmt.Sprintf("token.FileSet: base (%d) < token.FileSet.base (%d)", base, s.base))
	}
	if size < 0 {
		panic(fmt.Sprintf("token.FileSet: size (%d) < 0", size))
	}
	f := &File{
		set:   s,
		name:  filename,
		base:  base,
		size:  size,
		lines: []int{0},
	}
	base += size + 1
	if base < 0 {
		panic("token.FileSet: offset overflow (>2G of source code in file set)")
	}
	s.base = base
	s.files = append(s.files, f)
	s.last = f
	return f
}

func (s *FileSet) Position(pos Pos) Position {
	var file *File
	for _, f := range s.files {
		if f.base <= int(pos) && int(pos) <= f.base+f.size {
			file = f
			break
		}
	}

	if file == nil {
		return Position{Filename: "unknown", Line: 1, Column: 1}
	}

	return file.Position(pos)
}
