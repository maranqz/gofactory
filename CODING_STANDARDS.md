# Coding standards

## Comments

A comment says what the code can't show the next reader: a constraint, an invariant, a workaround,
a surprise. What the code does is for its names; why it changed goes in the commit message and the
PR. An exported identifier gets a doc comment for its caller; an unexported helper that says it all
gets none.

This one earns its place, because no reader could see the constraint in the code:

```go
// Pass.Module.Main cannot find the root module: go vet's unitchecker before
// Go 1.27 leaves it unset, and a go.work build sets it on every module.
```

## Tests

Subtest names, including the test-table keys that become them, have no spaces: words are joined
with `_`. `go test` prints a subtest name with each space replaced by `_`, so a case written as
`"invalid glob"` fails as `TestConfigurationErrors/invalid_glob`, and searching the source for that
name finds nothing. A camelCase word, such as a flag name, stays as written: `packageGlobs`,
`packageGlobsOnly_without_globs`.

## Testdata

A case the linter defers to later work goes in `testdata/module/unimplemented/`, with a `// want` for the diagnostic it should get.
