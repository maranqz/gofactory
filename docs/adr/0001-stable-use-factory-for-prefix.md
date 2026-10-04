# Keep `Use factory for pkg.T` as a stable message prefix

#51 adds a suffix naming accessible factories to the diagnostic message. A cleaner message could
have dropped or reworded the existing `Use factory for pkg.T` prefix instead of appending to it, but
users match on it with their own exclusion regexes (`severity.rules`, `nolint` patterns, CI
allowlists) and a changed prefix would silently break those. We decided to keep the prefix verbatim
and append the suffix after it, so `Use factory for pkg.T` on its own continues to match every
diagnostic for `T`, with or without a suggested factory.
