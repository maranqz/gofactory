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

// factory is only placement-checked so far.
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

// checkDirectives also collects //gofactory:trusted directives into the
// returned trustInfo: its effect is intra-pass (see trustInfo), so it has
// no fact to export and nothing to apply to a types.Object here.
func checkDirectives(pass *analysis.Pass) *trustInfo {
	trust := newTrustInfo()

	for _, file := range pass.Files {
		checkFileDirectives(pass, file, trust)
	}

	return trust
}

func checkFileDirectives(
	pass *analysis.Pass, file *ast.File, trust *trustInfo,
) {
	consumed := make(map[*ast.CommentGroup]bool)

	processDoc(pass, consumed, file.Doc, declPackage, nil, trust)

	for _, decl := range file.Decls {
		checkDeclDirectives(pass, decl, consumed, trust)
	}

	for _, cg := range file.Comments {
		if consumed[cg] {
			continue
		}

		processDoc(pass, consumed, cg, declOther, nil, trust)
	}
}

func checkDeclDirectives(
	pass *analysis.Pass,
	decl ast.Decl,
	consumed map[*ast.CommentGroup]bool,
	trust *trustInfo,
) {
	switch d := decl.(type) {
	case *ast.GenDecl:
		checkGenDeclDirectives(pass, d, consumed, trust)
	case *ast.FuncDecl:
		processDoc(pass, consumed, d.Doc, declFunc, d, trust)
	}
}

func checkGenDeclDirectives(
	pass *analysis.Pass,
	decl *ast.GenDecl,
	consumed map[*ast.CommentGroup]bool,
	trust *trustInfo,
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
			processDoc(pass, consumed, decl.Doc, kind, typeSpec, trust)
		}

		processDoc(pass, consumed, typeSpec.Doc, kind, typeSpec, trust)
	}
}

// node is the declaration the directive would apply to: *ast.TypeSpec for
// declType/declAlias, *ast.FuncDecl for declFunc, nil for declPackage and
// declOther, where no directive can take effect.
func processDoc(
	pass *analysis.Pass,
	consumed map[*ast.CommentGroup]bool,
	doc *ast.CommentGroup,
	kind declKind,
	node ast.Node,
	trust *trustInfo,
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

		applyDirective(pass, comment, name, kind, node, trust)
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
	node ast.Node,
	trust *trustInfo,
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

	switch name {
	case directiveIgnore:
		if typeSpec, ok := node.(*ast.TypeSpec); ok {
			applyIgnore(pass, typeSpec)
		}
	case directiveTrusted:
		applyTrusted(pass, trust, kind, node)
	}
}

func applyIgnore(pass *analysis.Pass, typeSpec *ast.TypeSpec) {
	obj := pass.TypesInfo.ObjectOf(typeSpec.Name)
	if obj == nil {
		return
	}

	pass.ExportObjectFact(obj, &ignoredFact{})
}

func applyTrusted(
	pass *analysis.Pass, trust *trustInfo, kind declKind, node ast.Node,
) {
	if kind == declPackage {
		trust.markPackage()

		return
	}

	decl, ok := node.(*ast.FuncDecl)
	if !ok {
		return
	}

	fn, ok := pass.TypesInfo.ObjectOf(decl.Name).(*types.Func)
	if !ok {
		return
	}

	trust.markFunc(fn)
}
