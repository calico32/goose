# Language basics

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

Variables are declared with the `let` or `const` keywords. Const variables cannot be reassigned (but are still mutable!), while let variables can.

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

| Type      | String        | Description                                   |
| --------- | ------------- | --------------------------------------------- |
| null      | `"null"`      | the absence of a value                        |
| int       | `"int"`       | 64-bit signed integer                         |
| float     | `"float"`     | IEEE 754 64-bit floating point number         |
| bool      | `"bool"`      | `true` or `false`                             |
| string    | `"string"`    | a string of characters                        |
| symbol    | `"symbol"`    | a unique module-specific identifier           |
| array     | `"array"`     | an array of any data type                     |
| composite | `"composite"` | generic "JS object" - a map of keys to values |
| function  | `"function"`  | a function                                    |

### Null

Null is the absence of a value. It is the default value of variables, function parameters, and function return values.

The null literal is `null`.

```js
let x // = null
let y = null

fn foo(a)
end

foo() // = null
```

### Integers

Integers are 64-bit signed integers. They can be written in decimal, hexadecimal, or octal notation, with or without separators.

```js
123
123_123_124
0xff
0xff_ab_ca
0o77123155464673
0b101010101_01010101
```

### Floats

Floats are IEEE 754 64-bit floating point numbers. They can be written in decimal or scientific notation, with or without separators.

```js
1.23
1.
.23
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
```js
"Hello world!"
"1234\t5678\b" // 1234    567
"\U0001F47D"   // "👽"
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

Symbols are unique identifiers that are local to a module. They are declared with the `symbol` statement, and are prefixed with `@`.

```js
symbol @foo
symbol @bar

let x = @foo
print(@foo)
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
arr[3]  // error
arr[2]  // 3
arr[1]  // 2
arr[0]  // 1
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
arr[-3:]   // [4, 5]
arr[-4:-2] // [3, 4]
arr[-2:-4] // error
arr[3:1]   // error
```

### Composites

Composites are "JS-like objects" that can be used to map keys to values. They are written using curly braces `{}`. Fields are unordered and do not maintain an insertion order.

Keys can be strings, numbers, or symbols. Keys must be unique in a literal. Fields are separated by commas `,`. A trailing comma is allowed. 

Unbracketed keys must be valid identifiers, and any other value must be enclosed in brackets `[]`. Any expression can be used as a key given it evaluates to a valid key type.

```js
let key1 = "some key value"
symbol @key2

let obj = {
  foo: 123,
  bar: "baz",
  qux: [1, 2, 3],
  [key1]: 456,
  [@key2]: "key2",
  [1 + 2]: "three",
}
```

Values can be accessed with dot notation `.` or square brackets `[]`. 
```js
obj.foo    // 123
obj["bar"] // "baz"
obj.qux[1] // 2
obj[key1]  // 456
obj[3]     // "three"
obj[@key2] // "key2"
```

Each composite has 3 key spaces—for each type of key, i.e. `composite["1"]` and `composite[1]` refer to different fields.

```js
let obj = {
  "1": "string key space",
  [1]: "number key space",
}

obj["1"] // "string key space"
obj[1]   // "number key space"
```

### Functions

Functions are first-class values. They can be assigned to variables, passed as arguments, and returned from functions. They can be declared using the `fn` keyword.

All function declarations are expressions. They can be given a name to assign the function to a constant variable, or they can be anonymous.

```js
fn foo(x, y)
  // code
end

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
  print(x)
  print(y)
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
  print(x)
  print(y)
end

multiPrint("123") // prints "123" and "y"
```

```js
fn caller(f, arg) -> f(arg)
fn square(x) -> x ** 2
caller(square, 2) // 4

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

print(fib(10)) // 55
print(fibCalls) // 177

print(fibMemo(10))  // 55
print(fibMemoCalls) // 11
```
