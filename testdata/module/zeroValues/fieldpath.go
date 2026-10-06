package zeroValues

import "factory/zeroValues/nested"

// Inner, Middle and Outer are plain types owned by this package: only
// nested.Struct behind them comes from another package, so a zero-valued
// Outer is reported exactly once, with the full path down to it.
type Inner struct {
	S nested.Struct
}

type Middle struct {
	W Inner
}

type Outer struct {
	V Middle
}

// globalOuter: a package-level var is always a candidate (checkPackageVars),
// so its nested field path is reported the same as a local var's.
var globalOuter Outer // want `Use factory for nested.Struct: zero value in V.W.S`

// EmptyOuterLiteralIsReportedWithFullPath: Outer{} leaves V unset, and so
// every field down to S, at any depth through plain by-value struct
// nesting.
func EmptyOuterLiteralIsReportedWithFullPath() Outer {
	return Outer{} // want `Use factory for nested.Struct: zero value in V.W.S`
}

// ZeroValuedVarReportsNestedFieldUnderFirstInteraction: v is never given a
// whole value, so its first interaction - the field read below - reports
// the same nested.Struct field path a zero-valued literal would.
func ZeroValuedVarReportsNestedFieldUnderFirstInteraction() nested.Struct {
	var v Outer

	return v.V.W.S // want `Use factory for nested.Struct: zero value in V.W.S`
}

type Wrapper struct {
	S nested.Struct
}

// EmptyWrapperLiteralIsReported: one level of nesting, the base case of
// the field-path rule.
func EmptyWrapperLiteralIsReported() Wrapper {
	return Wrapper{} // want `Use factory for nested.Struct: zero value in S`
}

// EmbeddedFieldLiteralIsReported: an embedded field's path segment is the
// embedded type's own name, Struct here.
func EmbeddedFieldLiteralIsReported() {
	_ = struct{ nested.Struct }{} // want `Use factory for nested.Struct: zero value in Struct`
}

type Pair struct {
	A nested.Struct
	B nested.Struct
}

// PartialKeyedLiteralReportsOnlyUnsetFields: A is keyed, so only B, left
// unset, is reported.
func PartialKeyedLiteralReportsOnlyUnsetFields() Pair {
	return Pair{A: nested.NewStruct(1)} // want `Use factory for nested.Struct: zero value in B`
}

// FullyPositionalLiteralReportsNone: Go requires a positional struct
// literal to supply every field, so nothing is left unset.
func FullyPositionalLiteralReportsNone() Pair {
	return Pair{nested.NewStruct(1), nested.NewStruct(2)}
}

type NotFollowed struct {
	Ptr   *nested.Struct
	Slice []nested.Struct
	Map   map[string]nested.Struct
	Chan  chan nested.Struct
	Array [2]nested.Struct
}

// FieldsBehindNotFollowedKindsAreSilent: a pointer, slice, map, chan or
// array field is not followed, so an empty literal of NotFollowed reports
// nothing.
func FieldsBehindNotFollowedKindsAreSilent() NotFollowed {
	return NotFollowed{}
}
