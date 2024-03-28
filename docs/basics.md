# Language basics

Goose is a dynamically typed and interpreted language. It is designed to be simple and easy to learn, while still being powerful and expressive. It is inspired by languages such as JavaScript, Go, Rust, and Lua.

Goose source code is written in UTF-8 text files. The file extension is `.goose`.

- [Language basics](#language-basics)
  - [Comments](#comments)
  - [Variables](#variables)
    - [Reserved names](#reserved-names)
  - [Basic data types](#basic-data-types)
    - [Null](#null)
    - [Integers](#integers)
    - [Floats](#floats)
    - [Booleans](#booleans)
    - [Strings](#strings)
    - [Symbols](#symbols)
    - [Arrays](#arrays)
    - [Composites](#composites)
    - [Functions](#functions)
      - [Memoization](#memoization)

## Comments

Single line comments start with `//`.

```js
// This is a comment
```

Multi-line comments start with `/*` and end with `*/`. They cannot be nested.

```js
/* This is a
   multi-line comment */
```

## Variables

Variables are declared with the `let` or `const` keywords. Const variables cannot be reassigned (but are still mutable!), while let variables can. Identifiers, such as variable names, must match `[a-zA-Z_][a-zA-Z0-9_]*`.

```js
let x // = null; initialization optional
let y = 5
const z = 10 // initialization required for consts
```

Variables are untyped and can hold any value at any time.

### Reserved names

The following names cannot be used as variable names:

```js
_ // single underscore
```

## Basic data types

| Type        | String          | Description                                   |
| ----------- | --------------- | --------------------------------------------- |
| null        | `"null"`        | the absence of a value                        |
| int         | `"int"`         | 64-bit signed integer                         |
| float       | `"float"`       | IEEE 754 64-bit floating point number         |
| bool        | `"bool"`        | `true` or `false`                             |
| string      | `"string"`      | a string of characters                        |
| symbol      | `"symbol"`      | a unique module-specific identifier           |
| array       | `"array"`       | an array of any data type                     |
| composite   | `"composite"`   | generic "JS object" - a map of keys to values |
| function    | `"function"`    | a function                                    |
| int range   | `"int range"`   | a range of integers                           |
| float range | `"float range"` | a range of floats                             |

### Null

Null is the absence of a value. It is the default value of variables, function parameters, and function return values.

The null literal is `null`.

```js
let x // = null
let y = null

fn foo(a)
  // no return statement
end

foo() // = null
```

### Integers

Integers are arbitrary-precision signed integers. They can be written in decimal, hexadecimal, or octal notation, with or without separators.

```js
123
123_123_124
0xff
0xff_ab_ca
0o77123155464673
0b101010101_01010101
```

### Floats

Floats are IEEE 754 64-bit floating point numbers (aka double-precision). They can be written in decimal or scientific notation, with or without separators.

```js
1.23
1
0.23
1_00.2
12e2
1_200e+2
12e-2
1.2e5
```

### Booleans

Booleans are either `true` or `false`. They are the result of comparisons and logical operations. Their literals are `true` and `false`.

```js
true
false
```

### Strings

String literals are written using "double quotes". Strings may include any character (`"` must be escaped).

Strings must be valid UTF-8.

```js
"Hello world!"
"1234\t5678\b" // 1234    567
"\U0001F47D" // "👽"
```

The following escapes can be used in a string:

| Escape       | Character                   |
| ------------ | --------------------------- |
| `\"`         | U+0022 double quote         |
| `\$`         | U+0024 dollar sign          |
| `\\`         | U+005C backslash            |
| `\a`         | U+0007 alert                |
| `\b`         | U+0008 backspace            |
| `\n`         | U+000A line feed (lf)       |
| `\t`         | U+0009 tab                  |
| `\v`         | U+000B vertical tab         |
| `\f`         | U+000C form feed (ff)       |
| `\r`         | U+000D carriage return (cr) |
| `\xHH`       | U+HH (0 - 0xFF)             |
| `\uHHHH`     | U+HHHH (0 - 0xFFFF)         |
| `\UHHHHHHHH` | U+HHHHHHHH (0 - 0x10FFFF)   |

Expressions and variables can be interpolated into a string using `${expr}` or `$name`. Braces can be omitted if interpolating the value of a variable using its name unabiguously.

```js
let x = 2
let y = 4
let z = x + y
"The value of x is $x"
"x squared is ${x ** 2}"
"$x + $y = $z"
"the $yth power of $x is ${x ** y}" // $yth is invalid, use ${y}th instead
"\$19.99" // escaped $s
```

### Symbols

Symbols are unique, unreproduceable identifiers. They are used to represent unique values, such as keys in a composite. They are written using the `@` symbol followed by a valid identifier.

```js
symbol @foo
symbol @bar

let x = @foo
let y = {
  [@foo]: 123,
}
println(y[x]) // 123
```

### Arrays

Arrays are ordered lists of values. They can hold any data type, including other arrays. Literals are written using square brackets `[]`.

```js
[1, 2, 3]
[1, "two", 3.0]
[1, [2, 3], 4]
```

Elements can be accessed using square brackets `[]`. indices start at 0, and negative indices are supported.

```js
let arr = [1, 2, 3]
arr[3] // error
arr[2] // 3
arr[1] // 2
arr[0] // 1
arr[-1] // 3
arr[-2] // 2
arr[-3] // 1
arr[-4] // error
```

Slice expressions can be used to access a portion of an array. Negative indices are supported.

```js
let arr = [1, 2, 3, 4, 5]
arr[1:3]   // [2, 3]
arr[:3]    // [1, 2, 3]
arr[3:]    // [4, 5]
arr[1:-1]  // [2, 3, 4]
arr[-3:]   // [4, 5]
arr[-4:-2] // [3, 4]
arr[-2:-4] // error
arr[3:1]   // error
```

### Composites

Composites are "JS-like objects" that can be used to map keys to values. They are written using curly braces `{}`. Fields are unordered and do not maintain an insertion order.

Keys can be strings, numbers, or symbols. Keys must be unique in a literal. Fields are separated by commas `,`. A trailing comma is allowed.

Unbracketed keys must be valid identifiers. String keys that aren't valid identifiers can be surrounded in quotes, and any other value must be enclosed in brackets `[]`. Any expression can be used as a key given it evaluates to a valid key type.

```js
let key1 = "some key value"
symbol @key2

let obj = {
  foo: 123,
  bar: "baz",
  "qux": [1, 2, 3],
  [key1]: 456,
  [@key2]: "key2",
  [1 + 2]: "three",
}
```

Values can be accessed with dot notation `.` for keys that are valid identifiers or square brackets `[]` for arbitrary keys.

```js
obj.foo    // 123
obj["bar"] // "baz"
obj.qux[1] // 2
obj[key1]  // 456
obj[3]     // "three"
obj[@key2] // "key2"
```

Each composite has a key space for each type of key (string, number, symbol). Fields with keys of different types are distinct, even if they "look" the same. For example, `composite["1"]` and `composite[1]` refer to different fields.

```js
let obj = {
  "1": "string key space",
  [1]: "number key space",
}

obj["1"] // "string key space"
obj[1] // "number key space"
```

### Functions

Functions are first-class values. They can be assigned to variables, passed as arguments, and returned from functions. They can be declared using the `fn` keyword.

All function declarations are expressions. They can be given a name to assign the function to a constant variable, or they can be anonymous.

```js
// function declaration
fn foo(x, y)
  // code
end

// anonymous function
let bar = fn(x, y)
  // code
end
```

A function with a single expression return body can be written as a single line using an arrow `->`.

```js
fn add(x, y) -> x + y
const mul = fn(x, y) -> x * y
```

Parameters are untyped and are `null` by default, but can be given any default value. Default value expressions are evaluated once at declaration time and then copied for each call.

```js
fn multiPrint(x, y)
  println(x)
  println(y)
end

multiPrint("123")
/* prints:
123
<null>
*/
```

```js
let value = "y"
fn multiPrint(x, y = value)
  println(x)
  println(y)
end

value = "z"

multiPrint("123") // prints "123" and "y" (not "z")
```

Using higher-order functions, functions can be passed as arguments and returned from other functions.

```js
fn caller(f, arg) -> f(arg)
fn square(x) -> x ** 2
caller(square, 2) // 4
```

```js
fn makeAdder(x) -> fn(y) -> x + y
plus2 = makeAdder(2)
plus2(3) // 5
```

#### Memoization

Functions can be memoized using the `memo` keyword. Memoized functions will only be called once for each unique set of arguments.

```js
let fibCalls = 0
fn fib(n)
  fibCalls++

  if n == 0 || n == 1
    return n
  end

  return fib(n - 1) + fib(n - 2)
end

let fibMemoCalls = 0
memo fn fibMemo(n)
  fibMemoCalls++

  if n == 0 || n == 1
    return n
  end

  return fibMemo(n - 1) + fibMemo(n - 2)
end

println(fib(10)) // 55
println(fibCalls) // 177

println(fibMemo(10))  // 55
println(fibMemoCalls) // 11
```

### Integer and float ranges

Ranges are used to represent a sequence of integers or floats. They are written using the `to` keyword. If either the start or stop value is a float, the range will be a float range.

The range is inclusive of the start value and exclusive of the stop value.

```js
1 to 10
1.0 to 10.0
```

An optional `step ...` clause can be used to specify the step between each value in the range.

```js
1 to 10 step 2 // 1, 3, 5, 7, 9
10.2 to 10.8 step 0.1 // 10.2, 10.3, 10.4, 10.5, 10.6, 10.7
```

Ranges are their own primitive type; they are not arrays or composite types.

## Operators

Operators are used to perform operations on values. They can be unary, binary, or ternary.

### Arithmetic operators

| Operator | Description    | Types     |
| -------- | -------------- | --------- |
| `+`      | Addition       | See below |
| `-`      | Subtraction    | See below |
| `*`      | Multiplication | See below |
| `/`      | Division       | See below |
| `%`      | Remainder      | int       |
| `**`     | Exponentiation | See below |

#### Arithmetic Types

The following table shows the result type of arithmetic operations based on the types of the operands.

| Left  | Right | Result |
| ----- | ----- | ------ |
| int   | int   | int    |
| int   | float | float  |
| float | int   | float  |
| float | float | float  |

### Increment and decrement

The increment and decrement operators are unary operators that add or subtract 1 from a variable.

| Operator | Description | Types      |
| -------- | ----------- | ---------- |
| `++`     | Increment   | int, float |
| `--`     | Decrement   | int, float |

### Comparison operators

| Operator | Description   | Types      |
| -------- | ------------- | ---------- |
| `==`     | Equal         | Any        |
| `!=`     | Not equal     | Any        |
| `<`      | Less than     | int, float |
| `<=`     | Less equal    | int, float |
| `>`      | Greater than  | int, float |
| `>=`     | Greater equal | int, float |

### Logical operators

| Operator                   | Description | Types              |
| -------------------------- | ----------- | ------------------ |
| `!`                        | Not         | Any (unary) → bool |
| `&&`                       | And         | Any → bool         |
| `\|\|`                     | Or          | Any → bool         |
| `??`                       | Coalesce    | Any → Any          |
| `if` ... `then` ... `else` | Ternary     | Any → Any          |

### Bitwise operators

| Operator | Description | Types |
| -------- | ----------- | ----- |
| `&`      | And         | int   |
| `\|`     | Or          | int   |
| `^`      | Xor         | int   |
| `~`      | Not         | int   |
| `<<`     | Shift left  | int   |
| `>>`     | Shift right | int   |

### Assignment operators

| Operator | Description | Types                                          |
| -------- | ----------- | ---------------------------------------------- |
| `=`      | Assign      | Any                                            |
| `+=`     | Add assign  | int, float (follows arithmetic operator rules) |
| `-=`     | Sub assign  | int, float                                     |
| `*=`     | Mul assign  | int, float                                     |
| `/=`     | Div assign  | int, float                                     |
| `%=`     | Mod assign  | int                                            |
| `**=`    | Exp assign  | int, float                                     |
| `&=`     | And assign  | int                                            |
| `\|=`    | Or assign   | int                                            |
| `^=`     | Xor assign  | int                                            |
| `<<=`    | Shl assign  | int                                            |
| `>>=`    | Shr assign  | int                                            |
| `&&=`    | And assign  | bool                                           |
| `\|\|=`  | Or assign   | bool                                           |
| `??=`    | Coalesce    | Any                                            |

### Other operators

| Operator | Description | Types |
| -------- | ----------- | ----- |
| `.`      | Access      | Any   |
| `[]`     | Index       | Any   |
| `()`     | Call        | Any   |

## Control flow

Goose has several control flow statements to alter the flow of a program.

### If

The `if` statement is used to execute a block of code if a condition is true.

```js
if condition
  // code
end
```

`else` and `else if` can be used to execute code if the condition is false.

```js
if condition
  // code
else
  // code
end
```

```js
if condition1
  // code
else if condition2
  // code
else
  // code
end
```

### If expressions

`if` can also be used as an expression like the ternary operator in other languages.

```js
let x = if temperature > 100 then "hot" else "cold"
```

### Loops

Goose has several types of loops.

#### `repeat x times`

The `repeat` statement is used to execute a block of code a fixed number of times.

```js
repeat 5 times
  println("Hello, world!")
end
```

#### `repeat while`

The `repeat while` statement is used to execute a block of code while a condition is true.

```js

let x = 0
repeat while x < 5
  println(x)
  x++
end
```

#### `repeat forever`

The `repeat forever` statement is used to execute a block of code indefinitely, until a `break` statement is encountered.

```js
repeat forever
  println("Hello, world!")
  if condition
    break
  end
end
```

#### `for ... in ...`

The `for ... in ...` statement is used to iterate over the elements of a sequence (array, composite, or range).

```js
for i in [1, 2, 3]
  println(i)
end
```

```js
for i in 1 to 11
  println(i) // 1, 2, 3, 4, 5, 6, 7, 8, 9, 10
end
```

```js
for i in 1 to 11 step 2
  println(i) // 1, 3, 5, 7, 9
end
```

#### `break` and `continue`

The `break` statement is used to exit a loop early.

```js
for i in 0 to 10
  if i == 5
    break
  end
  println(i) // 0, 1, 2, 3, 4
end
```

The `continue` statement is used to skip the rest of the current iteration and continue to the next.

```js
for i in 0 to 10
  if i % 2 == 0
    continue
  end
  println(i) // 1, 3, 5, 7, 9
end
```

### Match

The `match` statement is used to execute a block of code based on the value of an expression. It supports pattern matching and guards.

```js
match response
  200 -> println("OK")
  404 -> println("Not found")
  $x if x >= 500 && x < 600 -> println("Server error: $x")
  else -> println("Unknown status")
end
```

### Return

The `return` statement is used to exit a function early and return a value.

```js
fn foo(x)
  if x < 0
    return "negative"
  end
  return "positive"
end
```

## Module system

Goose has a simple module system. Each file is a module, and modules can import other modules.

```js
// file: foo.goose
export const x = 123

export fn greet()
  println("Hello, world!")
end
```

```js
// file: main.goose
import "./foo.goose"

println(foo.x) // 123
foo.greet() // "Hello, world!"
```

### `export`

The `export` keyword can be used to export variables and functions from a module.

```js
export const x = 123

export fn greet()
  println("Hello, world!")
end

export struct Point(x, y)

// functions with receivers don't need to be exported, they're always available
fn Point.distanceTo(other)
  return math.sqrt((other.x - x) ** 2 + (other.y - y) ** 2)
end
```

`export` can also be used to export a list of variables at the same time. Fields can also be renamed when exporting.

```js
// file: foo.goose
const x = 123
fn y() -> 456
export { x, y as z }
```

```js
// file: main.goose
import "./foo.goose"

println(foo.x) // 123
println(foo.z()) // 456
```

You can also use `export` to export other modules directly.

```js
// file: foo.goose
export "./bar.goose"

export "./baz.goose" show { x, y }

export "./qux.goose" show ...
```

### `import`

The `import` statement is used to import variables and functions from other modules.

A plain import will import all exported variables and functions from the module. They can be accessed using the module name, which is derived from the file name.

```js
import "./foo.goose"

foo.abc()
println(foo.xyz)
```

You can specify an explicit module name using the `as` clause.

```js
import "./foo.goose" as bar

bar.abc()
println(bar.xyz)
```

Import specific variables to the current module scope using the `show` clause.

```js
import "./foo.goose" show { abc, xyz }

abc()
println(xyz)
```

Import everything from a module using ellipsis `...`.

```js
import "./foo.goose" show ...

abc()
def()
println(xyz)
```

Import statements can be nested to import from submodules. This is useful for importing packages with a directory structure.

(The same can be done with `export "..." show { ... }`)

```js
import "pkg:fiber" show {
  App,
  Context,
  "client" show { Client }, // import "pkg:fiber/client" show { Client }
  "middleware" show {
    "cors" show { Cors }, // import "pkg:fiber/middleware/cors" show { Cors }
    "logger" show { Logger }, // import "pkg:fiber/middleware/logger" show { Logger }
  },
}
```

### Semantics

Modules are singletons. They are loaded and executed once, and their exports are cached for future imports.

Modules are executed in a separate scope from the importing module. This means that variables and functions in the imported module are not accessible from the importing module unless they are exported.

Modules are executed in the order they are imported. Circular imports are not allowed.

### Module schemes

Goose supports several module schemes for importing modules from different sources. In most cases, paths without schemes are assumed to be the same scheme as the importing module.

#### File system

The `file:` scheme is used to import modules from the file system. Relative paths are resolved relative to the importing module (as long as the importing module is also a file; otherwise, it is an error).

This is the scheme used in most Goose code. It is usually omitted, as most imports are in files on-disk and thus their imports are implicitly `file:` imports.

```js
import "file:./foo.goose"
import "./foo.goose" // same as above (in regular files)
import "file:/usr/local/lib/foo.goose" // absolute path
import "/usr/local/lib/foo.goose" // same as above
```

#### Package

The `pkg:` scheme is used to import modules from an installed package. The path is resolved relative to the package root.

```js
import "pkg:fiber"
import "pkg:discord/ws"
```

#### Standard library

The `std:` scheme is used to import modules from the standard library.

```js
import "std:json"
import "std:math"
import "std:http/server"
```

#### Network

The `https:` and `http:` schemes are used to import modules from the network.

```js
import "https://example.com/foo.goose"
import "http://example.com/bar/baz.goose"
```

### Standard library

Goose has a small standard library that is always available. It includes modules for common tasks such as I/O, math, and networking. See the [API reference](https://goose.calico.lol/docs/api-reference) for more information.

```js
import "std:math"
import "std:readline"

let input = readline.read().split(" ").map(int.parse)
let sum = math.sum(input)
println(sum)
```
