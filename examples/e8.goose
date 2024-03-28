// pipeline operator
// value on the left is passed as _ to the expression on the right
// evaluated left to right

fn add(x, y) -> x + y

let x = 10

add(x, 5)?
  -> add(_?, 10)?
  -> add(_?, 15)?
  -> add(_?, 20)?
  -> _?


async fn example() ->
  await http.get("https://people.example.com")
    -> await _.json()
    -> _.people.map(fn(p) -> p.name)
    -> random.choice(_)
    -> printf("hello %s", _)
