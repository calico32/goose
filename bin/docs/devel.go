//go:build !release

package main

import "embed"

var staticFS embed.FS

var tmplFS embed.FS
