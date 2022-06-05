// returns the n-th fibonacci number
// memoized for speed
let calls = 0
memo fn fib(n)
  calls++

  if n == 0 || n == 1
    return n
  else
    return fib(n - 1) + fib(n - 2)
  end
end

const max = 50
let n = 1

repeat while n < max
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
  print("the $n$suffix fibonacci number is ${fib(n)}")
  n++
end

print("calls: $calls")
