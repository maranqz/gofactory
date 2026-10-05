# Keep `Use factory for pkg.T` as a stable message prefix

#51 adds a suffix naming accessible factories to the diagnostic message. A cleaner message could
have dropped or reworded the existing `Use factory for pkg.T` prefix instead of appending to it, but
users match on it with golangci-lint's own text matchers, `linters.exclusions.rules[].text` and
`severity.rules[].text`, and a changed prefix would silently break those. We decided to keep the
prefix verbatim and append the suffix after it, so a regex that matches the prefix without anchoring
the end of the message (`^Use factory for`, not `^Use factory for pkg\.T$`) keeps matching every
diagnostic for `T`, with or without a suggested factory. An end-anchored regex written against the
old, suffix-less message stops matching once a factory suffix is appended.
