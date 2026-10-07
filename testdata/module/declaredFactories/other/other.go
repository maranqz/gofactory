// Package other declares factories of owner's types from outside owner,
// through //gofactory:factory; flagged.go in the sibling package covers
// -factories instead.
package other

import "factory/declaredFactories/owner"

// Build is owner.Solo's declared factory, living in another package. It
// may bypass Solo's own factory, since it is itself declared a factory of
// Solo, but not an unrelated type's: NoDirective still needs reporting
// here.
//
//gofactory:factory
func Build() owner.Solo { // want Build:"gofactory:factory"
	_ = owner.NoDirective{} // want `Use factory for owner.NoDirective$`

	return owner.Solo{}
}

type Repo struct{}

// Restore is owner.ViaMethod's declared factory, as a method of another
// type.
//
//gofactory:factory
func (Repo) Restore() owner.ViaMethod { // want Restore:"gofactory:factory"
	return owner.ViaMethod{}
}

// BuildBoth is a factory of both its protected results.
//
//gofactory:factory
func BuildBoth() (owner.Both, owner.Partner) { // want BuildBoth:"gofactory:factory"
	return owner.Both{}, owner.Partner{}
}

// BuildWithErr's error result is ignored: it is still a factory of
// WithErr.
//
//gofactory:factory
func BuildWithErr() (owner.WithErr, error) { // want BuildWithErr:"gofactory:factory"
	return owner.WithErr{}, nil
}

// BuildViaClosure exercises a closure's enclosing top-level declaration:
// the closure body shares BuildViaClosure's permission to bypass
// ClosureTarget's factory.
//
//gofactory:factory
func BuildViaClosure() owner.ClosureTarget { // want BuildViaClosure:"gofactory:factory"
	build := func() owner.ClosureTarget {
		return owner.ClosureTarget{}
	}

	return build()
}

// BuildPtr returns a pointer; &owner.PtrBuilt{} and new(owner.PtrBuilt) in
// it are both silent.
//
//gofactory:factory
func BuildPtr() (*owner.PtrBuilt, error) { // want BuildPtr:"gofactory:factory"
	_ = &owner.PtrBuilt{}

	return new(owner.PtrBuilt), nil
}

// NothingInteresting has no protected type among its results, so the
// directive is reported instead of taking effect.
//
//gofactory:factory // want `//gofactory:factory must have a protected type among NothingInteresting's results`
func NothingInteresting() int {
	return 0
}

// OtherBypass carries no directive: a declared factory's permission to
// bypass belongs to the function the directive is on, not to its whole
// package.
func OtherBypass() {
	_ = owner.Solo{} // want `Use factory for owner.Solo \(other.Build\)`
}

// inaccessibleRepo is unexported, so its exported method below is still
// owner.Inaccessible's declared factory (its fact reaches main.go, which
// imports this package directly), but nothing outside this package can
// name inaccessibleRepo to call it.
type inaccessibleRepo struct{}

//gofactory:factory
func (inaccessibleRepo) Restore() owner.Inaccessible { // want Restore:"gofactory:factory"
	return owner.Inaccessible{}
}
