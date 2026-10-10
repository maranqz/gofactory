# Ship the go/types rewrite's breaking behaviour as v1.1.0, not v2

The go/types rewrite (#43) changes what gofactory reports: module scope by default, the fence
intersection rule and gitignore-like globs, a configuration error for `-packageGlobsOnly` without
globs, new bypass routes and a message suffix. Under semantic versioning, behaviour that breaks
existing users would call for v2.

We decided to ship it as v1.1.0. Go's semantic import versioning ties a major version to the module
path: v2 would mean `github.com/maranqz/gofactory/v2`, and every import, `go install` line,
`.custom-gcl.yml` and golangci-lint integration would have to change for a linter whose Go API is
just `NewAnalyzer()`, which stays as it is. The breaks are in what is reported, not in the API, and
a linter's reports change between minor versions as a matter of course, as golangci-lint's own
linters do. The cost lands on users who depend on today's reports, so the v1.1.0 changelog opens
with a "Breaking changes" section and links the glob migration guide, and the message prefix that
exclusion rules match stays stable ([ADR 0001](0001-stable-use-factory-for-prefix.md)).

## Considered options

- **v2 at `github.com/maranqz/gofactory/v2`.** Rejected: it breaks every import and install line
  to signal a change no Go code would notice, and it would leave golangci-lint PR #4196 targeting a
  path the project no longer develops.
- **Keep the old behaviour behind flags and change nothing by default.** Rejected: the old defaults
  are the stdlib noise, the globs that disable each other and the silent `-packageGlobsOnly` the
  rewrite set out to fix.
- **v1.1.0 with a "Breaking changes" changelog section.** Chosen, for the reasons above.
