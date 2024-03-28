struct Node(value, next = null)

export struct LinkedList(head = null)

fn LinkedList.tail()
  let current = this.head
  if current == null
    return null
  end
  repeat while current.next
    current = current.next
  end
  return current
end

fn LinkedList.add(value)
  let tail = this.tail()
  if tail == null
    this.head = Node(value)
    return
  end
  tail.next = Node(value)
end

fn LinkedList.length()
  let current = this.head
  let i = 0
  repeat while current
    i++
    current = current.next
  end
  return i
end

fn LinkedList.get(index)
  let current = this.head
  let i = 0
  repeat while i < index
    i++
    current = current.next
  end
  return current.value
end

fn LinkedList.insert(index, value)
  let current = this.head
  let i = 0
  repeat while i < index - 1
    i++
    current = current.next
  end
  let next = current.next
  current.next = Node(value, next)
end

fn LinkedList.set(index, value)
  let current = this.head
  let i = 0
  repeat while i < index
    i++
    current = current.next
  end
  current.value = value
end

fn LinkedList.remove(index)
  let current = this.head
  let i = 0
  repeat while i < index - 1
    i++
    current = current.next
  end
  current.next = current.next.next
end

fn LinkedList.toString()
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

fn LinkedList.toArray()
  let current = this.head
  let a = []
  repeat while current
    a.push(current.value)
    current = current.next
  end
  return a
end

fn LinkedList.print()
  print(this)
end

fn LinkedList.drop(n)
  let current = this.head
  repeat n times
    current = current.next
  end
end

fn LinkedList.map(func)
  let clone = LinkedList()
  let current = this.head
  repeat while current
    clone.add(func(current.value))
    current = current.next
  end
end

fn LinkedList.filter(func)
  let clone = LinkedList()
  let current = this.head
  repeat while current
    if func(current.value)
      clone.add(current.value)
    end
    current = current.next
  end
end

fn LinkedList.reduce(func, acc)
  let current = this.head
  repeat while current
    acc = func(acc, current.value)
    current = current.next
  end
  return acc
end

fn LinkedList.every(func)
  let current = this.head
  repeat while current
    if !func(current.value)
      return false
    end
    current = current.next
  end
  return true
end
