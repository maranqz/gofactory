# Cross-package directives via analysis facts, on by default

`//gofactory:ignore` lives next to the type it describes, in that type's owner package. For the
directive to take effect in every importer, as #53 requires, the owner package's analysis must hand the
decision to every package and module that imports the type, even where the user neither
configures `--ignoreTypes` nor repeats the directive.

We export it as an `analysis.Fact` (`FactTypes = []analysis.Fact{new(ignoredFact)}`) and read it
back with `ImportObjectFact`, rather than, say, re-deriving it from source comments at each import
site (not possible: an importing package only has the dependency's export data, not its comments)
or asking users to repeat `--ignoreTypes` per module.

The trade-off: declaring `FactTypes` makes `go vet` and golangci-lint analyse the current package's
full transitive dependency graph, not just its own files, which is slower on a large monorepo. We
chose correctness-by-default over that cost; `-crossPackageDirectives=false` to opt back out of the
transitive analysis is deferred to a later ticket (#43, user story 35).

A second cost: a configuration error, which `run` returns, now fails the analysis of every
dependency as well. The standalone `gofactory` CLI prints the error once for each dependency without
imports, among `failed prerequisites` lines for all the others, so one bad setting can produce
hundreds of lines; `go vet` still prints it once per package, and golangci-lint once per run.
Checking settings while flags are parsed can't avoid this: `-packageGlobsOnly` without
`-packageGlobs` can only be caught once every flag is set. Reading directives in a separate analyzer
that takes no settings can: only it would declare `FactTypes` and run on dependencies, while the
main analyzer, which checks the settings, would run only on the packages being linted. We leave
that split to its own ticket under #43.
