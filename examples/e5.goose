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
let strings = fruits.map(fn(fruit, i) -> "${i + 1}. $fruit")

print(strings.join("\n"))
