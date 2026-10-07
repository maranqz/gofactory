package gofactory

import (
	"go/ast"
	"go/token"
	"go/types"
	"slices"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Like //go:build or //nolint:, a directive needs no space between "//" and
// the prefix: "// gofactory:ignore" is prose, neither applied nor reported.
const directivePrefix = "//gofactory:"

const misplacedDirectiveFormat = "%s%s must be in the doc comment of %s"

type declKind int

const (
	declOther declKind = iota
	declType
	declAlias
	declFunc
	declPackage
)

// factory and trusted are only placement-checked.
const (
	directiveIgnore  = "ignore"
	directiveFactory = "factory"
	directiveTrusted = "trusted"
)

type placement struct {
	kinds []declKind
	desc  string
}

func placementOf(name string) (placement, bool) {
	switch name {
	case directiveIgnore:
		return placement{
			kinds: []declKind{declType},
			desc:  "a single top-level type definition",
		}, true
	case directiveFactory:
		return placement{
			kinds: []declKind{declFunc},
			desc:  "a function or method",
		}, true
	case directiveTrusted:
		return placement{
			kinds: []declKind{declFunc, declPackage},
			desc:  "a function, a method, or a package",
		}, true
	default:
		return placement{}, false
	}
}

// checkDirectives applies every //gofactory: directive found in pass and
// returns the types.Object of each one ignored by a //gofactory:ignore in
// this package. isIgnored consults that set directly: with
// -crossPackageDirectives=false, Analyzer.FactTypes is empty and applyIgnore
// stops exporting ignoredFact, so a same-package ignore would otherwise be
// unreachable to this pass.
func checkDirectives(pass *analysis.Pass) map[types.Object]bool {
	ignored := make(map[types.Object]bool)

	for _, file := range pass.Files {
		checkFileDirectives(pass, file, ignored)
	}

	return ignored
}

func checkFileDirectives(
	pass *analysis.Pass,
	file *ast.File,
	ignored map[types.Object]bool,
) {
	consumed := make(map[*ast.CommentGroup]bool)

	processDoc(pass, consumed, ignored, file.Doc, declPackage, nil)

	for _, decl := range file.Decls {
		checkDeclDirectives(pass, decl, consumed, ignored)
	}

	for _, cg := range file.Comments {
		if consumed[cg] {
			continue
		}

		processDoc(pass, consumed, ignored, cg, declOther, nil)
	}
}

func checkDeclDirectives(
	pass *analysis.Pass,
	decl ast.Decl,
	consumed map[*ast.CommentGroup]bool,
	ignored map[types.Object]bool,
) {
	switch d := decl.(type) {
	case *ast.GenDecl:
		checkGenDeclDirectives(pass, d, consumed, ignored)
	case *ast.FuncDecl:
		processDoc(pass, consumed, ignored, d.Doc, declFunc, nil)
	}
}

func checkGenDeclDirectives(
	pass *analysis.Pass,
	decl *ast.GenDecl,
	consumed map[*ast.CommentGroup]bool,
	ignored map[types.Object]bool,
) {
	if decl.Tok != token.TYPE {
		return // the file-wide sweep reports these as misplaced
	}

	// The parser leaves a doc above "type" on the GenDecl, and a doc above a
	// spec inside "type ( ... )" on that TypeSpec; the GenDecl's doc can
	// only belong to a type when the group declares exactly one.
	for _, spec := range decl.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}

		kind := declType
		if typeSpec.Assign.IsValid() {
			kind = declAlias
		}

		if len(decl.Specs) == 1 {
			processDoc(pass, consumed, ignored, decl.Doc, kind, typeSpec)
		}

		processDoc(pass, consumed, ignored, typeSpec.Doc, kind, typeSpec)
	}
}

func processDoc(
	pass *analysis.Pass,
	consumed map[*ast.CommentGroup]bool,
	ignored map[types.Object]bool,
	doc *ast.CommentGroup,
	kind declKind,
	typeSpec *ast.TypeSpec,
) {
	if doc == nil {
		return
	}

	consumed[doc] = true

	for _, comment := range doc.List {
		name, ok := parseDirective(comment.Text)
		if !ok {
			continue
		}

		applyDirective(pass, comment, name, kind, typeSpec, ignored)
	}
}

// Text after the name is ignored, so a trailing comment such as
// analysistest's "// want" can share the line.
func parseDirective(text string) (string, bool) {
	if !strings.HasPrefix(text, directivePrefix) {
		return "", false
	}

	name := text[len(directivePrefix):]
	if i := strings.IndexAny(name, " \t"); i >= 0 {
		name = name[:i]
	}

	return name, true
}

func applyDirective(
	pass *analysis.Pass,
	comment *ast.Comment,
	name string,
	kind declKind,
	typeSpec *ast.TypeSpec,
	ignored map[types.Object]bool,
) {
	place, known := placementOf(name)
	if !known {
		pass.Reportf(comment.Pos(), "unknown directive %q", directivePrefix+name)

		return
	}

	// The fact would land on the alias, but the detector looks types up
	// unaliased, and ExportObjectFact cannot reach the aliased type when it
	// lives in another package.
	if kind == declAlias && slices.Contains(place.kinds, declType) {
		pass.Reportf(
			comment.Pos(),
			misplacedDirectiveFormat+", not an alias",
			directivePrefix, name, place.desc,
		)

		return
	}

	if !slices.Contains(place.kinds, kind) {
		pass.Reportf(
			comment.Pos(),
			misplacedDirectiveFormat,
			directivePrefix, name, place.desc,
		)

		return
	}

	if name == directiveIgnore {
		applyIgnore(pass, typeSpec, ignored)
	}
}

func applyIgnore(
	pass *analysis.Pass,
	typeSpec *ast.TypeSpec,
	ignored map[types.Object]bool,
) {
	obj := pass.TypesInfo.ObjectOf(typeSpec.Name)
	if obj == nil {
		return
	}

	ignored[obj] = true

	// An Analyzer that uses facts must declare their types (go/analysis's
	// documented rule); with -crossPackageDirectives=false, Analyzer.FactTypes
	// is empty, and go vet's gob encoder panics on an undeclared fact.
	if len(pass.Analyzer.FactTypes) > 0 {
		pass.ExportObjectFact(obj, &ignoredFact{})
	}
}
