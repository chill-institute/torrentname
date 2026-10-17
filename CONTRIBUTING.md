# Contributing

Parser changes should be deterministic, dependency-free at runtime, and backed
by representative release names.

## Start

```bash
mise install
mise run test
mise run verify
```

The [parser specification](./docs/SPEC.md) defines supported inference and
normalization. Add common aliases to `catalog_tokens.go`; avoid app-specific
special cases.

## Focused Checks

```bash
mise run corpus:metrics
mise run test:cover
mise run test:fuzz
mise run bench
```

Compare performance-sensitive changes with:

```bash
BENCH_OUT=tmp/bench/baseline.txt mise run bench:record
# make the change
BENCH_OUT=tmp/bench/current.txt mise run bench:record
mise run bench:compare
```

`mise run corpus:metrics` protects field-presence floors; it is not an accuracy
score. Review `ns/op`, `B/op`, and `allocs/op` for parser hot paths.

## Fixtures

```bash
JACKETT_API_KEY=... mise run fixtures:jackett
```

`JACKETT_BASE_URL` overrides the default local Jackett address.

Fixture generation strips configured URLs, credentials, download URLs, and
remote diagnostics. Commit only representative fixture changes and intentional,
minimized fuzz repro cases.

Keep the original upstream attribution visible when changing public framing.

## Release

After every [CI](./.github/workflows/ci.yml) job passes on `main`,
semantic-release tags `vX.Y.Z` and creates the immutable GitHub release; Go
tooling resolves the tag directly. Release types come from the
commit-analyzer `releaseRules` in [`.releaserc.json`](./.releaserc.json),
identical in every chill.institute package repo;
[`smoke.mjs`](./.github/release/smoke.mjs) fails when it drifts. Keep the
`release` Environment and release permissions limited to the `main` release
job.
