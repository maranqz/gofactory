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

// globalOuter: a package-level var is always a candidate, so its nested
// field path is reported the same as a local var's.
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
	return Wrapper{} // want `^Use factory for nested.Struct: zero value in S \(nested.NewStruct, nested.NewStructOrErr\)$`
}

// EmbeddedFieldLiteralIsReported: an embedded field's path segment is the
// embedded type's own name, Struct here.
func EmbeddedFieldLiteralIsReported() {
	_ = struct{ nested.Struct }{} // want `Use factory for nested.Struct: zero value in Struct`
}

type UnderlyingStruct struct {
	nested.Struct
}

// NamedEmbeddingLiteralIsReported is EmbeddedFieldLiteralIsReported's named
// form: the embedding type itself has a name rather than being anonymous.
func NamedEmbeddingLiteralIsReported() UnderlyingStruct {
	return UnderlyingStruct{} // want `Use factory for nested.Struct: zero value in Struct`
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
}

// FieldsBehindNotFollowedKindsAreSilent: a pointer, slice, map or chan
// field is not followed, so an empty literal of NotFollowed reports
// nothing.
func FieldsBehindNotFollowedKindsAreSilent() NotFollowed {
	return NotFollowed{}
}

type NamedContainers struct {
	G nested.Grid
	T nested.Tags
}

// NamedContainerFieldsAreReportedAtTheirOwnPath: a field whose own type is
// a protected named container is reported at its own path, same as any
// other field; the container's elements are not themselves followed.
func NamedContainerFieldsAreReportedAtTheirOwnPath() NamedContainers {
	return NamedContainers{} // want `Use factory for nested.Grid: zero value in G` `Use factory for nested.Tags: zero value in T`
}

// FieldWriteOnZeroValuedPairReportsBothFields: -zeroValues judges p's first
// interaction once for the whole variable (the previous ticket's rule), not
// each field path on its own, so writing A still reports both A and B at
// that write, even though the write already goes through A's factory.
func FieldWriteOnZeroValuedPairReportsBothFields() Pair {
	var p Pair

	p.A = nested.NewStruct(1) // want `Use factory for nested.Struct: zero value in A` `Use factory for nested.Struct: zero value in B`
	p.B = nested.NewStruct(2)

	return p
}

// FieldAddressPassedToCallOnZeroValuedPairReportsBothFields is
// FieldWriteOnZeroValuedPairReportsBothFields with &p.A passed to a call:
// only &p, the whole variable, is the silent form.
func FieldAddressPassedToCallOnZeroValuedPairReportsBothFields() Pair {
	var p Pair

	nested.Fill(&p.A) // want `Use factory for nested.Struct: zero value in A` `Use factory for nested.Struct: zero value in B`
	p.B = nested.NewStruct(2)

	return p
}

// WholeValueAssignOnZeroValuedVarIsSilent: w's first interaction is a
// whole-value assignment, the first-interaction rule's OK form, so no
// field-path report fires even though w started zero-valued.
func WholeValueAssignOnZeroValuedVarIsSilent() Wrapper {
	var w Wrapper

	w = Wrapper{S: nested.NewStruct(1)}

	return w
}

func discardWrapperPtr(*Wrapper) {}

// AddressOfWholeVarPassedToCallIsSilent: &w, the whole variable rather than
// one of its fields, is the first-interaction rule's other OK form.
func AddressOfWholeVarPassedToCallIsSilent() Wrapper {
	var w Wrapper

	discardWrapperPtr(&w)

	return w
}

type Blank struct {
	_ nested.Struct
	X int
}

// BlankFieldIsSilent: Go rejects _ as a literal key and as a selector, so
// a blank field of a protected type is never reported.
func BlankFieldIsSilent() int {
	var bl Blank

	_ = Blank{X: 1}

	return bl.X
}
