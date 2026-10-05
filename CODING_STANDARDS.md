# Coding standards

## Comments

A comment says what the code can't: a hidden constraint, an invariant, a workaround for a specific
bug, behaviour that would surprise the next reader. If deleting a comment would leave that reader no
more confused, delete it.

- **What the code does** is for its names and its body to say. A comment that retells the body is
  noise; one that a better name would make unnecessary asks for the rename instead.
- **Where a change came from** (the issue, the review, the caller it was written for) and why it is
  correct go in the commit message and the PR description. In the code they go stale as soon as the
  code around them changes.
- **Doc comments are comments.** An exported identifier gets the doc comment Go expects: what a
  caller needs, in a sentence or two. An unexported helper whose name and body say it all gets none.

This one earns its place, because no reader could see the constraint in the code:

```go
// Pass.Module.Main cannot find the root module: go vet's unitchecker before
// Go 1.27 leaves it unset, and a go.work build sets it on every module.
```

This one does not. Its first sentence retells the four lines under it, and the last is about how
fences treat a `go.work` sibling, which belongs with the fences:

```go
// belongsToModule reports whether pkgPath is the module at modulePath, or
// nested under it, so that a nested module in the same repository counts
// as the current module too. A go.work sibling module, whose path does not
// start with modulePath, does not belong, and is silent until a fence
// names it.
func belongsToModule(pkgPath, modulePath string) bool {
	if pkgPath == modulePath {
		return true
	}

	return strings.HasPrefix(pkgPath, modulePath+"/")
}
```

Without the comment, the function says the same.
