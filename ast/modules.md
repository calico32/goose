# Module conventions

## Definitions

- **Module**: A module is a file that contains a goose code. Modules are identified by their filesystem path.
- **Package**: A package is a directory that contains a `goose.toml` file.
- **Specifier**: A specifier is a string that identifies a module. It can be a relative path, an absolute path, or a package name.
- **`$GOOSEROOT`**: The root of the goose installation. This contains the standard library, binaries, and installed packages.

## Package naming

Package names should only use characters in the set `[a-z0-9_.-]`. The package name should be all lower case. 

Additionally, package names cannot be any of:
- keywords (`async`, `if`, `continue`)
- `_module`
- `_`

## Specifier rules

1. Specifiers must be valid UTF-8.
2. Specifiiers cannot end with `/`.

## Module resolution

1. If the specifier is a relative path (starts with `./` or `../`), resolve the module relative to the current module. If the specifier is an absolute path (starts with `/`), resolve the module relative to the root of the filesystem.
   1. The exact path given in the specifier is tried first.
      1. If the path is a directory, look for a `_module.goose` file in the directory.
         1. If it exists, use that file as the module.
         2. Otherwise, error; the path must be a file.
   2. Otherwise, if the path does not already end with `.goose`, try appending `.goose` to the path.
      1. If the path is a directory, look for a `_module.goose` file in the directory.
         1. If it exists, use that file as the module.
         2. Otherwise, error; the path must be a file.
   3. In all other cases, use the path as the module.
2. If the specifier begins with a package name, resolve the module relative to `$GOOSEROOT/packages`.
   1. The specifier is split on the first `/`. The first part is the package name, and the second part is the path within the package. 
   2. If the specifier is only the package name, look for `$GOOSEROOT/packages/<package name>/goose.toml`. If it exists, use the `main` field as the path.
      1. If the `main` field is not set, use `./_module.goose`.
      2. Treat it as a specifier relative to `$GOOSEROOT/packages/<package name>`. 
   3. Otherwise, treat the rest of the specifier as a relative path and resolve it relative to `$GOOSEROOT/packages/<package name>`.
   4. Follow the same rules as above for resolving the module.

### Examples
```js
// current module: "/home/user/project/foo.goose"
"./bar.goose"                   -> "/home/user/project/bar.goose/_module.goose", "/home/user/project/bar.goose"
"./bar"                         -> "/home/user/project/bar/_module.goose", "/home/user/project/bar", "/home/user/project/bar.goose"
"../project2/bar.goose"         -> "/home/user/project2/bar.goose/_module.goose", "/home/user/project2/bar.goose"
"/home/user/project2/bar.goose" -> "/home/user/project2/bar.goose/_module.goose", "/home/user/project2/bar.goose"
"discord"                       -> "$GOOSEROOT/packages/discord/goose.toml" -> "/home/user/.goose/packages/discord/_module.goose"
"discord/main.goose"            -> "$GOOSEROOT/packages/discord/main.goose/_module.goose", "$GOOSEROOT/packages/discord/main.goose"
"discord/commands"              -> "$GOOSEROOT/packages/discord/commands/_module.goose", "$GOOSEROOT/packages/discord/commands",k "$GOOSEROOT/packages/discord/commands.goose"
"discord/commands.goose"        -> "$GOOSEROOT/packages/discord/commands.goose/_module.goose", "$GOOSEROOT/packages/discord/commands.goose"
"discord/commands/main.goose"   -> "$GOOSEROOT/packages/discord/commands/main.goose/_module.goose", "$GOOSEROOT/packages/discord/commands/main.goose"
```

## Automatic module naming

1. Use the last path segment as the module name.
2. If the last path segment is `_module.goose`, use the second to last path segment as the module name.
3. If the path segment ends with `.goose`, remove the extension.
4. Remove any characters in the name that cannot be used in an identifier, leaving only `[a-z0-9_]`.
5. If the name: is a keyword, is empty, begins with a number, is `_module`, or is `_`: error; such modules must be named explicitly with `import <specifier> as <name>`.
6. If two modules have the same name, error; at least one of them must be named explicitly with `import <specifier> as <name>`.
