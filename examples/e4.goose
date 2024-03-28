// array initializer expression
//             ↓ value ↓ count
const nulls = [null;   10]
// [null, null, null, null, null, null, null, null, null, null]
print(nulls)

// eval an expression, use _ as the index of the current element
fn squares(n) -> [(_ + 1) ** 2; n]
// [1 ** 2, 2 ** 2, 3 ** 2, ..., n ** 2]

print("the first 10 squares are ${squares(10)}")
