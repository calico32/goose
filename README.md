# goose

**Go**ose is a minimal programming language built on top of [**Go**](https://go.dev/). It is a language that is designed to be easy to learn and easy to use.

# Table of Contents
- [goose](#goose)
- [Table of Contents](#table-of-contents)
- [Usage](#usage)
- [Roadmap](#roadmap)
  - [Less urgent](#less-urgent)
- [Language description](#language-description)
  - [Comments](#comments)
  - [Variables and data types](#variables-and-data-types)
  - [Strings](#strings)
  - [Arrays](#arrays)
  - [Composites](#composites)
  - [Literals](#literals)
  - [Control flow](#control-flow)
    - [If statements](#if-statements)
    - [For loops](#for-loops)
    - ["Repeat x times" loops](#repeat-x-times-loops)
    - [Repeat while loops](#repeat-while-loops)
    - [Repeat forever loops](#repeat-forever-loops)
    - [Branching](#branching)
  - [Functions](#functions)
    - [Memoization](#memoization)
  - [Standard library](#standard-library)

# Usage

Use `goose/ast`, `goose/scanner`, `goose/parser`, and `goose/interpreter` in your code:

```sh
go get github.com/wiisportsresort/goose
```

Run your code with `cli/goose.go`:

```sh
go run cli/goose.go [args]
go run cli/goose.go --help
```

Install the `goose` binary in your `$PATH`:

```sh
go install cli/goose.go
```

# Roadmap

- [ ] speed up loops with many iterations
- [ ] add composite types
- [ ] add anonymous functions
- [ ] explicit type annotations
- [ ] type casting

## Less urgent

- [ ] module system, imports/exports
- [ ] implement compiler

# Language description

## Comments

Single line `//` comments and multi-line `/* */` comments are supported. Multi-line comments cannot be nested.

```js
// this is comment

/*

      this
              
                is
            a
        
      comment

*/

/* 
  foo
  
  /*    bar     */
  
  invalid
*/
```

## Variables and data types

Variables are declared with `let` or `const`, in the following manner: 
```js
let y = "hello world"
const SIZE = 1024
```

Attempting to modify a constant variable (including modifying arrays and composites) is an error.

Variables are not restricted to one type and can be changed to any type later on.

The following data types are available:

| Type      | Description                                   |
| --------- | --------------------------------------------- |
| int       | 64-bit signed integer                         |
| float     | IEEE 754 64-bit floating point number         |
| string    | a string                                      |
| array     | an array of any data type                     |
| composite | generic "JS object" - a map of keys to values |

The following names are reserved and cannot be declared:
```
_
```

## Strings

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
| `\b`         | U+0008 backslash            |
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

## Arrays

Arrays are written using square brackets `[]`. 

Array literals can be empty, can contain any data type, and can be nested.
```js
[]
["abc", 123, false]
[1, 2, 3, [4, 5, 6]]
```

Array initializer expressions can be used to initialize an array. The special identifier `_` can be used to access the current index during initialization.
```js
[null; 10]         // array of 10 nulls
[1; 10]            // array of 10 1s
[_; 10]            // [0, 1, 2, 3, 4, 5, 6, 7, 8, 9]
[(_ + 1) ** 2; 10] // [1, 4, 9, 16, 25, 36, 49, 64, 81, 100]
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

## Composites

Coming soon!

## Literals

Integer literals can be written in decimal, hexadecimal, octal, or binary. Separators are allowed.
```js
123
123_123_124
0xff
0xff_ab_ca
0o77123155464673
0b101010101_01010101
```

Examples of floating point literals:
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

Boolean literals are written as `true` or `false`.
```js
true
false
```

Null literals are written as `null`.
```js
null
```

## Control flow

### If statements

If statements can be used to execute a block of code if a condition is true. If the condition is not a boolean, it is coerced to a boolean.
```js
if expr
  // block of code
else if expr
  // block of code
else
  // block of code
end
```

The following expressions are falsy:
```js
""
0
false
null
[] // empty array
{} // empty composite
```

All other expressions are truthy.

### For loops

For loops can be used to iterate over an array or composite.
```js
for i in arr
  // block of code
end
```

### "Repeat x times" loops

Repeat loops can be used to execute a block of code a number of times.
```js
repeat 5 times
  // block of code
end
```

### Repeat while loops

Repeat while loops can be used to execute a block of code while a condition is true.
```js
repeat while n < 12
  // block of code
end
```

### Repeat forever loops

Repeat forever loops can be used to execute a block of code forever (or until a return/branch statement is reached).
```js
repeat forever
  // block of code
end
```

### Branching

The `break` and `continue` statements can be used to end the current iteration or stop looping.
```js
for i in arr
  if i % 2 == 0
    continue
  end
end

repeat forever
  if condition
    break
  end
end
```

## Functions

Functions are declared using the `fn` keyword. Parameters can be specified using `(` and `)` and can be separated by commas.
```js
fn foo(x, y)
  // block of code
end

foo(1, 2)
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
<nil>
*/
```

```js
let value = "y"
fn foo(x = "x", y = value, z = [1])
  z += 2
  print(x, y, z)
  return z
end

foo()                  // prints: x y [1, 2]
foo(1)                 // prints: 1 y [1, 2]
foo(1, 2)              // prints: 1 2 [1, 2]
let z = foo(1, 2, [2]) // prints: 1 2 [2, 2]

value = "z"
z += 4
foo() // prints: x y [1, 2]

```

Declared functions are constant. They cannot be reassigned.

Functions are first-class values. They can be passed as arguments to other functions, and can be returned from other functions.

Anonymous functions will be supported in the future.
```js
fn caller(f, arg)
  f(arg)
end

fn square(x)
  return x ** 2
end

caller(square, 2) // 4
```

```js
fn makeAdder(x)
  fn adder(y)
    return x + y
  end
  
  return adder
end

plus2 = makeAdder(2)
plus2(3) // 5
```

### Memoization

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

## Standard library

The following functions are available in the global scope:

| Signature                    | Description                                                                      |
| ---------------------------- | -------------------------------------------------------------------------------- |
| `print(...any)`              | Prints the arguments to the console, separated by spaces, followed by a newline. |
| `printf(fmt, ...any)`        | Prints the arguments to the console, formatted according Golang's `fmt.Printf`.  |
| `join(arr, sep=",")`         | Joins the elements of an array into a string.                                    |
| `round(num)`                 | Rounds a number to the nearest integer.                                          |
| `floor(num)`                 | Rounds a number down to the nearest integer.                                     |
| `ceil(num)`                  | Rounds a number up to the nearest integer.                                       |
| `exit(code=0)`               | Exits the program with the given exit code.                                      |
| `nano()`                     | Returns the current time in nanoseconds.                                         |
| `milli()`                    | Returns the current time in milliseconds.                                        |
| `sleep(ms)`                  | Sleeps for the given number of milliseconds.                                     |
| `len(arr)`                   | Returns the length of an array.                                                  |
| `indices(arr)`               | Returns an array of integers from 0 to the length of the array.                  |
| `padLeft(str, len, ch=" ")`  | Pads the beginning of a string with the given character.                         |
| `padRight(str, len, ch=" ")` | Pads the end of a string with the given character.                               |

The standard library functions are not constants and can be overwritten.
