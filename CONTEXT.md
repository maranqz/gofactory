# gofactory

A Go linter that makes code obtain values of protected types through their factories instead of bypassing them, so that the invariants a factory establishes cannot be skipped.

## Language

Terms and their _Avoid_ lists describe the domain the analyzer reasons about — protected types and their factories in the code under lint. They don't govern prose about gofactory's own implementation (e.g. its golangci-lint plugin registration function): ordinary Go vocabulary like "constructor" is fine there.

### Types and factories

**Protected type**:
A named type, other than a func type or an interface, whose values outside its owner package must come from a factory.

_Avoid_: blocked type, structure, struct

**Ignored type**:
A type taken out of protection by a directive or a setting; anyone may create it without a factory.

_Avoid_: excluded type, allowed type

**Owner package**:
The package that declares a type. Its external test package (`foo_test`) counts as the same owner.

_Avoid_: local package, home package

**Factory**:
A function or method that is the sanctioned way to obtain a `T` or `*T`. It is either recognised in the owner package of `T` by its exported signature and name, or declared a factory anywhere by a directive or a setting. Factories for creation and for reconstitution are the same concept.

_Avoid_: constructor, builder

### Bypassing a factory

**Factory bypass**:
Obtaining a value of a protected type other than by calling one of its factories.

_Avoid_: direct creation, creation (DDD reserves it for an object's birth), instantiation (that is generic instantiation in Go), initialization, construction

**Bypass route**:
One of the ways a factory can be bypassed: a literal, a conversion, `new`, an implicit constant conversion, or a zero value.

_Avoid_: creation route, creation method (a "Creation Method" is a factory), creation kind (kinds belong to types), vector

**Zero-value bypass**:
A factory bypass by declaring a variable, field or element of a protected type without giving it a value.

_Avoid_: zero-value creation

**First interaction**:
The first mention of a zero-valued variable in its function; it decides whether that zero value counts as a factory bypass.

**Reconstitution**:
Rebuilding an existing object from storage, as a DDD repository does. For gofactory it is a factory bypass unless it goes through a factory or happens in trusted code.

_Avoid_: hydration, rehydration

### Where a factory may be bypassed

**Fence**:
A set of packages given by one package glob. If a type's package lies in one or more fences, only code inside all of those fences may bypass that type's factory.

_Avoid_: blocked package, glob package

**Producer**:
A top-level function or method in the owner package of `T` whose results include `T`; when the owner package itself is checked, it is the only place there where the factory of `T` may be bypassed. A producer is not necessarily a factory.

_Avoid_: factory (for producers that are not factories), constructor

**Trusted code**:
A package, function or method that may bypass the factory of any protected type.

_Avoid_: whitelist, allowed package

**Directive**:
A `//gofactory:` comment that marks a type, function, method or package as ignored, a factory, or trusted.

_Avoid_: annotation, pragma, marker
