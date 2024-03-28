package lib

import (
	"embed"

	std_canvas "github.com/calico32/goose/lib/std/canvas"
	std_collections "github.com/calico32/goose/lib/std/collections"
	std_crypto "github.com/calico32/goose/lib/std/crypto"
	std_fs "github.com/calico32/goose/lib/std/fs"
	std_json "github.com/calico32/goose/lib/std/json"
	std_math "github.com/calico32/goose/lib/std/math"
	std_platform "github.com/calico32/goose/lib/std/platform"
	std_random "github.com/calico32/goose/lib/std/random"
	std_readline "github.com/calico32/goose/lib/std/readline"
	std_time "github.com/calico32/goose/lib/std/time"
	"github.com/calico32/goose/lib/types"
)

//go:embed std/*
var Stdlib embed.FS

var StdlibDocs = []types.StdlibDoc{
	std_canvas.Doc,
	std_collections.Doc,
	std_crypto.Doc,
	std_fs.Doc,
	std_json.Doc,
	// std_language.Doc,
	std_math.Doc,
	std_platform.Doc,
	std_random.Doc,
	std_readline.Doc,
	std_time.Doc,
}
