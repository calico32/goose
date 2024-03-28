// single-line comment

/*
  multi
  line
  comment
*/

// variable declaration
let mutable = 1
const constant = 2

// variable assignment
mutable = 3

// strings
"string"
"string with \"escaped\" quotes"
"string " + "concatenation"
"string $interpolation"
"string ${interpolation + " and concatenation"}"
"string with \n newlines"
"string with \$ \\ \x00 \o000 \u0000 \U00000000 escaped characters"

// numbers
// integers (64-bit)
1
1_000
0x1
0x1_000
0o1
0o1_000
0b1
0b1_000

// floats (64-bit)
1.0
1e1
1e+1
1e-1
1.0e1
1.0e+1
1.0e-1
1_000.0
1_000.000_1
1_000.000_1e1
1_000.000_1e1_000



// function declaration
fn add(a, b = 2)
  return a + b
end

// function call
add(1)
add(1, 4)

// if statement
if true
  print("true")
else if false
  print("false")
else
  print("else")
end

// while loop
repeat while x >= 0
  print("loop")
  x--
end

// forever loop
repeat forever
  print("loop")
end

// times loop
repeat 10 times
  print("loop")
end

// for loop
for x in iterable
  print("loop")
end

// language constants
true
false
null

// reserved names
_

// keywords, current and future
let
const
symbol
if then else
repeat while forever times for in
break continue
fn end return memo
import export as show
generator yield to step
struct init operator
try catch finally throw
async await

// operators

// binary arithmetic
+   +=      // add
-   -=      // subtract
*   *=      // multiply
/   /=      // divide (integer and float)
%   %=      // modulo
**  **=     // exponent

// unary arithmetic
+           // no-op
-           // negative

// binary logical
&&  &&=     // and
||  ||=     // or
??  ??=     // nullish coalescing

// unary logical
!           // not

// binary comparison
==          // equal
!=          // not equal
<           // less than
<=          // less than or equal
>           // greater than
>=          // greater than or equal

// other
=           // assign
?           // debug operator (postfix)



/* ---------------------------------------------

  future proposed syntax plans

---------------------------------------------- */

// ranges
0 to 10
0 to 10 step 2
for i in 0 to 10 step 5
  print(i)
end



// generators
generator fib(n)
  let a = 0
  let b = 1
  yield a
  yield b
  repeat while n > 0
    let c = a + b
    yield c
    a = b
    b = c
    n--
  end
end

for i in fib(100)
  print(i)
end

let g = fib(100)
print(g.next())
print(g.next())
print(g.next())
if !g.done()
  print("not done")
end



// structs
struct Point(x, y)
// with initialization code
struct Point(x, y) init
  // property access with #name or this.name
  #sum = #x + #y // same as `this.sum = this.x + this.y`
end

