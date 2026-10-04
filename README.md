# Factory linter

[![CI](https://github.com/maranqz/gofactory/actions/workflows/ci.yml/badge.svg)](https://github.com/maranqz/gofactory/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/maranqz/gofactory)](https://goreportcard.com/report/github.com/maranqz/gofactory?dummy=unused)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
[![Coverage Status](https://coveralls.io/repos/github/maranqz/gofactory/badge.svg?branch=main)](https://coveralls.io/github/maranqz/gofactory?branch=main)

The linter checks that the Structures are created by the Factory, and not directly.

The checking helps to provide invariants without exclusion and helps avoid creating an invalid object.


## Usage

### Installation

    go install github.com/maranqz/gofactory/cmd/gofactory@latest

### Options

- `--packageGlobs` – list of glob packages, which can create structures without factories inside the glob package. 
By default, all structures from another package should be created by factories, [tests](testdata/module/packageGlobs).
- `--packageGlobsOnly` – use a factory to initiate a structure for glob packages only, 
[tests](testdata/module/packageGlobsOnly). Doesn't make sense without `--packageGlobs`.

### golangci-lint module plugin

gofactory can also run inside golangci-lint as a [module plugin](https://golangci-lint.run/docs/plugins/module-plugins/),
without waiting for it to be merged into golangci-lint itself.

Build a custom golangci-lint binary that includes gofactory with a `.custom-gcl.yml`:

```yaml
version: v2.12.2
plugins:
  - module: github.com/maranqz/gofactory
    import: github.com/maranqz/gofactory
    version: latest
```

```shell
golangci-lint custom
```

This writes `./custom-gcl`; run it in place of `golangci-lint`, e.g. `./custom-gcl run ./...`.

Then enable it in `.golangci.yml`, with settings in kebab-case mirroring the command-line flags:

```yaml
version: "2"
linters:
  enable:
    - gofactory
  settings:
    custom:
      gofactory:
        type: module
        settings:
          package-globs:
            - "mypkg/internal/**"
          package-globs-only: false
```

- `package-globs` – equivalent to `--packageGlobs`.
- `package-globs-only` – equivalent to `--packageGlobsOnly`.

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

1. Buffered channel. You can initialize struct in line `v, ok := <-bufCh` [example](testdata/module/unimplemented/chan.go).
2. Local initialization, [example](testdata/module/unimplemented/local/).
3. Named return. If you want to block that case, you can use [nonamedreturns](https://github.com/firefart/nonamedreturns) linter, [example](testdata/module/unimplemented/named_return.go).
4. var declaration, `var initilized nested.Struct` gives structure without factory, [example](testdata/module/unimplemented/var.go).
 To block that case, you can use [gopublicfield](github.com/maranqz/gopublicfield) to prevent fill of structure fields.

## TODO

### Possible Features

1. Catch nested struct in the same package, [example](testdata/module/unimplemented/local/nested_struct.go).
   ```go
   return Struct{
       Other: OtherStruct{}, // want `Use factory for nested.Struct`
   }
   ```
2. Resolve false negative issue with `var declaration`.

### Features that are difficult to implement and unplanned

1. Type assertion, type declaration and type underlying, [tests](testdata/module/simple/type_nested.go.skip).