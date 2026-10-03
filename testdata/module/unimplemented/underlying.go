package unimplemented

// not implemented

/*
Idea.

Protection is keyed by type identity (owner package + name), not by
structural layout. DeclStruct is a distinct, locally-owned type, so the
"owner package may build its own types" rule applies to it even though its
underlying layout is copied from invariant.Struct, including invariant's
unexported "secret" field. Building DeclStruct{} never calls
invariant.NewStruct, so "secret" is silently left at its zero value instead
of the 42 the factory establishes — a real invariant bypass that the
protected-type model, as specified, cannot see.

Catching this would mean protecting by structural reuse of a protected
type's underlying layout, not by type identity alone, which is a different
detection model than anything in the current spec (module scope, fences,
ownPackage, trusted code are all about *where* code lives, not about
whether it structurally reuses a protected type's layout).
*/

import "factory/unimplemented/invariant"

type DeclStruct invariant.Struct

func BypassPrivateInvariant() DeclStruct {
	return DeclStruct{Field: 1} // want `Use factory for invariant.Struct`
}
