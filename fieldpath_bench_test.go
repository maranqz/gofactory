package gofactory_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/maranqz/gofactory"
	"golang.org/x/tools/go/analysis/analysistest"
)

// BenchmarkFieldPaths runs the analyzer with -zeroValues over a generated
// package: package domain holds a chain of depth nested by-value struct
// types, and package root declares width distinct types that each embed
// the whole chain, so it is shared rather than duplicated. Without the
// per-type cache, computing the chain's field paths costs
// O(width*depth); with it, package domain's chain is expanded once and
// reused width times, O(width+depth).
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

	b.ResetTimer()

	for range b.N {
		analysistest.Run(b, dir, analyzer, root)
	}
}

// generateFieldPathPackage writes a self-contained module to a temporary
// directory and returns its root. Package domain declares Level0..LevelD,
// each wrapping the previous by value; package root declares Root0..RootW,
// each with a Chain field of type domain.LevelD, and a package-level var
// of each, so every field down the chain is a reported candidate
// (checkPackageVars): LevelD, LevelD-1, ..., Level0, one path segment
// "Next" deeper each time.
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

// rootSource writes width Root types and, for each, a package-level var
// whose want comment lists the depth+1 diagnostics fieldPaths(RootI)
// produces: RootI itself is silent (root is its own owner package), and
// each domain.LevelK from LevelD down to Level0 is reported once, with
// "Chain" followed by (depth-K) "Next" segments as its path.
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
