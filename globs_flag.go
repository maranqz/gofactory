package gofactory

import "strings"

// globsFlag collects raw -packageGlobs and -ignoreTypes patterns; they are
// compiled, and an invalid one reported, in run (factory.go).
type globsFlag struct {
	patterns []string
}

func (g *globsFlag) String() string {
	return strings.Join(g.patterns, ", ")
}

func (g *globsFlag) Set(pattern string) error {
	g.Append(pattern)

	return nil
}

func (g *globsFlag) Append(pattern string) {
	g.patterns = append(g.patterns, strings.TrimSpace(pattern))
}

func (g *globsFlag) Value() []string {
	return g.patterns
}
