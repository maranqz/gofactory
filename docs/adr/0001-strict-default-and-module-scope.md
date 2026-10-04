# Strict default, scoped to the current module

The go/types rewrite (#43) made every bypass of a protected type detectable, which made two
older defaults worth revisiting: should a type need an existing factory to be reported
(`-onlyWithFactory`), and should every other package stay protected the way it always was?

We decided strict mode (report a bypass whether or not the type has a factory) stays the
default, and narrowed what counts as protected instead: only types from the current module. A
report that ignores factory-less types would hide exactly the types a factory is missing for,
defeating the reason to migrate to the stronger detector; module scope removes stdlib and
dependency noise (`strings.Builder{}`, `sync.WaitGroup{}`, PR golangci/golangci-lint#4196's main
objection) without weakening what's reported about a developer's own code. Running without a
module (GOPATH, Bazel `nogo`) keeps the pre-#43 behaviour — every other package protected — so
the upgrade doesn't silently turn the linter off for those setups.

This is a breaking change, shipped as v1.1.0 rather than v2 (a separate decision, #43's
"breaking behaviour shipped as v1.1.0, not v2" candidate): existing users relying on types
outside the current module (stdlib, a dependency, a `go.work` sibling module) being reported
will see their diagnostics change on upgrade. `-onlyWithFactory`, planned in #43, will let teams
adopt the stronger detector gradually once it ships.

## Considered options

- **Keep protecting every package, as before.** Rejected: this is what produced the stdlib/
  dependency noise the rewrite set out to fix.
- **Default `-onlyWithFactory` to true.** Rejected: it would silently hide bypasses of types
  that have no factory yet, which is the situation most worth surfacing.
- **Module scope, strict by default.** Chosen, for the reasons above.
