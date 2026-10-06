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
