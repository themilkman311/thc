# THC

A small (**t**)ype-safe, (**h**)eterogeneous (**c**)ontainer. It allows you to store values, retrieve those values with typed keys, and delete stored values safely.

Module: `github.com/themilkman311/thc`

```
NewTHC() thc_container
Store[T any](container *thc_container, input T) (thc_key[T], error)
Fetch[T any](container *thc_container, key thc_key[T]) (T, error)
Update[T any](container *thc_container, key thc_key[T], input T) error
Remove[T](container *thc_container, key *thc_key[T]) error
```

Example

```go
package main

import (
    "fmt"
    "github.com/themilkman311/thc"
)

func main() {
    // Create a new container
    c := thc.NewTHC()

    // Store a string (or anything) in the container, get a key
    k, _ := thc.Store(&c, "hello, world")

    // Use the key to Fetch the value back
    v, _ := thc.Fetch(&c, k)
    fmt.Println("value:", v)

    // Update value (must be same type)
    thc.Update(&c, k, "goodbye, world")

    // Delete value and invalidate key
    if err := thc.Remove(&c, &k); err != nil {
        panic(err)
    }
}
```

Notes and design

- Keys only work on containers that created them.
- Attempting to store a container within itself will result in error.
