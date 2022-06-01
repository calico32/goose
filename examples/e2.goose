fn mul(a, b)
  let result = 0

  repeat b times
    result += a
  end

  return result
end

const multiplied = mul(21, 19)

print("the result is $multiplied") // parens optional if single name
