# Delivery

`torrentname` is a tagged Go module. Go tooling resolves `vX.Y.Z` directly from
this repository.

## CI

Pull requests, pushes to `main`, and manual dispatches run `mise run verify`.
Pushes to `main` and dispatches also run the fuzz and benchmark smoke, which
gate the release:

```bash
mise run test:fuzz
go test . -run=^$ -bench=BenchmarkParse -benchmem -count=1
```

`mise run verify` runs the [`tasks.verify`](../mise.toml) dependencies.

## Releases

After all checks pass on `main`, semantic-release creates the immutable tag and
GitHub release from Conventional Commits. There is no version file in the source
tree.

## Operator Checklist

- Keep `main` as the release branch.
- Use Conventional Commits such as `fix: ...`, `feat: ...`, and `perf: ...`.
- Keep the `release` Environment and release permissions limited to the `main`
  release lane.
