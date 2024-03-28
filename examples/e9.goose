import "std:readline"

println("Enter a number to convert.")

repeat forever
  print("> ")

  let line = readline.read()
  if line == "exit"
    break
  end

  try
    if len(input) != 0 && input[0].isDigit()
      const roman = toRoman(int.parse(input))
      println("Roman: $roman")
    else
      const arabic = toArabic(input)
      println("Arabic: $arabic")
    end
  catch e
    println("Error: ${e.message}")
  end
end
