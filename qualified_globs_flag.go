package gofactory

import (
	"fmt"
	"strings"

	"github.com/gobwas/glob"
)

// qualifiedGlobsFlag holds globs matched against qualified names
// (import/path.Name), gitignore-like: compiled with '/' as the only
// separator, so '*' crosses '.' but not '/'.
type qualifiedGlobsFlag struct {
	raw   []string
	globs []glob.Glob
}

func (g *qualifiedGlobsFlag) String() string {
	return strings.Join(g.raw, ", ")
}

func (g *qualifiedGlobsFlag) Set(pattern string) error {
	pattern = strings.TrimSpace(pattern)

	compiled, err := glob.Compile(pattern, '/')
	if err != nil {
		return fmt.Errorf("unable to compile glob %s: %w", pattern, err)
	}

	g.raw = append(g.raw, pattern)
	g.globs = append(g.globs, compiled)

	return nil
}

func (g *qualifiedGlobsFlag) Value() []glob.Glob {
	return g.globs
}
