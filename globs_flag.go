package gofactory

import "strings"

// globsFlag collects raw -packageGlobs patterns. It does not compile them:
// a pattern is only a fence once newFences compiles it at Pass.Analyze time,
// which is also where an invalid pattern becomes a configuration error
// (see run in factory.go). Compiling eagerly here would catch typos sooner
// for the flags entry point, but the plugin entry point decodes every
// setting before an Analyzer exists to catch anything, so both would not
// surface an invalid pattern the same way.
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

// Append records pattern without the error return flag.Value.Set needs,
// for the plugin entry point, which has no use for an error that can't
// occur.
func (g *globsFlag) Append(pattern string) {
	g.patterns = append(g.patterns, strings.TrimSpace(pattern))
}

func (g *globsFlag) Value() []string {
	return g.patterns
}
