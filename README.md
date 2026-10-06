# Factory linter

[![CI](https://github.com/maranqz/gofactory/actions/workflows/ci.yml/badge.svg)](https://github.com/maranqz/gofactory/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/maranqz/gofactory)](https://goreportcard.com/report/github.com/maranqz/gofactory?dummy=unused)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
[![Coverage Status](https://coveralls.io/repos/github/maranqz/gofactory/badge.svg?branch=main)](https://coveralls.io/github/maranqz/gofactory?branch=main)

The linter checks that the Structures are created by the Factory, and not directly.

The checking helps to provide invariants without exclusion and helps avoid creating an invalid object.

## Protection scope

By default, gofactory protects only types from your **current module**: a package belongs to
the current module if its import path equals the module path or starts with the module path
plus `/`, so a nested module under the same path (its own `go.mod`, but still under your
module's path) counts too. Stdlib and third-party types (`strings.Builder{}`, `sync.WaitGroup{}`,
`http.Header{}`, `time.Duration(5)`) are silent, and so are a `go.work` sibling module's types —
their import path is neither your module path nor under it. Bring a sibling module, or any other
package, into scope with `--packageGlobs='example.com/sibling/**'` (an exact path does not match
yet). Until fences land, code inside a `--packageGlobs` package may itself bypass any factory, so
don't reuse this setting when linting the sibling module.

Running without a module (GOPATH, Bazel's `nogo`, or a list of `.go` files instead of packages)
falls back to the previous behaviour: every package other than the current one is protected.

Within scope, gofactory is **strict by default**: every bypass of a protected type is reported,
whether or not the type has a factory. See [the ADR](docs/adr/0002-strict-default-and-module-scope.md)
for why.

## Usage

### Installation

    go install github.com/maranqz/gofactory/cmd/gofactory@latest

### Options

- `--packageGlobs` – list of glob packages, which can create structures without factories inside the glob package. 
By default, types from the current module must come from their factories (see [Protection scope](#protection-scope)), [tests](testdata/module/packageGlobs).
- `--packageGlobsOnly` – use a factory to initiate a structure for glob packages only, 
[tests](testdata/module/packageGlobsOnly). Doesn't make sense without `--packageGlobs`.
- `--ignoreTypes` – list of qualified name globs (`import/path.Name`) for types that may be created
without a factory everywhere, not just inside a fence, [tests](testdata/module/ignoreTypes).
For example, `mymod/geo.*` matches every type of package `mymod/geo`, and `mymod/a.Pair` matches
the generic `Pair` with any type arguments. Type arguments in a glob are not supported, so a single
instantiation can't be ignored: `mymod/a.Pair[bool, bool]` matches no `Pair`.
`*` stays within one path segment (it crosses `.` but not `/`) and `**` also crosses `/`:
`mymod/*` matches `mymod/a.T` but not `mymod/a/b.T`; `mymod/**` matches both.
Repeat the flag to give several globs.
See [Directives](#directives) for the equivalent `//gofactory:ignore` comment.
- `--zeroValues` – off by default; report a zero value of a protected type as a bypass too, not just a literal,
conversion or `new`, [tests](testdata/module/zeroValues).
  - A package-level `var x T` is always reported.
  - A local `var x T` and a named result are decided by their first interaction anywhere in the function: a
    whole-value assignment (`x = …`, `x, err = …`) or `&x` passed to a call is silent; anything else, such as a
    read, a field access, a method call or `return x`, is reported. A naked `return` counts as an interaction with
    every named result. The function is read top to bottom, except that an assignment's right-hand side and a
    range expression count before the assignment itself (so `x = x.Paid()` and `for _, x = range x.Items()` are
    reported), and a `for` loop's body counts before its post statement. A function literal counts where it is
    written, even under `defer`.
  - Pointers, slices, maps and chans are not followed, and `make([]T, n)` and arrays are not reported yet, pending
    fill analysis: `var a [N]T` and `make([]T, n)` stay silent. A defined type such as
    `type Grid [3]T` or `type Tags []T` is a protected type itself, so `var g Grid` is reported.
  - Goes through the same owner-package and fences policy as every other route.
- `--factoryPatterns` – extra factory-name regex, appended to the default `^New` pattern; repeatable,
e.g. `--factoryPatterns=^Make --factoryPatterns=^Restore`, [tests](testdata/module/factoryPatterns).
- `--useDefaultFactoryPattern` – recognise the default `^New` pattern, `true` by default; set to
`false` to drop it, following the append-plus-bool-to-drop-builtins convention also used by errcheck,
asasalint and canonicalheader. With `--useDefaultFactoryPattern=false` and no `--factoryPatterns` at
all, no factory is ever recognised, [tests](testdata/module/useDefaultFactoryPattern).
- `--onlyWithFactory` – report only types that have a factory accessible from the reported site, so a
team can adopt the linter gradually, starting from the types that already have one,
[tests](testdata/module/onlyWithFactory).

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

### Message format

Every factory-bypass diagnostic starts with the stable prefix `Use factory for pkg.T`; this prefix
is a public contract that does not change. A diagnostic about an unknown or misplaced
[directive](#directives) does not carry it. A golangci-lint `linters.exclusions.rules[].text` or
`severity.rules[].text` regex that matches the prefix without anchoring the end of the message (for
example `^Use factory for`) keeps matching; one anchored to the end of the old, suffix-less message
(`^Use factory for pkg\.T$`) stops matching once a factory suffix is appended.

When `T` has a factory the reported site can call, the message gets a suffix naming up to three of
them, `New…` first, in a deterministic order: `Use factory for order.Order (order.NewOrder)`. A type
with no accessible factory keeps the bare prefix. A zero value reported under `--zeroValues` adds
`: zero value` right after the prefix, before any factory list:
`Use factory for order.Order: zero value (order.NewOrder)`.

A factory is recognised automatically in `T`'s owner package when it is an exported function, or an
exported method of another type than `T`, that returns `T` or `*T` among its results, takes no `T` or
`*T` parameter, and is named `New…`. Generic instantiations count too, e.g. `NewBox[T]() Box[T]`.
Methods of `T` itself are never factories, so withers and clones are not suggested, and a method
declared on an interface type is never recognised either. A factory function is named `pkg.NewT` in
the suffix; a factory method of another type `U` is named `pkg.U.NewT`, with a generic `U` rendered
without its type arguments, [tests](testdata/module/factories).

`--factoryPatterns` and `--useDefaultFactoryPattern` widen or replace which names count: a candidate
is recognised as soon as its name matches any configured pattern, default `^New` included unless
dropped. A `New…` candidate still sorts first in the suggestion even when another pattern also
matches it.

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
          zero-values: false
          factory-patterns:
            - "^Make"
          use-default-factory-pattern: true
          only-with-factory: false
```

- `package-globs` – equivalent to `--packageGlobs`.
- `package-globs-only` – equivalent to `--packageGlobsOnly`.
- `ignore-types` – equivalent to `--ignoreTypes`.
- `zero-values` – equivalent to `--zeroValues`.
- `factory-patterns` – equivalent to `--factoryPatterns`.
- `use-default-factory-pattern` – equivalent to `--useDefaultFactoryPattern`.
- `only-with-factory` – equivalent to `--onlyWithFactory`.

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
3. Unnamed composite literal implicitly converted to a named type, `var s nested.Struct = struct{ Field int }{-1}`, [example](testdata/module/unimplemented/implicit.go).
4. Conversion of an untyped non-constant expression, `nested.MyInt(1 << n)` or `nested.Flag(a == b)`, [example](testdata/module/unimplemented/untyped.go).
5. Type parameter whose constraint admits a single protected type, `func F[T nested.Struct]() T { return T{} }`, [example](testdata/module/unimplemented/typeparam.go).
6. `--zeroValues` reports a field-by-field fill after `var` (the first field write is the first interaction), but not
   elements filled after `make([]T, n)` or in arrays, which wait for fill analysis, [example](testdata/module/unimplemented/fill.go); use
   [gopublicfield](https://github.com/maranqz/gopublicfield) to prevent that.

## TODO

### Possible Features

1. Catch nested struct in the same package, [example](testdata/module/unimplemented/local/nested_struct.go).
   ```go
   return Struct{
       Other: OtherStruct{}, // want `Use factory for nested.Struct`
   }
   ```

### Features that are difficult to implement and unplanned

1. Reusing a protected type's underlying layout, including its unexported fields, to silently skip an invariant the factory establishes on them, [example](testdata/module/unimplemented/underlying.go). A local `type D ext.T` belongs to the current package, so it only hides an invariant when `ext.T` has unexported fields. Protection is keyed by type identity, not structural layout, so this needs a different detection model than anything currently planned.