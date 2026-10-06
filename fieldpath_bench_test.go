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
// duplicated. Without the per-type cache, computing the chain's field
// paths costs O(width*depth); with it, package domain's chain is expanded
// once and reused width times, O(width+depth).
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

	// analysistest.Run is the correctness check, against the want comments
	// below; it also loads the packages benchmarked by the loop, which
	// times only checker.Analyze so that package loading isn't counted.
	analysistest.Run(b, dir, analyzer, root)

	pkgs := loadFieldPathPackage(b, dir)

	b.ResetTimer()

	for range b.N {
		_, err := checker.Analyze([]*analysis.Analyzer{analyzer}, pkgs, nil)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// loadFieldPathPackage loads the root package of the module at dir the way
// analysistest.Run does for a module-mode root: GOPROXY=off keeps it from
// reaching the network, and GOWORK=off keeps it from picking up this
// repo's go.work.
func loadFieldPathPackage(b *testing.B, dir string) []*packages.Package {
	b.Helper()

	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles |
			packages.NeedImports | packages.NeedTypes | packages.NeedTypesSizes |
			packages.NeedSyntax | packages.NeedTypesInfo |
			packages.NeedDeps | packages.NeedModule,
		Dir: dir,
		Env: append(os.Environ(), "GO111MODULE=on", "GOPROXY=off", "GOWORK=off"),
	}

	pkgs, err := packages.Load(cfg, "./root")
	if err != nil {
		b.Fatal(err)
	}

	if packages.PrintErrors(pkgs) > 0 {
		b.Fatal("errors loading the generated package")
	}

	return pkgs
}

// generateFieldPathPackage writes a self-contained module to a temporary
// directory and returns its root.
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
