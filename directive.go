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

// With -crossPackageDirectives=false, applyIgnore and applyFactory export no
// fact, so ignored and factories are the only way isIgnored (as
// detector.locallyIgnored) and declaredFactoryIndex see this package's own
// //gofactory:ignore and //gofactory:factory.
type directiveState struct {
	ignored   map[types.Object]bool
	factories []*types.Func
	trust     *trustInfo
}

func checkDirectives(pass *analysis.Pass) *directiveState {
	state := &directiveState{
		ignored: make(map[types.Object]bool),
		trust:   newTrustInfo(),
	}

	for _, file := range pass.Files {
		checkFileDirectives(pass, file, state)
	}

	return state
}

func checkFileDirectives(
	pass *analysis.Pass, file *ast.File, state *directiveState,
) {
	consumed := make(map[*ast.CommentGroup]bool)

	processDoc(pass, consumed, file.Doc, declPackage, nil, state)

	for _, decl := range file.Decls {
		checkDeclDirectives(pass, decl, consumed, state)
	}

	for _, cg := range file.Comments {
		if consumed[cg] {
			continue
		}

		processDoc(pass, consumed, cg, declOther, nil, state)
	}
}

func checkDeclDirectives(
	pass *analysis.Pass,
	decl ast.Decl,
	consumed map[*ast.CommentGroup]bool,
	state *directiveState,
) {
	switch d := decl.(type) {
	case *ast.GenDecl:
		checkGenDeclDirectives(pass, d, consumed, state)
	case *ast.FuncDecl:
		processDoc(pass, consumed, d.Doc, declFunc, d, state)
	}
}

func checkGenDeclDirectives(
	pass *analysis.Pass,
	decl *ast.GenDecl,
	consumed map[*ast.CommentGroup]bool,
	state *directiveState,
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
			processDoc(pass, consumed, decl.Doc, kind, typeSpec, state)
		}

		processDoc(pass, consumed, typeSpec.Doc, kind, typeSpec, state)
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
	state *directiveState,
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

		applyDirective(pass, comment, name, kind, node, state)
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
	state *directiveState,
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
			applyIgnore(pass, typeSpec, state.ignored)
		}
	case directiveFactory:
		if funcDecl, ok := node.(*ast.FuncDecl); ok {
			applyFactory(pass, comment, funcDecl, state)
		}
	case directiveTrusted:
		applyTrusted(pass, state.trust, kind, node)
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
	exportFact(pass, obj, &ignoredFact{})
}

func applyFactory(
	pass *analysis.Pass, comment *ast.Comment, funcDecl *ast.FuncDecl,
	state *directiveState,
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

	state.factories = append(state.factories, factory)
	exportFact(pass, factory, &factoryFact{})
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
