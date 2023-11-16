# Factory linter

[![CI](https://github.com/maranqz/go-factory-lint/actions/workflows/ci.yml/badge.svg)](https://github.com/maranqz/go-factory-lint/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/maranqz/go-factory-lint)](https://goreportcard.com/report/github.com/maranqz/go-factory-lint?dummy=unused)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
[![Coverage Status](https://coveralls.io/repos/github/maranqz/go-factory-lint/badge.svg?branch=main)](https://coveralls.io/github/maranqz/go-factory-lint?branch=main)

The linter checks that the Structures are created by the Factory, and not directly.

The checking helps to provide invariants without exclusion and helps avoid creating an invalid object.


## Use

Installation

    go install github.com/maranqz/go-factory-lint/cmd/go-factory-lint@latest

### Options

- `--packageGlobs` - list of glob packages, which can create structures without factories inside the glob package. By default, all structures from another package should be created by factories, [tests](testdata/src/factory/packageGlobs).
    - `onlyPackageGlobs` - use a factory to initiate a structure for glob packages only, [tests](testdata/src/factory/onlyPackageGlobs).

## Example

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
package main

import (
	"fmt"

	"bad"
)

func main() {
	// Use factory for bad.User
	u := &bad.User{
		ID: -1,
	}

	fmt.Println(u.ID) // -1
	fmt.Println(u.CreatedAt) // time.Time{}
}

```

```go
package bad

import "time"

type User struct {
	ID        int64
	CreatedAt time.Time
}

var sequenceID = int64(0)

func NextID() int64 {
	sequenceID++

	return sequenceID
}


```

</td><td>

```go
package main

import (
	"fmt"

	"good"
)

func main() {
	u := good.NewUser()
	
	fmt.Println(u.ID)        // auto increment
	fmt.Println(u.CreatedAt) // time.Now()
}

```

```go
package user

import "time"

type User struct {
	ID        int64
	CreatedAt time.Time
}

func NewUser() *User {
	return &User{
		ID: nextID(),
		CreatedAt: time.Now(),
	}
}

var sequenceID = int64(0)

func nextID() int64 {
	sequenceID++

	return sequenceID
}

```

</td></tr>
</tbody></table>

## False Negative

Linter doesn't catch some cases.

1. Buffered channel. You can initialize struct in line `v, ok := <-bufCh` [example](testdata/src/factory/unimplemented/chan.go).
2. Local initialization. There is not single ability to understand which factory is right [example](testdata/src/factory/unimplemented/local.go).
3. Named return. If you want to block that case, you can use [nonamedreturns](https://github.com/firefart/nonamedreturns) linter, [example](testdata/src/factory/unimplemented/named_return.go).
4. var declaration, `var initilized nested.Struct` gives structure without factory, [example](testdata/src/factory/unimplemented/var.go).

## TODO

### Possible Features

1. Resolve false negative issue with `var declaration`.
2. Catch nested struct in the same package, [example](testdata/src/factory/unimplemented/nested_struct.go).
   ```go
   return Struct{
   	Other: OtherStruct{}, // want `Use factory for nested.Struct`
   }
   ```

### Features that are difficult to implement and unplanned

1. Type assertion, type declaration and type underlying, [tests](testdata/src/factory/simple/type_nested.go.skip).