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
package, into scope with a fence, `--packageGlobs` (see [Fences](#fences) below).

Running without a module (GOPATH, Bazel's `nogo`, or a list of `.go` files instead of packages)
falls back to the previous behaviour: every package other than the current one is protected.

Within scope, gofactory is **strict by default**: every bypass of a protected type is reported,
whether or not the type has a factory. See [the ADR](docs/adr/0002-strict-default-and-module-scope.md)
for why.

An untyped constant implicitly converted to a protected type is a bypass too, reported where it is
stored: a `var` declaration or assignment (`var st ext.Status = 3`, `st = 3`), a call argument
(`f(3)`, `append(s, 3)`), a `return` and a composite literal element (`[]ext.Status{3}`). Using one
is not: comparisons (`st == 3`), `case 3:`, arithmetic (`st + 1`, `st += 1`) and `const`
declarations stay silent, and so does a typed constant such as `ext.Active`. An explicit conversion
`ext.Status(3)` is reported once, as a conversion, [tests](testdata/module/implicitConstants).

## Fences

Each `--packageGlobs` pattern is its own **fence**: the set of packages matching it. If a
type's package lies in one or more fences, only code inside **all** of those fences may bypass
that type's factory — several fences no longer disable each other, [tests](testdata/module/twofences).
A type in no fence gains nothing from fences; module scope alone decides whether it's protected,
the same as if `--packageGlobs` were never given — except under `--packageGlobsOnly`, where such a
type is not protected at all. See [the ADR](docs/adr/0003-fence-intersection-rule.md) for why.

Nesting two fences guards the inner one more tightly. With `--packageGlobs='a/**'
--packageGlobs='a/domain/**'`, `a/domain`'s type lies in both fences, so only code that is itself
in both — i.e. also under `a/domain/**` — may bypass it. `a/infra`, under `a/**` but not
`a/domain/**`, cannot bypass `a/domain`'s factory; `a/domain`, under both, can still bypass
`a/infra`'s, since `a/infra`'s type lies only in the outer fence and `a/domain` satisfies that
one, [tests](testdata/module/nestedfence).

A fence can name any package, including stdlib and third-party ones that module scope leaves
unprotected, [tests](testdata/module/stdlibfence).

### Glob syntax

A `--packageGlobs` pattern is gitignore-like: it is compiled with `/` as the path separator and
matched against both the package path and the path plus `/`. An exact package path therefore
matches on its own, with no wildcard needed, and `*` does not cross a `/`, so `a/*` matches `a/b`
but not `a/b/c`; `a/**` matches both. The path-plus-`/` match also means `a/*` and `a/**` match
`a` itself, not just what's inside it, unlike a `.gitignore` pattern, [tests](testdata/module/globsyntax).
Unlike `.gitignore`, a leading `/` has no "from the root" meaning — a Go package path never
starts with `/` — so it is rejected as a configuration error rather than silently matching
nothing.

### Migrating from pre-fence globs

Before fences, a package matching **any** `--packageGlobs` pattern was exempt from every check,
for every type, anywhere — this held even for a single pattern: code in that one subtree could
bypass any factory, not just the factories of types that are themselves in a fence. An exact path
like `pkg` never matched, `pkg/` matched only `pkg`, `pkg/**` matched `pkg` and its subpackages,
and `*` crossed `/`.

Under the intersection rule, the factory of a type in fences may be bypassed only by code inside
all of them, and a type in no fence is checked as if `--packageGlobs` were never given — or, under
`--packageGlobsOnly`, not protected at all. So a fence no longer exempts its code from checks on
types outside it, however many `--packageGlobs` patterns you pass. The migration is to put the
bypassing and the bypassed packages in one fence, e.g.
`--packageGlobs='{app/infra/**,app/domain/**}'` in place of separate `app/infra/**` and
`app/domain/**` patterns (gobwas/glob brace syntax), so that code in either package still lies
inside the same fence as the other's types. `--trusted` (see [Trusted code](#trusted-code) below)
is the direct replacement for "this code may bypass anything, anywhere." Also replace any `*` you
relied on crossing `/` with `**`. `pkg/` and `pkg/**` still match the same packages as before, and
`pkg` alone now matches exactly `pkg`.

### Recipe: protecting `go.work` sibling modules

A `go.work` sibling module's types are outside the current module, so they're silent by default.
Fence the sibling module's path with a trailing `/**` to protect every one of its packages, not
just its root package:

```
--packageGlobs='example.com/sibling/**'
```

This also matches the sibling module's root package itself, since a pattern is matched against
both the path and the path plus `/` (see [Glob syntax](#glob-syntax) above). Code inside the
current module is not in that fence, so it is blocked from bypassing any of the sibling's
factories; code inside the sibling module itself still may bypass any of them, which is looser
than module scope's per-package strictness — don't reuse this setting when linting the sibling
module itself, [tests](testdata/module/siblingfence).

## Owner package

By default, a package may bypass the factory of its own exported or unexported types anywhere in
its own code — that is what "owner package" means throughout this README. `--ownPackage` (off by
default) tightens that for exported types: inside the owner package, the factory of an **exported**
protected type may then be bypassed only inside a **producer** — a top-level function or method
whose results include:

- the type itself, `T`, or `*T`;
- a named interface `T` implements — not `any` or `interface{}`, which unalias to an unnamed
  interface and so never qualify, but a user-declared named interface does, even an empty one,
  since every type trivially implements it;
- a container holding `T` at any depth: a slice, an array, a map key or value, a `chan`, or an
  `iter.Seq`/`iter.Seq2`.

A method of a value object that returns the same type, a wither such as `func (o Order) Paid()
Order`, is a producer by the same `T` rule; so is a method of some other type that returns `T`, and
a function whose only connection to `T` is a parameter is not. A closure is judged by its enclosing
top-level declaration, not by its own signature, so a helper closure written inside a producer may
still bypass, while one written inside a non-producer is reported even if the closure itself
returns `T`. Every bypass route is checked — literal, conversion, `new` and implicit constant
conversion alike — except a `const` declaration in the owner package, at package scope or local to
a function, which this setting always leaves alone; elsewhere, a conversion in a `const`
declaration, such as `const C = order.Status(1)`, is still reported. A `_test.go` file is exempt, so
test code can still build fixtures directly. A type declared inside a function body is unaffected
too, even when capitalised, since Go exports only package-scope identifiers; no producer could ever
name such a type in its signature anyway. An unexported protected type is unaffected: outside code
could never name it anyway, so the owner-package rule for it stays as unrestricted as it always
was. With `--packageGlobsOnly`, `--ownPackage` tightens only the fence packages; a type outside
every fence stays unprotected, as `--packageGlobsOnly` already makes it. `--ownPackage` and fences
are otherwise independent rules: a fence still decides who outside the owner package may bypass a
type's factory, unchanged by this setting, [tests](testdata/module/ownPackage).

## Trusted code

Infrastructure code, such as a repository reconstituting an aggregate from storage, legitimately
needs to bypass a protected type's factory: a factory shaped for fresh input usually can't accept a
persisted row, field by field, the way reconstitution needs to. If the domain package can offer a
factory that takes a persisted row instead, such as `order.Restore`, declare it rather than trusting
the repository; see [Declared factories](#declared-factories) and its [reconstitution
recipe](#recipe-reconstitution-through-a-factory). When no such factory fits, trusted code may
bypass any protected type through every route, in every mode — module scope, fences,
`--onlyWithFactory`, `--zeroValues` and `--ownPackage` included. Trust is checked before any mode,
so a mode added later respects it automatically.

Mark it one of three ways:

- `//gofactory:trusted` in the doc comment of a function or a method trusts that function or
  method, and any closure written inside it.
- `//gofactory:trusted` in a package's doc comment trusts the whole package: every function,
  method and package-level var in it, without marking each one,
  [tests](testdata/module/directive) (`trustedfunc/` and `trustedpkg/`).
- `--trusted`, a repeatable glob, trusts code by name instead of by directive — a package path or a
  qualified function/method name — to trust a whole subtree at once (e.g.
  `--trusted='example.com/app/infra/**'`), or to keep every trust decision in your linter
  configuration instead of next to the code, [tests](testdata/module/trusted). See
  [Options](#options) for its glob syntax.

Code outside trusted code is still checked as usual.

### Recipe: trusted-repository reconstitution

```go
package postgres

import "example.com/app/domain/order"

type OrderRepository struct{ /* ... */ }

//gofactory:trusted
func (r *OrderRepository) Load(id string) (*order.Order, error) {
    row, err := r.queryRow(id)
    if err != nil {
        return nil, err
    }

    return &order.Order{ID: row.ID, Status: order.Status(row.Status)}, nil
}
```

`Load` may build `order.Order` directly from `row` on every route, while every other caller still
needs `order.NewOrder` or whichever factory the domain package declares. Prefer
`--trusted='example.com/app/infra/postgres.OrderRepository.Load'` over the directive to keep the
trust decision in your linter configuration rather than next to the code, or
`--trusted='example.com/app/infra/**'` to trust every repository under `infra` at once.

## Usage

### Installation

    go install github.com/maranqz/gofactory/cmd/gofactory@latest

### Options

- `--packageGlobs` – repeatable; each occurrence is its own fence (see [Fences](#fences)),
[tests](testdata/module/packageGlobs).
- `--packageGlobsOnly` – protect exactly the fence packages named by `--packageGlobs`, instead of
every current-module type, [tests](testdata/module/packageGlobsOnly). A configuration error
without at least one `--packageGlobs` pattern.
- `--ignoreTypes` – repeatable; a qualified name glob (`import/path.Name`) for types that may be
created without a factory everywhere, not just inside a fence, [tests](testdata/module/ignoreTypes).
For example, `mymod/geo.*` matches every type of package `mymod/geo`, and `mymod/a.Pair` matches
the generic `Pair` with any type arguments. Type arguments in a glob are not supported, so a single
instantiation can't be ignored: `mymod/a.Pair[bool, bool]` matches no `Pair`.
`*` stays within one path segment (it crosses `.` but not `/`) and `**` also crosses `/`:
`mymod/*` matches `mymod/a.T` but not `mymod/a/b.T`; `mymod/**` matches both.
Name the type where it is defined: an alias's name does not match, and neither does a bare package
path.
Repeat the flag to give several globs: a comma does not separate them. As with `--packageGlobs`, an
empty glob or one starting with `/` is a configuration error.
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
  - An unset by-value field of a protected type is reported at any depth through struct fields and embedding, with
    the field path in the message, e.g. `Use factory for pkg.T: zero value in V.W.S`. A composite literal that
    leaves such a field unset is reported (an empty literal, a partially keyed one for its unset fields only, or an
    embedded field left out), and so is `new(W)`, which sets none; a fully positional literal sets every field, so
    it reports nothing. A zero-valued variable of the enclosing type is reported the same way, under the first-interaction rule above, judged once for
    the whole variable: writing one field, or passing that field's address to a call, still reports every protected
    field of the variable at that write, the written one included, even though the write itself used a factory; only
    a whole-value assignment or `&x` on the variable itself is silent. A path through an unexported field of another
    package's type is reported too, though the reporting code can't set that field: the type's own methods can still
    hand out its zero value, so the fix lies with a factory for the enclosing type. Fields behind a pointer, slice,
    map, chan or array are not followed, and an array element in a field path is deferred like every other array,
    pending fill analysis. Field-path analysis uses a per-type cache, benchmarked by `BenchmarkFieldPaths`.
  - Goes through the same owner-package and fences policy as every other route.
- `--ownPackage` – off by default; restrict an exported protected type's owner package to
producers — see [Owner package](#owner-package), [tests](testdata/module/ownPackage).
- `--factoryPatterns` – extra factory-name regex, appended to the default `^New` pattern; repeatable,
e.g. `--factoryPatterns=^Make --factoryPatterns=^Restore`, [tests](testdata/module/factoryPatterns).
- `--useDefaultFactoryPattern` – recognise the default `^New` pattern, `true` by default; set to
`false` to drop it, following the append-plus-bool-to-drop-builtins convention also used by errcheck,
asasalint and canonicalheader. With `--useDefaultFactoryPattern=false` and no `--factoryPatterns` at
all, no factory is ever recognised, [tests](testdata/module/useDefaultFactoryPattern).
- `--onlyWithFactory` – report only types that have a factory accessible from the reported site, so a
team can adopt the linter gradually, starting from the types that already have one,
[tests](testdata/module/onlyWithFactory).
- `--factories` – repeatable; a qualified function- or method-name glob (`import/path.Func` or
`import/path.Type.Method`) declaring a factory outside the usual recognition rule,
[tests](testdata/module/declaredFactories). See [Declared factories](#declared-factories) below,
including its reach limit.
`*` stays within one path segment: as in `--ignoreTypes`, it crosses `.` but not `/`, so
`mymod/order.*` also matches a method of any type in `mymod/order`. `**` also crosses `/`.
Repeat the flag to give several globs: a comma does not separate them. An empty glob, or one
starting with `/`, is a configuration error.
- `--trusted` – repeatable; a package-path glob (matched like `--packageGlobs`, see [Glob
syntax](#glob-syntax)) or a qualified function- or method-name glob (`import/path.Name`,
`import/path.Type.Method`, matched like `--ignoreTypes`) for code that may bypass any protected
type's factory through any route, in every mode, [tests](testdata/module/trusted). See [Trusted
code](#trusted-code) for the equivalent `//gofactory:trusted` directive and a reconstitution
recipe.
- `--crossPackageDirectives` – propagate `//gofactory:ignore` and `//gofactory:factory` to importing
packages and modules, `true` by default; see [Directives](#directives) below for what that means
(`//gofactory:trusted` never crosses a package boundary, so this setting does not concern it). Set
to `false` on a large monorepo to trade that off for less analysis of dependencies: with it off, the
analyzer declares no `FactTypes`, so the standalone `gofactory` command no longer parses,
type-checks and analyses dependencies from source (the `go` command still compiles them for their
export data), and golangci-lint no longer runs gofactory on them (it still parses and type-checks
them from source if another enabled linter uses facts). `go vet` runs the tool on, and type-checks,
every dependency either way, so there the setting saves only gofactory's own pass.
`//gofactory:ignore` and `//gofactory:factory` still take effect in the package that declares them,
and settings such as `--ignoreTypes` and `--packageGlobs` still apply everywhere, as does a
`--factories` match within its usual [reach](#declared-factories); only propagation of directives to
importers is turned off, [tests](testdata/module/crossPackageDirectives).

### Directives

A `//gofactory:` comment, written in a declaration's doc comment (the comment group directly above
it) with no space after the slashes (like `//go:build`), marks that declaration for the linter.
gofactory exports directives as [analysis facts](https://pkg.go.dev/golang.org/x/tools/go/analysis#Fact),
so a package that imports the one carrying the directive sees it even though it never parses the
file with the comment, not just the package that writes it. `//gofactory:ignore` reaches every
transitive importer this way; `//gofactory:factory` reaches only direct importers of its package
(see [Declared factories](#declared-factories) below). `//gofactory:trusted` needs no fact; see
below.

- `//gofactory:ignore`, in the doc comment of a single top-level type definition (not an alias,
  not a trailing comment, not above a `type ( … )` group of several types), takes that type out of
  protection: nothing in any package needs a factory to obtain a value of it, on any bypass route.
  It is the directive form of `--ignoreTypes`, [tests](testdata/module/directive). It takes effect
  in every package and module that imports the type, not just the one that writes it: gofactory
  exports it as an [analysis fact](https://pkg.go.dev/golang.org/x/tools/go/analysis#Fact), so an
  importing package's analysis sees it even though it never parses the file that carries the
  comment, unless `--crossPackageDirectives=false` turns that propagation off.

  ```go
  //gofactory:ignore
  type Point struct {
      X, Y int
  }
  ```

- `//gofactory:factory`, in the doc comment of a function or a method, declares it a factory. See
  [Declared factories](#declared-factories) below for its reach, which `--crossPackageDirectives=false`
  limits the same way it does `//gofactory:ignore`'s.

- `//gofactory:trusted`, in the doc comment of a function, a method or a package, trusts that
  declaration — see [Trusted code](#trusted-code). Unlike `ignore`, its effect never crosses a
  package boundary: it describes who may bypass a factory, not which type is exempt, and who is
  always decided inside the package being linted, so it needs no fact.

An unknown directive name, or a known one in the wrong place (for example `//gofactory:ignore` on a
function or an alias), is reported as a diagnostic at the comment, so a typo does not silently
disable protection.

### Declared factories

A function or a method is a factory of a type whether or not it is recognised by name pattern in
that type's owner package: mark it with a `//gofactory:factory` doc comment, or match it with a
`--factories` glob, and it becomes a factory of every protected type among its results (`T` or
`*T`; an `error` result is just ignored, not disqualifying), wherever it lives. This is how a
package that wraps generated code, say `pb`, can provide the factory for `pb.Order` from outside
`pb` itself, something the owner-package-only recognition rule can't do:

```go
package order

//gofactory:factory
func New(id string) *pb.Order {
    return &pb.Order{Id: id}
}
```

A declared factory may itself bypass the factories of the types it is a factory of, inside its own
body (including a closure it defines, which shares its enclosing top-level declaration's
permission), the same way a recognised factory may bypass its own type's factory in its owner
package; it gains no permission over any other type. This holds even when the type lies in a
fence (see [Fences](#fences)) that the factory's package is outside of,
[tests](testdata/module/declaredFactoriesFence): a declared factory is a statement about that
function, independent of where it lives. In a package that imports its package directly, calling
it is suggested in messages the same way a recognised factory is,
[tests](testdata/module/declaredFactories). There it also counts for `--onlyWithFactory`, even
with `--useDefaultFactoryPattern=false`, [tests](testdata/module/declaredFactoriesOnlyWithFactory).

**Reach.** A declared factory is known only to a package that imports its package directly, for a
function and a method alike. A `--factories` glob needs no fact: every package matches it against
its own functions and methods and those of the packages it imports directly. A
`//gofactory:factory` directive travels as a fact on a function object, and go/analysis's fact
machinery (`checker.exportedFrom`, `facts.Encode`) forwards a function's fact only to direct
importers, unlike a type's `//gofactory:ignore` fact, which every transitive importer sees. For a
method, both `checker.exportedFrom` and go vet's `facts.Encode` can hand its fact to more than the
direct importers — `checker.exportedFrom` over-approximates outright, and `facts.Encode` forwards
it whenever an importer's export data includes the declaring package; gofactory discards what
either adds, so every driver agrees on the same direct-importer rule. A package's `_test.go`
imports count only when its tests are analysed: the standalone CLI, which analyses a package both
without and with its tests, then prints such a line twice, the first time without the suggestion.

A package that uses the protected type without importing the declaring package directly doesn't
see that factory: it is left out of the message and doesn't count for `--onlyWithFactory`. A type
whose only factory it is gets the bare `Use factory for pkg.T` message there, or no report at all
under `--onlyWithFactory` — see [False Negative](#false-negative).

A `//gofactory:factory` directive on a function or method with no result that could ever be
protected (a named type other than a func type or an interface, as `T` or `*T`) is reported as a
diagnostic, since it would otherwise do nothing while looking like it did something. Whether a
result is protected in this module, or fenced, doesn't matter: `func Now() time.Time` is accepted.
A `--factories` glob that matches such a function is not an error: like `--ignoreTypes`, it may
simply match nothing relevant.

### Recipe: reconstitution through a factory

Creation and reconstitution are the same concept: a function that a repository calls to rebuild an
aggregate from storage, `order.Restore(id, status, paidAt)`, is just a factory, enabled the same
way any other unusually-named factory is — a `--factoryPatterns=^Restore` pattern if the convention
is project-wide, or a `//gofactory:factory` directive if `Restore` is just this one repository's
name for it:

```go
package order

//gofactory:factory
func Restore(id string, status Status, paidAt time.Time) *Order {
    return &Order{id: id, status: status, paidAt: paidAt}
}
```

When no such factory fits — the repository needs a bypass that a factory shaped for fresh input
couldn't offer — trust the repository instead; see [Trusted code](#trusted-code).

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
`Use factory for order.Order: zero value (order.NewOrder)`. A zero value reported through a field
path adds where the field is, e.g. `Use factory for order.Order: zero value in Shipment.Order (order.NewOrder)`.

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
          trusted:
            - "mypkg/infra/**"
          zero-values: false
          own-package: false
          factory-patterns:
            - "^Make"
          use-default-factory-pattern: true
          only-with-factory: false
          factories:
            - "mymod/order.Restore"
          cross-package-directives: true
```

- `package-globs` – equivalent to `--packageGlobs`.
- `package-globs-only` – equivalent to `--packageGlobsOnly`.
- `ignore-types` – equivalent to `--ignoreTypes`.
- `trusted` – equivalent to `--trusted`.
- `zero-values` – equivalent to `--zeroValues`.
- `own-package` – equivalent to `--ownPackage`.
- `factory-patterns` – equivalent to `--factoryPatterns`.
- `use-default-factory-pattern` – equivalent to `--useDefaultFactoryPattern`.
- `only-with-factory` – equivalent to `--onlyWithFactory`.
- `factories` – equivalent to `--factories`.
- `cross-package-directives` – equivalent to `--crossPackageDirectives`.

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
2. Local initialization inside the owner package: catching this needs `--ownPackage`,
   [tests](testdata/module/ownPackage).
3. Unnamed composite literal implicitly converted to a named type, `var s nested.Struct = struct{ Field int }{-1}`, [example](testdata/module/unimplemented/implicit.go).
4. Conversion of an untyped non-constant expression, explicit or implicit, `nested.MyInt(1 << n)`, `nested.Flag(a == b)` or `var f nested.Flag = a == b`, [example](testdata/module/unimplemented/untyped.go).
5. Type parameter whose constraint admits a single protected type, `func F[T nested.Struct]() T { return T{} }` or `func F[T nested.MyInt]() T { return 3 }`, [example](testdata/module/unimplemented/typeparam.go).
6. `--zeroValues` reports a field-by-field fill after `var` (the first field write is the first interaction), but not
   elements filled after `make([]T, n)` or in arrays, which wait for fill analysis, [example](testdata/module/unimplemented/fill.go); use
   [gopublicfield](https://github.com/maranqz/gopublicfield) to prevent that.
7. A declared factory (`//gofactory:factory` or `--factories`) is suggested, and counted by
   `--onlyWithFactory`, only in a package that imports its package directly; a package that imports
   the declaring package only through another package, or not at all, never sees it,
   [example](testdata/module/unimplemented/visibility/).
8. Untyped constant stored by a channel send, as a map index key or by a range assignment,
   `ch <- 3`, `m[3] = v` or `for st = range 3`, [example](testdata/module/unimplemented/constant.go).

## TODO

### Features that are difficult to implement and unplanned

1. Reusing a protected type's underlying layout, including its unexported fields, to silently skip an invariant the factory establishes on them, [example](testdata/module/unimplemented/underlying.go). A local `type D ext.T` belongs to the current package, so it only hides an invariant when `ext.T` has unexported fields. Protection is keyed by type identity, not structural layout, so this needs a different detection model than anything currently planned.