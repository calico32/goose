// returns the n-th fibonacci number
// memoized for speed
let calls = 0
memo fn fib(n)
  calls++

  if n <= 1
    return n
  else
    return fib(n - 1) + fib(n - 2)
  end
end


for n in 1 to 101
  let suffix = "th"
  if n % 100 <= 10 || n % 100 >= 20
    if n % 10 == 1
      suffix = "st"
    else if n % 10 == 2
      suffix = "nd"
    else if n % 10 == 3
      suffix = "rd"
    end
  end
  println("the ${n}$suffix fibonacci number is ${fib(n)}") // $ is valid in identifiers, so the n needs to be wrapped in {}
end

println("calls: $calls")
