// more array initializers

const fruits = [
  "apple",
  "banana",
  "orange",
  "pear",
  "mango",
  "grape",
  "pineapple",
  "strawberry",
]

// ["1. apple", "2. banana", ..., "8. strawberry"]
let strings = ["${_ + 1}. ${fruits[_]}"; len(fruits)]

print(join(strings, "\n"))
