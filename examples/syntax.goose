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
1
1.0
1e1
1e+1
1e-1
1.0e1
1.0e+1
1.0e-1
1_000
1_000.0
1_000.000_1
1_000.000_1e1
1_000.000_1e1_000

0x1
0x1_000
0o1
0o1_000
0b1
0b1_000

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
let const
if then else
repeat while forever times for in
break continue
fn end return memo
import export as from
generator yield to step
struct init operator
try catch finally throw

// operators
// binary arithmetic
+ - * / % **
+= -= *= /= %= **=
// unary arithmetic
+ -
// binary logical
&& ||
&&= ||=
// unary logical
!
// binary comparison
== != < <= > >=
// assignment
=

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
  // property access with #prefix
  #sum = #x + #y
end

// receiver functions (and generators)
fn Point.foo()
  for prop in ["x", "y"]
    // property access with arbitrary expressions
    print(#[prop + 1])
  end
end

// creation
let p = Point(1, 2)
p.foo()

// operator overloading
operator Point +(other)
  return Point(#x + other.x, #y + other.y)
end

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
fn add(a, b, :c, :d=4)
  return a + b + c + d
end
// parameters without : become positional or named
// parameters with : are always named
// named and rest parameters cannot be mixed?
// named parameters can be used in any order
add(   1,    2, c: 3) // 10
add(   1, b: 2, c: 3) // 10
add(a: 1, b: 2, c: 3) // 10
add(b: 2, c: 3, a: 1) // 10

struct Vec3d(:x = 0, :y = 0, :z = 0)
let v = Vec3d(x: 1, y: 2) // Vec3d(1, 2, 0)



// bitwise operators
~ & | ^ << >>
&= |= ^= <<= >>=


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



// type casting
float(1)
int(2.5)
string(3) // bad idea?
bool(expr) // truthy/falsy



// import/export
import "./foo.goose" // imports all symbols to `foo.{symbol}` (some name solving algorithm)
import "./foo.goose" as whatever // imports all symbols to `whatever.{symbol}`
import { add, sub } from "./foo.goose" // imports `add` and `sub` to global
import { add as add2, sub as sub2 } from "./foo.goose" // imports `add` and `sub` to `add2` and `sub2` on global

let x = 1
export x // exports to `x`

export fn add(a, b)
  return a + b
end // exports to `add`

export let y = 1 // exports to `y`
export const z = 2 // exports to `z`
export y as y2, z as z2 // exports to `y2` and `z2`



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



// introspction (names subject to change)
struct.name(Point) // "Point"
struct.properties(Point) // ["x", "y"]
struct.receivers(Point) // ["foo"]
struct.receivers(Point, "foo") // <function>
struct.operators(Point) // ["+"]
struct.operators(Point, "+") // <function>

fn.name(add) // "add"
fn.parameters(add) // ["a", "b"]
fn.restParameter(add) // null
fn.namedParameters(add) // ["c", "d"]
fn.parameterDefault(add, "d") // 4
fn.restParameterDefault(add) // null

symbol.name(@foo) // "foo"
symbol.id(@foo) // 0x12345678 (some unique id)

int.max // max int64
int.min // min int64
float.infinity // +inf
float.nan // NaN



// fleshed out standard library
import "std:math" // -> `math.{symbol}`, etc
import "std:io"
import "std:time"
import { get } from "std:http"

math.PI // π
math.E // e
math.sin(1) // sin(1)

const f = io.open("foo.txt", "r")
f.read(1024) // read 1024 bytes
f.close()

// move to std:time?
time.sleep(1000) // sleep for 1000 milliseconds
time.milli() // current time in milliseconds
time.micro() // current time in microseconds
time.nano() // current time in nanoseconds

// move to std:os?
os.exit(1) // exit with code 1
os.args // command line arguments
os.env // environment variables

let res = get("https://example.com")
res.status // 200
res.headers // { "Content-Type": "text/html", ... }
res.body // <buffer?>
