package gofactory

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Like //go:build or //nolint:, a directive needs no space between "//" and
// the prefix: "// gofactory:ignore" is prose, neither applied nor reported.
const directivePrefix = "//gofactory:"

type declKind int

const (
	declOther declKind = iota
	declType
	declFunc
	declPackage
)

// factory and trusted are only placement-checked.
const (
	directiveIgnore  = "ignore"
	directiveFactory = "factory"
	directiveTrusted = "trusted"
)

// directiveAllowedIn reports whether name is a known directive valid in
// kind. The second result is false for an unknown name, regardless of kind.
func directiveAllowedIn(name string, kind declKind) (bool, bool) {
	switch name {
	case directiveIgnore:
		return kind == declType, true
	case directiveFactory:
		return kind == declFunc, true
	case directiveTrusted:
		return kind == declFunc || kind == declPackage, true
	default:
		return false, false
	}
}

func placementDesc(name string) string {
	switch name {
	case directiveIgnore:
		return "a type declaration"
	case directiveFactory:
		return "a function or method declaration"
	case directiveTrusted:
		return "a function, a method, or a package declaration"
	default:
		return ""
	}
}

func checkDirectives(pass *analysis.Pass) {
	for _, file := range pass.Files {
		checkFileDirectives(pass, file)
	}
}

func checkFileDirectives(pass *analysis.Pass, file *ast.File) {
	consumed := make(map[*ast.CommentGroup]bool)

	processDoc(pass, file.Doc, declPackage, nil)
	consumed[file.Doc] = true

	for _, decl := range file.Decls {
		checkDeclDirectives(pass, decl, consumed)
	}

	for _, cg := range file.Comments {
		if consumed[cg] {
			continue
		}

		processDoc(pass, cg, declOther, nil)
	}
}

// Mark every doc processed here in consumed, or the file sweep reports it
// again as misplaced.
func checkDeclDirectives(
	pass *analysis.Pass,
	decl ast.Decl,
	consumed map[*ast.CommentGroup]bool,
) {
	switch d := decl.(type) {
	case *ast.GenDecl:
		checkGenDeclDirectives(pass, d, consumed)
	case *ast.FuncDecl:
		processDoc(pass, d.Doc, declFunc, nil)
		consumed[d.Doc] = true
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
	if len(decl.Specs) == 1 {
		processDoc(pass, decl.Doc, declType, decl.Specs[0])
		consumed[decl.Doc] = true
	}

	for _, spec := range decl.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}

		processDoc(pass, typeSpec.Doc, declType, typeSpec)
		consumed[typeSpec.Doc] = true
	}
}

func processDoc(
	pass *analysis.Pass,
	doc *ast.CommentGroup,
	kind declKind,
	spec ast.Spec,
) {
	if doc == nil {
		return
	}

	for _, comment := range doc.List {
		name, ok := parseDirective(comment.Text)
		if !ok {
			continue
		}

		applyDirective(pass, comment, name, kind, spec)
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
	spec ast.Spec,
) {
	allowed, known := directiveAllowedIn(name, kind)
	if !known {
		pass.Reportf(comment.Pos(), "unknown directive %q", directivePrefix+name)

		return
	}

	if !allowed {
		pass.Reportf(
			comment.Pos(),
			"%s%s must be on %s",
			directivePrefix, name, placementDesc(name),
		)

		return
	}

	if name == directiveIgnore {
		applyIgnore(pass, spec)
	}
}

func applyIgnore(pass *analysis.Pass, spec ast.Spec) {
	typeSpec, ok := spec.(*ast.TypeSpec)
	if !ok {
		return
	}

	obj := pass.TypesInfo.ObjectOf(typeSpec.Name)
	if obj == nil {
		return
	}

	pass.ExportObjectFact(obj, &ignoredFact{})
}
