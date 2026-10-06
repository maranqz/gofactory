package gofactory_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/maranqz/gofactory"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
	"golang.org/x/tools/go/analysis/checker"
	"golang.org/x/tools/go/packages"
)

// BenchmarkFieldPaths runs the analyzer with -zeroValues over a generated
// package: package domain holds a chain of depth nested by-value struct
// types, and package root declares width distinct types that each hold
// the whole chain in a Chain field, so it is shared rather than
// duplicated. With the per-type cache, package domain's chain is expanded
// once and reused by every root, instead of once per root.
func BenchmarkFieldPaths(b *testing.B) {
	const (
		depth = 20
		width = 30
	)

	dir := generateFieldPathPackage(b, depth, width)

	analyzer := gofactory.NewAnalyzer()

	err := analyzer.Flags.Set("zeroValues", "true")
	if err != nil {
		b.Fatal(err)
	}

	root := filepath.Join(dir, "root")

	// analysistest.Run is both the correctness check against the want
	// comments below, and the one load of the generated package the timed
	// loop below analyzes.
	results := analysistest.Run(b, dir, analyzer, root)
	pkgs := []*packages.Package{results[0].Action.Package}

	b.ResetTimer()

	for range b.N {
		_, err := checker.Analyze([]*analysis.Analyzer{analyzer}, pkgs, nil)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func generateFieldPathPackage(b *testing.B, depth, width int) string {
	b.Helper()

	dir := b.TempDir()

	mustWriteFile(b, filepath.Join(dir, "go.mod"), "module bench\n\ngo 1.26\n")
	mustWriteFile(b, filepath.Join(dir, "domain", "domain.go"), domainSource(depth))
	mustWriteFile(b, filepath.Join(dir, "root", "root.go"), rootSource(depth, width))

	return dir
}

func domainSource(depth int) string {
	var src strings.Builder

	src.WriteString("package domain\n\n")
	src.WriteString("type Level0 struct{ Field int }\n\n")

	for i := 1; i <= depth; i++ {
		fmt.Fprintf(&src, "type Level%d struct{ Next Level%d }\n\n", i, i-1)
	}

	return src.String()
}

// RootI's own zero value is silent: root is its own owner package.
func rootSource(depth, width int) string {
	var src strings.Builder

	src.WriteString("package root\n\n")
	src.WriteString(`import "bench/domain"` + "\n\n")

	for i := range width {
		fmt.Fprintf(&src, "type Root%d struct{ Chain domain.Level%d }\n\n", i, depth)
	}

	src.WriteString("var (\n")

	for i := range width {
		fmt.Fprintf(&src, "\troot%d Root%d // %s\n", i, i, chainWantComment(depth))
	}

	src.WriteString(")\n")

	return src.String()
}

func chainWantComment(depth int) string {
	patterns := make([]string, 0, depth+1)

	path := "Chain"
	for level := depth; level >= 0; level-- {
		patterns = append(patterns, fmt.Sprintf(
			"`Use factory for domain.Level%d: zero value in %s`",
			level, path,
		))

		path += ".Next"
	}

	return "want " + strings.Join(patterns, " ")
}

func mustWriteFile(b *testing.B, path, content string) {
	b.Helper()

	err := os.MkdirAll(filepath.Dir(path), 0o750)
	if err != nil {
		b.Fatal(err)
	}

	err = os.WriteFile(path, []byte(content), 0o600)
	if err != nil {
		b.Fatal(err)
	}
}
