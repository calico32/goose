type numeric = int | float

native type any
native type void
native type never


native type int
native operator int +(other: numeric)
native operator int +(other: numeric)
native operator int -(other: numeric)
native operator int *(other: numeric)
native operator int /(other: numeric)
native operator int %(other: numeric)
native operator int ==(other: any)
native operator int <=>(other: numeric)
native operator int !()
native operator int ~()
native operator int &(other: int)
native operator int |(other: int)
native operator int ^(other: int)
native operator int **(other: numeric)
native fn int.toString(base: int): string


native type float
native operator float +(other: numeric)
native operator float -(other: numeric)
native operator float *(other: numeric)
native operator float /(other: numeric)
native operator float %(other: numeric)
native operator float ==(other: any)
native operator float <=>(other: numeric)
native operator float !()
native operator float **(other: numeric)
native fn float.toString(): string


native type bool
native operator bool ==(other: any)
native fn bool.toString(): string


native type null
native operator null ==(other: any)
native operator null !()
native fn null.toString(): string


native type string
native operator string ==(other: any)
native operator string +(other: string)
native fn string.toString(): string
native fn string.split(sep: string): Array

native type Array
native operator Array ==(other: any)
native operator Array +(other: Array): Array
native fn Array.toString(): string
native fn Array.push(item: any): void
native fn Array.sort(comparator: fn(a: any, b: any) -> int,): void


native type Function


native fn string(obj: any): string
native fn int(obj: any): int
native fn float(obj: any): float
native fn bool(obj: any): bool

native fn len(obj: any): int
native fn sleep(ms: int): void
native fn print(...args: any): void
native fn printf(fmt: string, ...args: any): void
native fn exit(code: int): never

native fn indices(obj: any): Array
native fn keys(obj: any): Array
native fn values(obj: any): Array

native fn typeof(obj: any): string
