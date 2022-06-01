let powers = []

const N = 10
let i = 0

repeat while i < N
  // append to array using +=
  powers += 2 ** i
  i++
end

// print powers
for n in powers
  print(n)
end
