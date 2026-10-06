package gofactory

import (
	"errors"
	"fmt"
	"go/types"
	"strings"

	"github.com/gobwas/glob"
)

// errEmptyPackageGlobPattern is the configuration error for a -packageGlobs
// pattern that is empty after TrimSpace: it can never match a package, so
// it would silently protect nothing instead of forming a fence.
var errEmptyPackageGlobPattern = errors.New("packageGlobs pattern must not be empty")

// errLeadingSlashGlobPattern is the configuration error for a -packageGlobs
// pattern starting with '/': gitignore gives a leading '/' a special
// "from the root" meaning, but a Go package path never starts with '/', so
// such a pattern would compile and silently match nothing.
var errLeadingSlashGlobPattern = errors.New("packageGlobs pattern must not start with '/'")

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
// path plus "/": the bare-path match is what lets a pattern with no
// wildcard match that exact path, and the path-plus-"/" match is what lets
// `a/**` and `a/*` also match `a` itself.
type fence struct {
	glob glob.Glob
}

func newFence(pattern string) (fence, error) {
	if pattern == "" {
		return fence{}, errEmptyPackageGlobPattern
	}

	if strings.HasPrefix(pattern, "/") {
		return fence{}, errLeadingSlashGlobPattern
	}

	compiled, err := glob.Compile(pattern, '/')
	if err != nil {
		return fence{}, fmt.Errorf("unable to compile packageGlobs pattern %q: %w", pattern, err)
	}

	return fence{glob: compiled}, nil
}

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
// those fences. A type in no fence gains nothing from fences and falls
// back to defaultStrategy.
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

	for _, candidate := range s.fences {
		if !candidate.contains(identPkgPath) {
			continue
		}

		inAnyFence = true

		if !candidate.contains(currentPkg.Path()) {
			return true
		}
	}

	if !inAnyFence {
		return s.defaultStrategy.IsBlocked(currentPkg, identObj)
	}

	return false
}
