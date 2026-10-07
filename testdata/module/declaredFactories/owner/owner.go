// Package owner holds protected types with no recognised factory of their
// own: every factory suggested for them, elsewhere in this test tree, comes
// from a declared factory living in a different package (other, flagged),
// except for RedundantlyDeclared, declared where it already matches the
// default pattern.
package owner

// Solo's only factory is other.Build, a //gofactory:factory function in
// another package.
type Solo struct{}

// ViaMethod's only factory is other.Repo.Restore, a //gofactory:factory
// method in another package.
type ViaMethod struct{}

// Both and Partner share one declared factory, other.BuildBoth, which
// returns both.
type Both struct{}

type Partner struct{}

// WithErr's declared factory, other.BuildWithErr, also returns an error:
// that extra result is ignored, not disqualifying.
type WithErr struct{}

// ClosureTarget's declared factory, other.BuildViaClosure, bypasses it
// through a closure; the closure takes its enclosing top-level
// declaration's permission.
type ClosureTarget struct{}

// Flagged's only factory, flagged.MakeFlagged, is declared through
// -factories rather than a directive.
type Flagged struct{}

// FlaggedMethod's only factory, flagged.Box.RestoreFlaggedMethod, is
// declared through -factories too, over a method this time.
type FlaggedMethod struct{}

// FlaggedExact's only factory, flagged.Box.RestoreExact, is declared
// through a -factories glob with no wildcard, receiver type included.
type FlaggedExact struct{}

// PtrBuilt's declared factory, other.BuildPtr, returns a pointer.
type PtrBuilt struct{}

// Globbed's declared factory, sub.MakeGlobbed, is matched by a -factories
// wildcard; deep.MakeGlobbed, matched against the identical glob but one
// path segment further, is not (glob/sub, glob/sub/deep).
type Globbed struct{}

// NoDirective has no factory anywhere, declared or recognised.
type NoDirective struct{}

// Inaccessible's declared factory, other.inaccessibleRepo.Restore, is a
// method of an unexported type: outside other, nothing can name the
// receiver, so it is never suggested from there, the same accessibility
// rule a recognised factory follows.
type Inaccessible struct{}

// Unmarked has no factory anywhere either, like NoDirective.
type Unmarked struct{}

// RedundantlyDeclared already matches the default ^New pattern; marking
// its factory with the directive too must not list it twice.
type RedundantlyDeclared struct{}

//gofactory:factory
func NewRedundantlyDeclared() RedundantlyDeclared { // want NewRedundantlyDeclared:"gofactory:factory"
	return RedundantlyDeclared{}
}
