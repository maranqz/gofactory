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
- `--ignoreTypes` – list of qualified name globs (`import/path.Name`) for types that may be created
without a factory everywhere, not just inside a fence, [tests](testdata/module/ignoreTypes).
`*` stays within one path segment and may cross `.`; `**` also crosses `/`. For example,
`mymod/*` matches `mymod/a.T` but not `mymod/a/b.T`, and `mymod/**` matches both but not `mymod.T`.
Repeat the flag to give several globs.
See [Directives](#directives) for the equivalent `//gofactory:ignore` comment.

### Directives

A `//gofactory:` comment, written in a declaration's doc comment (the comment group directly above
it) with no space after the slashes (like `//go:build`), marks that declaration for the linter.
It takes effect in every package and module that imports the declaration, not just the one that
writes it: gofactory exports directives as [analysis facts](https://pkg.go.dev/golang.org/x/tools/go/analysis#Fact),
so an importing package's analysis sees them even though it never parses the file that carries the
comment.

- `//gofactory:ignore`, in the doc comment of a single top-level type definition (not an alias,
  not a trailing comment, not above a `type ( … )` group of several types), takes that type out of
  protection: nothing in any package needs a factory to obtain a value of it, on any bypass route.
  It is the directive form of `--ignoreTypes`, [tests](testdata/module/directive).

  ```go
  //gofactory:ignore
  type Point struct {
      X, Y int
  }
  ```

`//gofactory:factory` (on a function or method) and `//gofactory:trusted` (on a function, a method,
or in a package's doc comment) are reserved: they are placement-checked like `ignore` but have no
effect yet.

An unknown directive name, or a known one in the wrong place (for example `//gofactory:ignore` on a
function or an alias), is reported as a diagnostic at the comment, so a typo does not silently
disable protection.

### golangci-lint module plugin

gofactory can also run inside golangci-lint as a [module plugin](https://golangci-lint.run/docs/plugins/module-plugins/),
without waiting for it to be merged into golangci-lint itself.

It needs golangci-lint v2.14.0 or newer: directives travel between packages as analysis facts, and
the fact-cache fixes they rely on landed in v2.13 and v2.14.

Build a custom golangci-lint binary that includes gofactory with a `.custom-gcl.yml`:

```yaml
version: v2.14.0
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
          ignore-types:
            - "mypkg.Point"
```

- `package-globs` – equivalent to `--packageGlobs`.
- `package-globs-only` – equivalent to `--packageGlobsOnly`.
- `ignore-types` – equivalent to `--ignoreTypes`.

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

None of `testdata/module/unimplemented/` is wired into `analysistest` (it's
excluded from `lint_test.go`'s test packages), so the `// want` comments in
files under `testdata/module/unimplemented/` document the diagnostic we want
once/if a case is handled — they are not asserted against the linter's
actual output. Every case added there must carry such a `// want` comment.

1. Buffered channel. You can initialize struct in line `v, ok := <-bufCh` [example](testdata/module/unimplemented/chan.go).
2. Local initialization, [example](testdata/module/unimplemented/local/).
3. Named return. If you want to block that case, you can use [nonamedreturns](https://github.com/firefart/nonamedreturns) linter, [example](testdata/module/unimplemented/named_return.go).
4. Unnamed composite literal implicitly converted to a named type, `var s nested.Struct = struct{ Field int }{-1}`, [example](testdata/module/unimplemented/implicit.go).
5. Conversion of an untyped non-constant expression, `nested.MyInt(1 << n)` or `nested.Flag(a == b)`, [example](testdata/module/unimplemented/untyped.go).
6. Type parameter whose constraint admits a single protected type, `func F[T nested.Struct]() T { return T{} }`, [example](testdata/module/unimplemented/typeparam.go).
7. var declaration, `var initilized nested.Struct` gives structure without factory, [example](testdata/module/unimplemented/var.go).
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

1. Reusing a protected type's underlying layout, including its unexported fields, to silently skip an invariant the factory establishes on them, [example](testdata/module/unimplemented/underlying.go). A local `type D ext.T` belongs to the current package, so it only hides an invariant when `ext.T` has unexported fields. Protection is keyed by type identity, not structural layout, so this needs a different detection model than anything currently planned.