/// Read a line from standard input. The trailing newline is removed.
export native fn read()
export fn prompt(str)
  print(str)
  return read()
end
