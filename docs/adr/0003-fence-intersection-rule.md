# Fences: intersection rule over gitignore-like globs

Before this ticket (#50), `-packageGlobs` had one rule: a package matching *any* given glob was
exempt from every check, for every type, anywhere. With one glob this reads as "this subtree may
bypass any factory," which is what it was for. With two or more globs, it reads as "being inside
*some* bounded context exempts you from *every* bounded context's protection" — the globs disable
each other, and the more fences a team adds, the less any of them actually protects. The globs
were also plain string globs, not gitignore-like: an exact path never matched (you always needed
a trailing `/**`), and `*` crossed `/`.

We decided each `-packageGlobs` pattern is its own **fence**, and a type whose package lies in
one or more fences may be bypassed only by code inside *all* of those fences. A type in no fence
is untouched by fences entirely; module scope alone decides whether it's protected. Nesting two
fences (`a/**` and `a/domain/**`) then guards the inner one more tightly without extra syntax:
`a/domain`'s type lies in both, so only code also under `a/domain/**` may bypass it, while
`a/infra` (under `a/**` only) cannot. We also moved to gitignore-like compilation — `/` as the
separator, matched against both the path and the path plus `/` — so an exact path matches on its
own and `*` no longer crosses a package boundary.

This is a breaking change for anyone passing one or more `-packageGlobs` patterns today, not only
users of several patterns: code inside a fence can no longer bypass the factory of a type outside
every fence that type's package lies in, so even a single pattern that used to mean "this subtree
may bypass any factory" now means "this subtree may bypass only its own fence's factories." The
migration guide in the README covers putting the bypassing and the bypassed packages in one fence.
It ships as part of v1.1.0 alongside the other breaking changes in #43, not as a separate major
version (a separate decision, #43's "breaking behaviour shipped as v1.1.0, not v2" candidate).

## Considered options

- **Keep "any glob, exempt from everything."** Rejected: this is the bug the ticket exists to fix
  — fences actively undermine each other instead of composing.
- **Union instead of intersection: a type is exempt if either its own fence or the bypassing
  code's fence allows it.** Rejected: this still lets an unrelated fence's membership leak
  protection from a type it doesn't name, just with one fewer step than the original bug.
- **Intersection rule, gitignore-like globs.** Chosen, for the reasons above: fences compose
  without weakening each other, nesting falls out of the rule for free, and the glob syntax
  matches what most users already expect from `.gitignore`.
