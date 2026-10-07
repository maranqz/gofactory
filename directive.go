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

const noProtectedResultFormat = "%s%s must have a protected type among %s's results"

type declKind int

const (
	declOther declKind = iota
	declType
	declAlias
	declFunc
	declPackage
)

// trusted is only placement-checked.
const (
	directiveIgnore  = "ignore"
	directiveFactory = "factory"
	directiveTrusted = "trusted"
)

// apply is nil for a directive with no effect yet (trusted): applyDirective
// then stops at the placement check.
type placement struct {
	kinds []declKind
	desc  string
	apply func(pass *analysis.Pass, comment *ast.Comment, node ast.Node)
}

func placementOf(name string) (placement, bool) {
	switch name {
	case directiveIgnore:
		return placement{
			kinds: []declKind{declType},
			desc:  "a single top-level type definition",
			apply: func(pass *analysis.Pass, _ *ast.Comment, node ast.Node) {
				if typeSpec, ok := node.(*ast.TypeSpec); ok {
					applyIgnore(pass, typeSpec)
				}
			},
		}, true
	case directiveFactory:
		return placement{
			kinds: []declKind{declFunc},
			desc:  "a function or method",
			apply: func(pass *analysis.Pass, comment *ast.Comment, node ast.Node) {
				if funcDecl, ok := node.(*ast.FuncDecl); ok {
					applyFactory(pass, comment, funcDecl)
				}
			},
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

func checkDirectives(pass *analysis.Pass) {
	for _, file := range pass.Files {
		checkFileDirectives(pass, file)
	}
}

func checkFileDirectives(pass *analysis.Pass, file *ast.File) {
	consumed := make(map[*ast.CommentGroup]bool)

	processDoc(pass, consumed, file.Doc, declPackage, nil)

	for _, decl := range file.Decls {
		checkDeclDirectives(pass, decl, consumed)
	}

	for _, cg := range file.Comments {
		if consumed[cg] {
			continue
		}

		processDoc(pass, consumed, cg, declOther, nil)
	}
}

func checkDeclDirectives(
	pass *analysis.Pass,
	decl ast.Decl,
	consumed map[*ast.CommentGroup]bool,
) {
	switch d := decl.(type) {
	case *ast.GenDecl:
		checkGenDeclDirectives(pass, d, consumed)
	case *ast.FuncDecl:
		processDoc(pass, consumed, d.Doc, declFunc, d)
	}
}

func checkGenDeclDirectives(
	pass *analysis.Pass,
	decl *ast.GenDecl,
	consumed map[*ast.CommentGroup]bool,
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
			processDoc(pass, consumed, decl.Doc, kind, typeSpec)
		}

		processDoc(pass, consumed, typeSpec.Doc, kind, typeSpec)
	}
}

// node is the declaration the doc comment belongs to: a *ast.TypeSpec for
// declType/declAlias, a *ast.FuncDecl for declFunc, nil otherwise. It is
// only consulted by the directives whose placement allows it.
func processDoc(
	pass *analysis.Pass,
	consumed map[*ast.CommentGroup]bool,
	doc *ast.CommentGroup,
	kind declKind,
	node ast.Node,
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

		applyDirective(pass, comment, name, kind, node)
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

	if place.apply != nil {
		place.apply(pass, comment, node)
	}
}

func applyIgnore(pass *analysis.Pass, typeSpec *ast.TypeSpec) {
	obj := pass.TypesInfo.ObjectOf(typeSpec.Name)
	if obj == nil {
		return
	}

	pass.ExportObjectFact(obj, &ignoredFact{})
}

func applyFactory(
	pass *analysis.Pass, comment *ast.Comment, funcDecl *ast.FuncDecl,
) {
	factory, ok := pass.TypesInfo.ObjectOf(funcDecl.Name).(*types.Func)
	if !ok {
		return
	}

	if len(protectedResultTargets(factory.Signature())) == 0 {
		pass.Reportf(
			comment.Pos(), noProtectedResultFormat,
			directivePrefix, directiveFactory, factory.Name(),
		)

		return
	}

	pass.ExportObjectFact(factory, &factoryFact{})
}
