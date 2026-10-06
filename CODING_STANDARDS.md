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

## Testdata

A case the linter defers to later work goes in `testdata/module/unimplemented/`, with a `// want` for the diagnostic it should get.
