package gofactory

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// directivePrefix marks a directive comment. Like //go:build or //nolint:,
// it takes effect only with no space between "//" and the prefix; "// gofactory:ignore"
// is prose, not a directive.
const directivePrefix = "//gofactory:"

// directiveContext names the kind of declaration a directive comment was
// found on, which decides whether that directive is valid there.
type directiveContext int

const (
	contextOther directiveContext = iota
	contextType
	contextFunc
	contextPackage
)

// Known directive names. factory and trusted are recognised and validated
// here, the machinery the issue says this ticket builds for them to reuse,
// but only ignore has an effect today; the others are future tickets.
const (
	directiveIgnore  = "ignore"
	directiveFactory = "factory"
	directiveTrusted = "trusted"
)

// directiveAllowedIn reports whether name is a known directive valid in
// ctx. The second result is false for an unknown name, regardless of ctx.
func directiveAllowedIn(name string, ctx directiveContext) (bool, bool) {
	switch name {
	case directiveIgnore:
		return ctx == contextType, true
	case directiveFactory:
		return ctx == contextFunc, true
	case directiveTrusted:
		return ctx == contextFunc || ctx == contextPackage, true
	default:
		return false, false
	}
}

// placementDesc describes where name belongs, for the misplaced-directive
// diagnostic.
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

// checkDirectives parses, validates and applies every //gofactory: comment
// in the package. Unknown directives and known directives in the wrong
// place are reported as diagnostics at the comment. A valid
// //gofactory:ignore exports the ignoredFact on its type, which is what
// propagates the directive to importing packages and modules: FactTypes
// makes drivers analyse dependencies, and ImportObjectFact later sees the
// fact whether the type was declared in this package or an imported one.
func checkDirectives(pass *analysis.Pass) {
	for _, file := range pass.Files {
		checkFileDirectives(pass, file)
	}
}

func checkFileDirectives(pass *analysis.Pass, file *ast.File) {
	consumed := make(map[*ast.CommentGroup]bool)

	processDoc(pass, file.Doc, contextPackage, nil)
	consumed[file.Doc] = true

	for _, decl := range file.Decls {
		checkDeclDirectives(pass, decl, consumed)
	}

	for _, cg := range file.Comments {
		if consumed[cg] {
			continue
		}

		processDoc(pass, cg, contextOther, nil)
	}
}

// checkDeclDirectives processes the doc comments a declaration can carry,
// and records each one in consumed so the file-wide sweep over every
// comment group does not process it a second time as a free-floating,
// always-misplaced comment.
func checkDeclDirectives(
	pass *analysis.Pass,
	decl ast.Decl,
	consumed map[*ast.CommentGroup]bool,
) {
	switch d := decl.(type) {
	case *ast.GenDecl:
		checkGenDeclDirectives(pass, d, consumed)
	case *ast.FuncDecl:
		processDoc(pass, d.Doc, contextFunc, nil)
		consumed[d.Doc] = true
	}
}

func checkGenDeclDirectives(
	pass *analysis.Pass,
	decl *ast.GenDecl,
	consumed map[*ast.CommentGroup]bool,
) {
	if decl.Tok != token.TYPE {
		processDoc(pass, decl.Doc, contextOther, nil)
		consumed[decl.Doc] = true

		for _, spec := range decl.Specs {
			doc := specDoc(spec)
			processDoc(pass, doc, contextOther, nil)
			consumed[doc] = true
		}

		return
	}

	// A lone "type T struct{}" carries its doc on the GenDecl; a grouped
	// "type ( T struct{} )" carries it on each TypeSpec instead.
	if len(decl.Specs) == 1 {
		processDoc(pass, decl.Doc, contextType, decl.Specs[0])
		consumed[decl.Doc] = true
	}

	for _, spec := range decl.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}

		processDoc(pass, typeSpec.Doc, contextType, typeSpec)
		consumed[typeSpec.Doc] = true
	}
}

func specDoc(spec ast.Spec) *ast.CommentGroup {
	switch spec := spec.(type) {
	case *ast.ValueSpec:
		return spec.Doc
	case *ast.ImportSpec:
		return spec.Doc
	case *ast.TypeSpec:
		return spec.Doc
	default:
		return nil
	}
}

func processDoc(
	pass *analysis.Pass,
	doc *ast.CommentGroup,
	ctx directiveContext,
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

		applyDirective(pass, comment, name, ctx, spec)
	}
}

// parseDirective extracts the directive name from a single-line comment's
// raw text, e.g. "//gofactory:ignore" -> "ignore", stopping at the first
// space so a trailing "// want ..." testdata annotation on the same line
// is not taken for part of the name. It does not match multi-line /* */
// comments or a comment with a space after "//".
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

// applyDirective reports name as unknown or misplaced, or, once it is
// confirmed valid, applies its effect: today that is only
// //gofactory:ignore exporting ignoredFact on its type.
func applyDirective(
	pass *analysis.Pass,
	comment *ast.Comment,
	name string,
	ctx directiveContext,
	spec ast.Spec,
) {
	allowed, known := directiveAllowedIn(name, ctx)
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
