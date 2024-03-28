struct Node(value, next)

export struct Stack(head)

fn Stack.push(value)
  if !this.head
    this.head = Node(value, null)
  else
    this.head = Node(value, this.head)
  end
end

fn Stack.pop()
  if !this.head
    return null
  end
  let value = this.head.value
  this.head = this.head.next
  return value
end

fn Stack.peek()
  if !this.head
    return null
  end
  return this.head.value
end

fn Stack.isEmpty()
  return !this.head
end

fn Stack.toString()
  let current = this.head
  let s = "["
  repeat while current
    s += current.value
    s += ", "
    current = current.next
  end
  s = s.slice(0, -2)
  s += "]"
  return s
end

fn Stack.toArray()
  let current = this.head
  let a = []
  repeat while current
    a.push(current.value)
    current = current.next
  end
  return a
end

fn Stack.print()
  print(this)
end

fn Stack.length()
  let count = 0
  let current = this.head
  repeat while current
    count++
    current = current.next
  end
  return count
end
