package gofactory

import (
	"fmt"
	"go/types"

	"github.com/gobwas/glob"
)

type blockedStrategy interface {
	IsBlocked(currentPkg *types.Package, identObj types.Object) bool
}

type nilPkg struct{}

func newNilPkg() nilPkg {
	return nilPkg{}
}

func (nilPkg) IsBlocked(_ *types.Package, _ types.Object) bool {
	return false
}

type anotherPkg struct{}

func newAnotherPkg() anotherPkg {
	return anotherPkg{}
}

func (anotherPkg) IsBlocked(
	currentPkg *types.Package,
	identObj types.Object,
) bool {
	return currentPkg.Path() != identObj.Pkg().Path()
}

// fence is one -packageGlobs pattern: the set of packages matching its
// glob. It is compiled with '/' as the separator, so a '*' does not cross
// a package boundary, and is tested against both the package path and the
// path plus "/", so a pattern with no wildcard matches that exact path.
type fence struct {
	glob glob.Glob
}

func newFence(pattern string) (fence, error) {
	compiled, err := glob.Compile(pattern, '/')
	if err != nil {
		return fence{}, fmt.Errorf("unable to compile packageGlobs pattern %q: %w", pattern, err)
	}

	return fence{glob: compiled}, nil
}

// newFences compiles every -packageGlobs pattern into a fence, in order,
// stopping at the first invalid pattern.
func newFences(patterns []string) ([]fence, error) {
	fences := make([]fence, 0, len(patterns))

	for _, pattern := range patterns {
		f, err := newFence(pattern)
		if err != nil {
			return nil, err
		}

		fences = append(fences, f)
	}

	return fences, nil
}

func (f fence) contains(pkgPath string) bool {
	return f.glob.Match(pkgPath) || f.glob.Match(pkgPath+"/")
}

// fencedPkgs applies the intersection rule: a type whose package lies in
// one or more fences may be bypassed only by code inside every one of
// those fences, so several fences no longer disable each other. A type in
// no fence gains nothing from fences and falls back to defaultStrategy.
type fencedPkgs struct {
	fences          []fence
	defaultStrategy blockedStrategy
}

func newFencedPkgs(
	fences []fence,
	defaultStrategy blockedStrategy,
) fencedPkgs {
	return fencedPkgs{
		fences:          fences,
		defaultStrategy: defaultStrategy,
	}
}

func (s fencedPkgs) IsBlocked(
	currentPkg *types.Package,
	identObj types.Object,
) bool {
	identPkgPath := identObj.Pkg().Path()

	inAnyFence := false

	for _, fence := range s.fences {
		if !fence.contains(identPkgPath) {
			continue
		}

		inAnyFence = true

		if !fence.contains(currentPkg.Path()) {
			return true
		}
	}

	if !inAnyFence {
		return s.defaultStrategy.IsBlocked(currentPkg, identObj)
	}

	return false
}
