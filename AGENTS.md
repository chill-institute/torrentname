# Torrentname

Zero-dependency Go library that parses torrent-style release names into
`TorrentInfo`. The parser contract is [docs/spec.md](./docs/spec.md).

```bash
mise install
mise run verify
```

## Rules

- Parsing stays deterministic and dependency-free at runtime.
- Add release-token aliases to the catalogs in `catalog_tokens.go`; keep
  heuristics general to release names.
- Keep the upstream attribution in the README and the upstream copyright in
  [LICENSE](./LICENSE).
- Release types come from `releaseRules` in [`.releaserc.json`](./.releaserc.json),
  shared by every chill.institute package repo and checked by
  [`smoke.mjs`](./.github/release/smoke.mjs).

## Proof map

| Change | Check | Runs |
| --- | --- | --- |
| Parser, catalogs, normalization | `mise run verify`, then `mise run test:fuzz` | local; [CI](./.github/workflows/ci.yml) `verify`, `fuzz` on `main` |
| Hot paths | `BENCH_OUT=tmp/bench/baseline.txt mise run bench:record` before, `mise run bench:record` after, then `mise run bench:compare` | local; CI `benchmarks` on `main` |
| Fixture corpus | `JACKETT_API_KEY=... mise run fixtures:jackett`, then `mise run verify` | local, against a live Jackett |
| Workflows | `mise run actions`, `go test ./internal/workflowpolicy` (in `verify`) | local, CI |
| Release rules | `npm ci --prefix .github/release --ignore-scripts && node .github/release/smoke.mjs` | local, CI |
| Release | CI `release` job after all jobs pass on `main` | CI only: tag, GitHub release |

Gaps: `FuzzParse` checks only for panics and errors; benchmark regressions are
compared by hand.
