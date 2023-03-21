// random-ish language features


// anonymous functions
const add1 = fn(x, y)
  return x + y
end

// single line functions
const add2 = fn(x, y) -> x + y


// if expressions
fn diff(x, y) -> if x > y then x - y else y - x


// do blocks
let x = 1
do
  let x = 2
  let y = 3
  print(x?) // x = 2
  print(y?) // y = 3
end
print(x?) // x = 1
print(y?) // error: y is not defined
