let people = []

people.push({
  name: "John",
  age: 30,
})

people.push({
  ["name"]: "Jane",
  ["age"]: 25,
})

print(people[0].name) // John
print(people[1].name) // Jane

print(people[0]["age"]) // 30
print(people[1]["age"]) // 25


for person in people
  print()
  for key in keys(person)
    print(key + ": " + person[key])
  end
end
