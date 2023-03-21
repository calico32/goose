// Code generated for package lib by go-bindata DO NOT EDIT. (@generated)
// sources:
// std/_module.goose
// std/fs/_module.goose
// std/math/_module.goose
package lib

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _libStd_moduleGoose = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00")

func libStd_moduleGooseBytes() ([]byte, error) {
	return bindataRead(
		_libStd_moduleGoose,
		"lib/std/_module.goose",
	)
}

func libStd_moduleGoose() (*asset, error) {
	bytes, err := libStd_moduleGooseBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "lib/std/_module.goose", size: 0, mode: os.FileMode(420), modTime: time.Unix(1673274169, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _libStdFs_moduleGoose = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x8d\x31\x0a\xc4\x30\x0c\x04\x7b\xbd\x42\xe5\x19\xf4\x8e\xfb\x87\x40\x7b\x9c\xc1\x38\xc2\x88\x24\x10\xf2\xf7\x74\xb6\xaa\xd4\xb3\x3b\xd3\x35\xea\x0e\xfe\x75\x1e\x50\xfb\xd6\x86\x8f\x6b\xfc\x0b\x2d\x70\x8c\x1a\x98\x44\xd8\x34\x34\x73\x75\x47\xb7\x97\x81\xa1\x21\x19\x0a\x11\xe1\xf4\x6d\x04\x5f\xc4\xb3\x2b\x2b\x24\xc9\x29\xe9\x4e\x37\x3d\x01\x00\x00\xff\xff\x27\xae\x4a\x7d\xb0\x00\x00\x00")

func libStdFs_moduleGooseBytes() ([]byte, error) {
	return bindataRead(
		_libStdFs_moduleGoose,
		"lib/std/fs/_module.goose",
	)
}

func libStdFs_moduleGoose() (*asset, error) {
	bytes, err := libStdFs_moduleGooseBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "lib/std/fs/_module.goose", size: 176, mode: os.FileMode(420), modTime: time.Unix(1673356416, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _libStdMath_moduleGoose = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x64\x92\x51\x6e\x03\x21\x0c\x44\xff\x39\x05\x9f\x89\x34\x52\x92\x3d\x44\xa5\x1e\xc3\x49\xc9\x2e\x12\x01\xba\xd0\x76\xa3\xaa\x77\xaf\x6c\xd4\x26\x86\x9f\x95\xf6\x79\x6c\x60\x3c\x91\xaa\xff\x74\xf6\x1a\x6d\xf1\x71\xb7\xed\xcd\x03\x5c\x52\xd1\xa0\x52\xa7\xa0\xa1\x87\x86\x26\x1a\xbb\x2a\xc5\x69\x77\x87\x55\xb4\xf8\xb8\x0c\xe7\x2f\xc3\x05\x84\x3c\xa1\x35\x7d\xc4\x37\xad\xba\x86\x94\xd6\x6e\x94\xf3\x41\x93\xf2\xbe\xd6\x6e\x54\x48\xb3\xd6\x84\x34\x4f\x03\x39\x1d\x35\x72\x5b\xd6\x20\xa7\xaf\xdd\x06\x7b\x57\xb3\xe9\xdc\xd9\x52\xfc\xdc\xd9\x72\x63\x2f\xa5\xef\x89\xd1\x36\xb0\x4b\xa0\x5b\x66\x7a\xf3\x11\xac\x50\x07\x65\x5a\x8b\x7b\x8d\x95\x05\x67\x2a\x6e\x6f\x0e\x07\xdb\x95\x5f\x42\x22\xf5\xf8\x4b\x8a\xa5\xda\xec\xf5\xbf\xeb\xca\x4b\x57\x67\x0b\x27\x8d\x42\xec\x41\x9a\x27\xd7\x6b\x4e\xc7\x41\x74\x3a\x3a\x63\xdc\x96\xd3\x5a\xed\xb7\xb1\x9c\x06\x70\x00\xc0\x3b\x87\x04\x0d\x12\x2e\x48\x7c\xda\x77\x82\xa4\x46\x84\x8b\x28\x17\x18\xdb\x22\x81\x16\x03\xc8\xea\x21\x77\xe5\x5a\x48\x33\xe4\x52\x68\xa7\x82\xf7\x07\xde\x19\x57\xe9\x5c\x20\x9b\x79\x98\x8b\xe6\x37\x8c\xb1\xff\xde\xb6\x1f\x0f\xeb\xc0\xae\xb4\xe9\x3c\x31\xb6\xb1\x93\x83\xbc\x12\x7f\x2f\xfb\x31\xbf\x01\x00\x00\xff\xff\x7a\x57\xe8\x42\x67\x03\x00\x00")

func libStdMath_moduleGooseBytes() ([]byte, error) {
	return bindataRead(
		_libStdMath_moduleGoose,
		"lib/std/math/_module.goose",
	)
}

func libStdMath_moduleGoose() (*asset, error) {
	bytes, err := libStdMath_moduleGooseBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "lib/std/math/_module.goose", size: 871, mode: os.FileMode(420), modTime: time.Unix(1673541685, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"lib/std/_module.goose":      libStd_moduleGoose,
	"lib/std/fs/_module.goose":   libStdFs_moduleGoose,
	"lib/std/math/_module.goose": libStdMath_moduleGoose,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"lib": &bintree{nil, map[string]*bintree{
		"std": &bintree{nil, map[string]*bintree{
			"_module.goose": &bintree{libStd_moduleGoose, map[string]*bintree{}},
			"fs": &bintree{nil, map[string]*bintree{
				"_module.goose": &bintree{libStdFs_moduleGoose, map[string]*bintree{}},
			}},
			"math": &bintree{nil, map[string]*bintree{
				"_module.goose": &bintree{libStdMath_moduleGoose, map[string]*bintree{}},
			}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
