package interpreter

import (
	"os"
	"path/filepath"
)

func CreateGooseRoot(gooseRoot string) error {
	if err := os.MkdirAll(gooseRoot, 0755); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Join(gooseRoot, "pkg"), 0755); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Join(gooseRoot, "std"), 0755); err != nil {
		return err
	}

	// unpack stdlib
	// for _, name := range lib.Stdlib.ReadDir() {
	// 	if strings.HasPrefix(name, "lib/std/") {
	// 		data, err := lib.Asset(name)
	// 		if err != nil {
	// 			return err
	// 		}

	// 		path := filepath.Join(gooseRoot, "std", strings.TrimPrefix(name, "lib/std/"))

	// 		err = os.MkdirAll(filepath.Dir(path), os.FileMode(0755))
	// 		if err != nil {
	// 			return err
	// 		}
	// 		err = os.WriteFile(path, data, os.FileMode(0644))
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	return nil
}
