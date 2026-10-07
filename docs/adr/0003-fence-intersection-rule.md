# Fences: intersection rule over gitignore-like globs

Before this ticket (#50), `-packageGlobs` had one rule: a package matching *any* given glob was
exempt from every check, for every type, anywhere. With one glob this reads as "this subtree may
bypass any factory," which is what it was for. With two or more globs, it reads as "being inside
*some* bounded context exempts you from *every* bounded context's protection" — the globs disable
each other, and the more fences a team adds, the less any of them actually protects. The globs
were also plain string globs, not gitignore-like: an exact path like `pkg` never matched (`pkg/`
matched only `pkg`, and `pkg/**` matched `pkg` and its subpackages), and `*` crossed `/`.

We decided each `-packageGlobs` pattern is its own **fence**, and a type whose package lies in
one or more fences may be bypassed only by code inside *all* of those fences. A type in no fence
is untouched by fences entirely; module scope alone decides whether it's protected, except under
`-packageGlobsOnly`, where it is not protected at all. Nesting two fences (`a/**` and
`a/domain/**`) then guards the inner one more tightly without extra syntax: `a/domain`'s type lies
in both, so only code also under `a/domain/**` may bypass it, while `a/infra` (under `a/**` only)
cannot. We also moved to gitignore-like compilation — `/` as the separator, matched against both
the path and the path plus `/` — so an exact path matches on its own and `*` no longer crosses a
package boundary.

This is a breaking change for anyone passing one or more `-packageGlobs` patterns today, not only
users of several patterns: a fence no longer exempts its code from checks on types outside it, so
even a single pattern that used to mean "this subtree may bypass any factory" now means "this
subtree may bypass only its own fence's factories." The migration guide in the README covers
putting the bypassing and the bypassed packages in one fence. It ships as part of v1.1.0 alongside
the other breaking changes in #43, not as a separate major version (a separate decision, #43's
"breaking behaviour shipped as v1.1.0, not v2" candidate).

## Considered options

- **Keep "any glob, exempt from everything."** Rejected: this is the bug the ticket exists to fix
  — fences actively undermine each other instead of composing.
- **Inside *any* of the type's fences, not all of them.** Rejected: this only ever consults the
  type's own fences, so an unrelated fence can't leak protection from a type it doesn't name — but
  nesting stops tightening. With `a/**` and `a/domain/**`, `a/infra` sits inside one of
  `a/domain`'s fences (`a/**`), so it could bypass `a/domain`'s factory freely, the opposite of
  what nested fences are for.
- **Intersection rule, gitignore-like globs.** Chosen, for the reasons above: fences compose
  without weakening each other, nesting falls out of the rule for free, and the glob syntax is
  gitignore-like — though a path-plus-`/` match also makes `a/*` and `a/**` match `a` itself,
  unlike a real `.gitignore` pattern (see [Glob syntax](../../README.md#glob-syntax)).