// receiver functions (and generators), must be defined in the same module as the struct
fn Point.foo()
  for prop in ["x", "y"]
    // property access with arbitrary expressions
    print(#[prop + ""]) // or print(this[prop])
  end

  print(this) // prints the entire struct
end

// other functions can use the bind operator ::
fn bar()
  print(#x)
end

// creation
let p = Point(1, 2)
p.foo()
p::bar()

// operator overloading
operator Point +(other)
  return Point(#x + other.x, #y + other.y)
end

operator Point :(low, high) // slice; point[low:high], point[low:], point[:high]
operator Point <=>(other) // spaceship operator; returns -1, 0, or 1 depending on the comparison

let p1 = Point(1, 2)
let p2 = Point(3, 4)
let p3 = p1 + p2





// single-line function declaration
fn add(a, b) -> a + b



// anonymous functions
let foo = fn(x, y) -> x + y // single-line
let bar = fn(a, b, c) // multi-line
  let d = a + b + c
  return "d² = ${d ** 2}"
end



// rest parameters
fn add(a, b, ...rest)
  let sum = a + b
  for x in rest
    sum += x
  end
  return sum
end

fn something(...params) -> somethingElse(...params)



// named parameters (functions, generators, structs)
fn add(a, b, :c, :d=4, ...e)
  sum = a + b + c + d
  for x in e
    sum += x
  end
  return sum
end
// parameters without : become positional or named
// parameters with : are always named
// rest parameters are positional (after everything else) or named
// named parameters can be used in any order
add(   1,    2, c: 3) // 10
add(   1, b: 2, c: 3) // 10
add(a: 1, b: 2, c: 3) // 10
add(c: 3, a: 1, b: 2, d: 4, /* e: */ 5, 6) // 21
add(b: 2, c: 3, a: 1,       e: [5, 6]) // 21

struct Vec3d(:x = 0, :y = 0, :z = 0)
let v = Vec3d(x: 1, y: 2) // Vec3d(1, 2, 0)



// bitwise operators
~ & | ^ << >>
&= |= ^= <<= >>=



// do expressions
let x = do
  let y = 1
  let z = 2
  y + z // end block with bare expression to return a value
end
print(x?) // x = 3
print(y?) // error: y is not defined



// if expressions
if true then 1 else 2

if true then 1
  else if false then 2
  else 3


// symbols
symbol @foo
symbol @bar

// use in composites, maps, etc.
let comp = { @foo: 1, @bar: 2 }
comp[@foo] // 1

// modules
export @bar
export symbol @baz

import "./foo.goose"
foo.@bar // the symbol @bar from module foo
comp[foo.@bar]

import "./foo.goose" show { @baz }
import "./foo.goose" show { @baz as @qux }
comp[@baz]
comp[@qux]



// type casting
float(1)
int(2.5)
string(3)
bool(expr) // truthy/falsy




// parsing
int.parse("123")          // 123
int.parse("0x123")        // 291
int.parse("1001001", 2)   // 73
float.parse("123.456")    // 123.456
float.parse("1e2")        // 100
bool.parse("true")        // true

int.parse("not a number") // throws an error
int.tryParse("not a number") // null
int.tryParse("1234not a number") // null
float.tryParse("123.456not a number") // null
bool.parse("not a boolean") // throws an error
bool.tryParse("not a boolean") // null



// constants
int.max // max int64
int.min // min int64
float.infinity // +inf
float.nan // NaN



// import/export
import "./foo.goose" // imports all symbols to `foo.{symbol}` (some name solving algorithm)
import "./foo.goose" as whatever // imports all symbols to `whatever.{symbol}`
import "./foo.goose" show { add, sub } // imports `add` and `sub` to global
import "./foo.goose" show { add as add2, sub as sub2 } // imports `add` and `sub` to `add2` and `sub2` on global
import "./foo.goose" show { add, ...foo } // imports `add` and everything else to `foo.{symbol}`
import "./foo.goose" show ... // imports everything to global

let x = 1
export { x } // exports to `x`

export fn add(a, b)
  return a + b
end // exports to `add`

export let y = 1 // exports to `y`
export const z = 2 // exports to `z`
export { y as y2, z as z2 } // exports to `y2` and `z2`

export "./foo.goose" // exports to `foo`
// equivalent to:
import "./foo.goose"
export { foo }

export "./foo.goose" as whatever // exports to `whatever`
export "./foo.goose" show { add, sub }  // exports to `add` and `sub`
export "./foo.goose" show ... // exports everything as top-level symbols

import "std:io" // imports the io module via the std protocol
import "pkg:discord" // imports the discord package via the pkg protocol
import "discord" // no protocol & not relative, defaults to the pkg protocol
import "discord/foo.goose" // foo.goose from the discord package
import "file:./foo.goose" // imports the foo.goose file via the file protocol
import "./foo.goose" // no protocol, relative: if the current module is a file, defaults to the file protocol, otherwise defaults to whatever network protocol is in use (http/https/ftp)
import "https://example.com/foo.goose" // downloads and imports foo.goose from the internet
import "mem:export const x = 1" as name // imports the string as a module (name required)
name.x // 1

// custom protocols
import.defineProtocol("foo", fn(url) {
  // ...
  // eventually return:
  // - a specifier with a built-in protocol
  // - a name to use for the module
  return {
    name: url,
    specifier: "mem:export const foo = \"foo\"",
  }
})

import "foo:bar" // uses the foo protocol (url == "bar")
bar.foo // "foo"



// exceptions
try
  // ...
catch e
  // ...
finally
  // ...
end

throw Error(message: "something went wrong")
throw Error(message: "something went wrong", cause: otherError)



// introspection (names subject to change)
import "std:reflect" show { Struct, Function, Symbol, Import, Export }

Struct.name(Point) // "Point"
Struct.properties(Point) // ["x", "y"]
Struct.receivers(Point) // ["foo"]
Struct.receivers(Point, "foo") // <function>
Struct.operators(Point) // ["+"]
Struct.operators(Point, "+") // <function>

Function.name(add) // "add"
Function.parameters(add) // ["a", "b"]
Function.restParameter(add) // null
Function.namedParameters(add) // ["c", "d"]
Function.parameterDefault(add, "d") // 4
Function.restParameterDefault(add) // null
Function.isMemoized(add) // false
Function.hasMemoized(add, 1, 2) // false

Symbol.name(@foo) // "foo"
Symbol.id(@foo) // 0x12345678 (some unique id)

// modules
Import.name("./foo.goose") // "foo" (the algorithm used to name modules)
Import.resolve("./foo.goose") // "/absolute/path/to/foo.goose"
Import.resolve("std:foo") // "<std/foo.goose>" or similar (stdlib has no absolute path)
Import.self.name // "bar" (the name of the current module)
Import.self.path // "/absolute/path/to/bar.goose"

export const x = 5
const y = 5
Export.name(x) // "x"
Export.name(y) // null
Export.name(5) // throws an error (literals can't be exported)
Export.name("foo") // throws an error
Export.list() // ["x"]





// fleshed out standard library
import "std:math"      // -> `math.{symbol}`: math constants, functions, etc
import "std:io"        // io operations
import "std:time"      // time and date util
import "std:http"      // http client
import "std:json"      // json parser and serializer
import "std:hash"      // hash functions
import "std:crypto"    // crypto functions (pbkdf2, secure random, etc)
import "std:random"    // pseudo random numbers and data
import "std:platform"  // language, runtime, and platform information (os, arch, etc)
import "std:canvas" show {  // 2d drawing canvas
  Canvas,
  rect,
  circle,
  ellipse,
  line,
  text,
  image,
}

math.PI // π
math.E // e
math.sin(1) // sin(1)

const f = io.open("foo.txt", "r") // file handle
f.read(1024) // read 1024 bytes
f.close()

const text = io.readSync("foo.txt", "utf-8") // read file synchronously
io.writeSync("bar.txt", text, "utf-8") // write file synchronously

// move to std:time?
time.sleep(1000) // sleep for 1000 milliseconds
time.milli() // current time in milliseconds
time.micro() // current time in microseconds
time.nano() // current time in nanoseconds

// move to std:os?
os.exit(1) // exit with code 1
os.args // command line arguments
os.env // environment variables

let res = await http.get("https://example.com")
res.status // 200
res.headers // { "Content-Type": "text/html", ... }
res.body // byte array
res.text() // unicode string (throws if not valid unicode)
res.tryText() // unicode string or null
res.json() // parsed json (throws if not valid json)
res.tryJson() // parsed json or null

let json = json.parse("[1, 2, 3]")
json[0] // 1
let str = json.stringify(json) // "[1,2,3]"

hash.sha256("foo") // "0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33"
hash.md5("foo") // "acbd18db4cc2f85cedef654fccc4a4d8"

crypto.randomBytes(16) // 16 secure random bytes
crypto.pbkdf2("password", "salt", 1000, 32) // 32 byte key

random.int(1, 10) // random int in range [1, 10)
random.float(1, 10) // random float in range [1, 10)
random.bool() // random bool
random.choice([1, 2, 3]) // random element from array
random.shuffle([1, 2, 3]) // random permutation of array
random.shuffle([1, 2, 3], 2) // random permutation of array with 2 elements
random.random() // random float in [0, 1)

platform.os // "windows", "linux", "darwin", "wasm", etc
platform.arch // "x86", "x64", "arm", "wasm", etc
platform.version // language version, semver string "1.0.0+whatever"
platform.runtime // implementation language (go or goose)


const c = Canvas(width: 100, height: 100) // create canvas
c.display() // platform-specific windowing (x11 on linux, html5 canvas on web, etc)

c.width // 100
c.height // 100

c.fill("#ff0000")
c.stroke("#00ff00")

c.draw(rect(0, 0, 100, 100))

c.draw(
  rect(x: 0, y: 0, w: 100, h: 100)
    .fill("#1e90ff"), // change color on the fly

  circle(x: 50, y: 50, r: 50)
    .stroke()
    .blur(radius: 5),

  ellipse(x: 50, y: 50, rx: 50, ry: 50),

  line(x1: 0, y1: 0, x2: 100, y2: 100)
    .shadow(dx: 5, dy: 5, blur: 5, color: "#ffffff"),

  text(x: 0, y: 0, text: "hello world", font: "sans-serif", size: 12),

  image(x: 0, y: 0, w: 100, h: 100, data: rgbaData),
)

c.clear() // clear canvas
c.save() // save canvas state
c.restore() // restore canvas state
c.translate(x, y) // translate canvas
c.rotate(angle) // rotate canvas
c.scale(x, y) // scale canvas
c.transform(a, b, c, d, e, f) // transform matrix
c.clip() // clip canvas
c.resetClip() // reset clip


// async/await

// async function
async fn foo()
  await bar()
  await baz()
end

let x = Promise(fn(resolve, reject)
  // ...
end)

foo() // returns a promise

// await expression (top level supported)
let x = await foo()

let x = async fn()
  let y = await foo()
  return y
end

// async generators
async generator foo()
  yield 1
  yield 2
  yield 3
end

for await x in foo()
  // ...
end

// async operators
async operator Point +(a, b)
  // ...
end

let r = p await+ q
let r = await (p + q)

// pipeline operator

printf("hello %s\n", random.choice((await (await http.get("https://people.example.com")).json()).people.map(fn(p) -> p.name)))

// becomes

await http.get("https://people.example.com")
  -> await _.json()
  -> _.people.map(fn(p) -> p.name)
  -> random.choice(_)
  -> printf("hello %s\n", _)



// pattern matching
let y = 5
let m = match x
  1 -> 10      // x is 1 -> 10
  2 -> 20      // x is 2 -> 20
  3 -> do      // x is 3 -> execute block
    return x * 10     // -> 30
  end
  [1] -> x[0]  // x is an array with 1 element, the number 1         -> 1
  [1, y] -> y  // x is an array of length 2, and is [1, 5]           -> 5
  [1, $y] -> y  // x is an array of length 2, the first element is 1 -> the second element
  [$x, $y, $z] -> x + y + y + z // x is an array of length 3, bind   -> x[0] + x[1] + x[1] + x[2]
  int($x) -> x // x is an integer                                    -> x
  float($x) -> x // x is a float                                     -> x
  { foo: $y } -> y // x has a property foo                           -> x.foo
  Array<int>($y) -> 5 // x is an array of integers
  arbitrary.Type($y) -> 5 // x is an instance of arbitrary.Type
  else -> 0 // match anything
end



// assertions
assert x == 5, "x is not 5"
assert x != 5, "x is 5"
assert true, "this is always true"
assert false, "this is always false" // throws AssertionError



// macros
// declarative macros receive expressions
macro foo!(x, y) -> x * y

foo!(1, 2)         // 1 * 2         -> 2
foo!(2 - 3, 5 + 4) // 2 - 3 * 5 + 4 -> -7 (no escaping!)

// procedural macros receive tokens
procmacro html!(tokens)
  for token in tokens
    print(token)
  end
end

let headingClass = "foo bar baz"
let subtitle = "hello world"

const document = html!(
  <html>
    <head>
      <title>hello world</title>
    </head>
    <body>
      <h1 class={headingClass}>hello world</h1>
      <h2>{subtitle}</h2>
    </body>
  </html>
)

fn sqlEscape(obj)
  // ...
end

procmacro sql!(tokens)
  let sql = ""
  for t in tokens
    match t
      when macro.Group($group) -> do
        sql += "\${sql_escape(${macro.eval(group)})}"
      end
      else -> sql += token
    end
    sql += " "
  end
  return "\"${sql.trim()}\""
end

sql!(select * from users where id = {id}) // "select * from users where id = ${sql_escape(id)}"



// type system
let x: int = 5
const y: float = 5.0
let z: string = "hello world"
let a: bool = true
let b: int[] = [1, 2, 3]
type Point = { x: int, y: int }
let c: Point = { x: 5, y: 5 }
let d: Point[] = [{ x: 5, y: 5 }, { x: 5, y: 5 }]
let e: any = 5

fn foo<T>(x: T): T
  return x
end

let e = foo(5) // e is an int
let f = foo<int>(5.0) // f is a float

macro foo!(tokens: macro.TokenTree[]): string
  return x
end

struct Node<T>(value: T, prev: Node<T>, next: Node<T>)

fn Node<T>.insertAfter(value: T): Node<T>
  let n = Node(value, this, this.next)
  this.next = n
  n.next.prev = n
  return n
end

let n = Node(5)

operator Node<T> +(a: Node<T>, b: Node<T>): Node<T>
  return a.insertAfter(b)
end
