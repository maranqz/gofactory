# Changelog

## v1.1.0

Detection is rebuilt on go/types, protection is scoped to the current module, package globs become
fences with gitignore-like syntax, and directives and settings let you describe your design. The
module path `github.com/maranqz/gofactory` and `NewAnalyzer()` are unchanged, so imports and
tooling keep working, but the set of reports changes: read the breaking changes below before
upgrading. [ADR 0005](docs/adr/0005-breaking-changes-as-v1-1-0.md) explains why this is v1.1.0 and
not v2.

### Breaking changes

- **Module scope by default.** Only types from the current module are protected: a package whose
  import path equals the module path or starts with it plus `/`, nested modules included. Stdlib,
  third-party and `go.work` sibling types are no longer reported (`strings.Builder{}`,
  `sync.WaitGroup{}`, `time.Duration(5)`). To keep protecting one, name it in a fence, e.g.
  `-packageGlobs='example.com/sibling/**'`. Without a module (GOPATH, Bazel `nogo`), every package
  other than the current one is still protected, as before.
  See [Protection scope](README.md#protection-scope) and
  [ADR 0002](docs/adr/0002-strict-default-and-module-scope.md).
- **Fence semantics.** Each `-packageGlobs` pattern is now its own fence. A type in one or more
  fences may be bypassed only by code inside all of them; a type in no fence is checked as if no
  globs were given (unprotected under `-packageGlobsOnly`). Before, a package matching any glob was
  exempt from every check, for every type, so globs disabled each other's protection. A setup that
  relied on one subtree bypassing another's factories needs both in one fence, or `-trusted` for
  code that may bypass anything.
  See [Fences](README.md#fences) and [ADR 0003](docs/adr/0003-fence-intersection-rule.md).
- **Glob syntax.** Package globs are gitignore-like: an exact package path matches on its own, `*`
  no longer crosses `/` (use `**`). An empty pattern, or one starting with `/`, is a configuration
  error.
  **Migrate your globs with the [glob migration guide](README.md#migrating-from-pre-fence-globs).**
- **`-packageGlobsOnly` without `-packageGlobs` is a configuration error.** It used to do nothing,
  silently.
- **New bypass routes.** Detection resolves types instead of matching syntax, so more is reported:
  - literals through local aliases, generic aliases and aliases of instantiations
    (`type IntG = ext.Generic[int]`), dot-imports, generics with several type parameters
    (`ext.Pair[int, string]{}`), elided and keyed elements at any depth (`[][]ext.Struct{{{}}}`,
    `[]ext.Struct{0: {}}`), parenthesised element types, and each element of a named container
    (`ext.Structs{{}}`, `type Ss []ext.Struct; Ss{{}}`);
  - parenthesised conversions, `(ext.Struct)(l)` and `(*ext.Struct)(&l)`;
  - `new(T)`;
  - an untyped constant implicitly converted to a protected type where it is stored:
    `var st ext.Status = 3`, `st = 3`, `f(3)`, `return 3`, `[]ext.Status{3}`;
  - with `-zeroValues` (off by default), zero values: `var x T`, named results and unset fields.
- **Fewer reports where no factory could help.** Conversions that create nothing (`ext.Mp(nil)`,
  to an identical type, to an interface, to `unsafe.Pointer`), func and interface types
  (`http.HandlerFunc(f)`), files with a `// Code generated … DO NOT EDIT.` header and the test main
  `go test` synthesizes are no longer reported.
- **Aliases resolve to the aliased type.** A literal or conversion through an alias is reported,
  and protected or not, as the type it aliases: with `type Alias = other.Real`, `ext.Alias{}` is
  now `Use factory for other.Real`, not `Use factory for ext.Alias`, and an alias of a stdlib type
  is silent. Update `text:` rules that name an alias.
- **Report position.** A report now starts at the beginning of the expression, not at the type
  name: for `_ = ext.Struct{}` it is at `ext`, not at `Struct`, so the line stays and the column
  moves left (column 10 becomes column 6 when the statement is indented by one tab).
- **Dependencies are analysed by default.** Directives travel as analysis facts, so the analyzer
  now runs over every dependency of the linted packages, which is slower than v1.0.x;
  `-crossPackageDirectives=false` restores the old cost and drops cross-package directives.
- **Message suffix.** Messages keep the `Use factory for pkg.T` prefix and may now add a suffix: the
  accessible factories, `Use factory for order.Order (order.NewOrder)`, and for a zero value the
  route and field path, `Use factory for order.Order: zero value in Shipment.Order`. An exclusion
  or severity regex anchored to the end of the old message (`^Use factory for pkg\.T$`) stops
  matching; `^Use factory for pkg\.T` keeps matching. See
  [Message format](README.md#message-format) and
  [ADR 0001](docs/adr/0001-stable-use-factory-for-prefix.md).
- **Requirements.** Go 1.26 or newer to build and install. golangci-lint v2.14.0 or newer for the
  module plugin.

### Added

- `-onlyWithFactory`: report only types that have a factory accessible from the reported site.
- Factory recognition: an exported `New…` function, or method of another type, in the owner package
  that returns `T` or `*T` is suggested in messages. `-factoryPatterns` adds name patterns and
  `-useDefaultFactoryPattern=false` drops `^New`.
- Declared factories: `//gofactory:factory` and `-factories`, in any package.
- Ignored types: `//gofactory:ignore` and `-ignoreTypes`.
- Trusted code: `//gofactory:trusted` on a function, a method or a package, and `-trusted`.
- Unknown or misplaced `//gofactory:` directives are reported.
- Directives propagate to importing packages and modules as analysis facts;
  `-crossPackageDirectives=false` turns that off and stops analysing dependencies.
- `-zeroValues`: zero-value bypasses, with the field path in the message.
- `-ownPackage`: inside the owner package, an exported type's factory may be bypassed only in a
  producer.
- `foo_test` counts as owner package `foo`.
- A golangci-lint module plugin with kebab-case settings mirroring the flags.

### Fixed

- The analyzer no longer writes to the log.
