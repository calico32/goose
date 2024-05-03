//go:build release

package main

import "embed"

//go:embed static
var staticFS embed.FS

//go:embed all:views
var tmplFS embed.FS
