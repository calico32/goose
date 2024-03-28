package std_platform

import (
	"runtime"
	"time"

	. "github.com/calico32/goose/interpreter/lib"
	"github.com/calico32/goose/lib/types"
)

var Doc = types.StdlibDoc{
	Name:        "platform",
	Description: "Information about the platform Goose is running on.",
}

var Version = "1.0.0"
var Runtime = "go"
var BuildTime = time.Now().Format(time.RFC3339)
var OS = runtime.GOOS
var Arch = runtime.GOARCH

var Index = map[string]Value{
	"C/os":        NewString(OS),
	"C/arch":      NewString(Arch),
	"C/version":   NewString(Version),
	"C/runtime":   NewString(Runtime),
	"C/buildTime": NewString(BuildTime),
}
