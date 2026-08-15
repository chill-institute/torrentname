# Delivery

`torrentname` is a tagged Go module. Go tooling resolves `vX.Y.Z` directly from
this repository.

## CI

Pull requests and pushes to `main` run the same guardrails:

```bash
mise run verify
mise run test:fuzz
go test . -run=^$ -bench=BenchmarkParse -benchmem -count=1
```

`mise run verify` covers formatting, module tidiness, static analysis, workflow
linting, parser tests, and corpus field-presence floors.

## Releases

After all checks pass on `main`, semantic-release creates the immutable tag and
GitHub release from Conventional Commits. There is no version file in the source
tree.

## Operator Checklist

- Keep `main` as the release branch.
- Use Conventional Commits such as `fix: ...`, `feat: ...`, and `perf: ...`.
- Keep the `release` Environment and release permissions limited to the `main`
  release lane.
