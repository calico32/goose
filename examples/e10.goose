import "std:collections" show { Stack }
import "std:readline"

print("Enter a stack expression: ")
let input = readline.read()

let stack = Stack()

for token in input.split(" ")
  if token == "+" || token == "-" || token == "*" || token == "/"
    let b = stack.pop()
    let a = stack.pop()
    if token == "+"
      stack.push(a + b)
    else if token == "-"
      stack.push(a - b)
    else if token == "*"
      stack.push(a * b)
    else if token == "/"
      stack.push(a / b)
    end
  else
    stack.push(float.parse(token))
  end
end

println("Result: ${stack.pop()}")
