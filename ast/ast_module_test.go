package ast

import (
	"testing"

	"github.com/calico32/goose/token"
)

func TestModuleName(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		// global module names
		"discord":                  "discord",
		"discord/":                 "", // trailing slash is not allowed
		"discord/commands":         "commands",
		"discord/commands.goose":   "commands", // if a file ends in .goose, it is removed from the module name
		"discord/commands/":        "",
		"discord/commands/abc":     "abc",
		"discord/commands/abc/":    "",
		"discord.org":              "discordorg",
		"discord.goose":            "discord",
		"discord.goose/commands":   "commands",
		"discord.goose/commands/":  "",
		"discord_commands":         "discord_commands",
		"_discord":                 "_discord",
		"discord_":                 "discord_",
		"a-bunch[of]illegal+chars": "abunchofillegalchars",
		"0123456789":               "", // module names must start with a letter

		"std":    "std",
		"std/":   "",
		"std/io": "io",
		// reserved
		"index": "",
		"_":     "",

		// relative imports
		"./discord":                          "discord",
		"./discord.commands":                 "discordcommands", // periods removed from module name
		"./discord.goose":                    "discord",         // if a file ends in .goose, it is removed from the module name
		"./discord.commands.goose":           "discordcommands",
		"./discord/":                         "",
		"./discord/index.goose":              "discord", // index.goose is treated as its parent directory
		"./discord/commands/index.goose":     "commands",
		"./discord/commands/abc.goose":       "abc",
		"./discord/commands/abc.foo.goose":   "abcfoo",
		"./discord/commands/abc/index.goose": "abc",
		"../discord.goose":                   "discord",
		"../discord/commands.goose":          "commands",
		"/path/to/discord.goose":             "discord",
		"/path/to/discord/commands.goose":    "commands",
		".":                                  "", // no module name for current directory
		"..":                                 "", // no module name
		"/":                                  "", // no module name
		"./.././../.././//discord.goose":     "discord",
		"./.././discord/../discord.goose":    "discord",
		"////discord.goose":                  "discord",
		"./discord/123.goose":                "",
		"./discord/123":                      "",
		"./discord/123/a.goose":              "a",
		"./discord/commands1234.goose":       "commands1234",
	}

	// expect empty string for all keywords
	for i := token.KeywordStart + 1; i < token.KeywordEnd; i++ {
		tests[i.String()] = ""
	}

	for name, expected := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := ModuleName(name)

			if expected == "" && err == nil {
				t.Errorf("expected error, got nil")
			} else if actual != expected {
				t.Errorf("expected '%s', got '%s'", expected, actual)
			}
		})
	}
}

func TestModulePath(t *testing.T) {
	t.Parallel()

	root := "/path/to/gooseroot"
	cwd := "/path/to/project"

	tests := map[string][]string{
		// global module names resolve to gooseroot/pkg/<name>
		"discord":                  {root + "/pkg/discord"}, // top level names do not get resolved to a .goose file
		"discord/":                 {},
		"discord/commands":         {root + "/pkg/discord/commands", root + "/pkg/discord/commands.goose"}, // look for .goose files after directories
		"discord/commands.goose":   {root + "/pkg/discord/commands.goose"},                                 // .goose extensions do not get duplicated
		"discord/commands/":        {},
		"discord/commands/abc":     {root + "/pkg/discord/commands/abc", root + "/pkg/discord/commands/abc.goose"},
		"discord/commands/abc/":    {},
		"discord.org":              {root + "/pkg/discord.org"},
		"discord.goose":            {root + "/pkg/discord.goose"},
		"discord.goose/commands":   {root + "/pkg/discord.goose/commands", root + "/pkg/discord.goose/commands.goose"},
		"discord.goose/commands/":  {},
		"discord_commands":         {root + "/pkg/discord_commands"},
		"_discord":                 {root + "/pkg/_discord"},
		"discord_":                 {root + "/pkg/discord_"},
		"a-bunch[of]illegal+chars": {root + "/pkg/a-bunch[of]illegal+chars"},
		"0123456789":               {}, // module names must start with a letter

		"std":    {root + "/std"}, // std is resolved to gooseroot/std
		"std/":   {},
		"std/io": {root + "/std/io"}, // top level std modules do not get resolved to a .goose file
		// reserved
		"index": {},
		"_":     {},

		// relative imports
		"./discord":                       {cwd + "/discord", cwd + "/discord.goose"},
		"./discord.commands":              {cwd + "/discord.commands", cwd + "/discord.commands.goose"},
		"./discord.goose":                 {cwd + "/discord.goose"},
		"./discord/":                      {},
		"./discord/index":                 {cwd + "/discord/index", cwd + "/discord/index.goose"},
		"./discord/commands":              {cwd + "/discord/commands", cwd + "/discord/commands.goose"},
		"./discord/commands/":             {},
		"./discord/commands/abc":          {cwd + "/discord/commands/abc", cwd + "/discord/commands/abc.goose"},
		"./discord/commands/abc/":         {},
		"../discord.goose":                {"/path/to/discord.goose"},
		"../discord/commands.goose":       {"/path/discord/commands.goose"},
		"/path/to/discord.goose":          {"/path/to/discord.goose"},
		"/path/to/discord/commands.goose": {"/path/to/discord/commands.goose"},
		".":                               {cwd},
		"..":                              {"/path/to"},
		"/":                               {"/"},
		"./.././../.././//discord.goose":  {"/discord.goose"},
		"./.././discord/../discord.goose": {"/path/to/discord.goose"},
		"////discord.goose":               {"/discord.goose"},
		"./discord/123.goose":             {cwd + "/discord/123.goose"},
		"./discord/123":                   {cwd + "/discord/123", cwd + "/discord/123.goose"},
		"./discord/123/a.goose":           {cwd + "/discord/123/a.goose"},
		"./discord/commands1234.goose":    {cwd + "/discord/commands1234.goose"},
	}

	return

	for name, expected := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := ModuleName(name)

			if len(expected) == 0 && err == nil {
				t.Errorf("expected error, got nil")
			} else {
				for _, path := range expected {
					if path != actual {
						t.Errorf("resolution of '%s' failed, expected '%s', got '%s'", name, path, actual)
					}
				}
			}
		})
	}
}

func ResolveModule(specifier, gooseroot, cwd string) ([]string, error) {
	return []string{}, nil
}
