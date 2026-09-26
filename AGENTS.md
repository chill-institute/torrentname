# Torrentname

`torrentname` is the `chill.institute` Go library for parsing torrent-style release names into structured metadata.

## Work

- `mise install`

## Conventions

- Keep parsing behavior deterministic and dependency-free.
- Keep release-token coverage catalog-driven: add common aliases to the relevant catalog in `catalog_tokens.go`, normalize through lookup helpers, and avoid one-off app-specific parser heuristics.
- Preserve visible attribution to the original upstream library in public docs and package framing.

## Proof map

| Change | Check | Runs | Leaves |
| --- | --- | --- | --- |
| Docs | `mise run verify` | local, [CI](./.github/workflows/ci.yml) `verify` | exit status |
| Parser logic, catalogs, normalization | `mise run verify`, then `mise run test:fuzz` | local, [CI](./.github/workflows/ci.yml) `verify`; `fuzz` on push to `main` and dispatch | exit status; `FuzzParse` failures write repros to `testdata/fuzz/FuzzParse/` |
| Parser hot paths | `BENCH_OUT=<file> mise run bench:record` before and after, then `mise run bench:compare` | local; [CI](./.github/workflows/ci.yml) `benchmarks` runs one pass on push to `main` and dispatch | `tmp/bench/*.txt`, benchstat table |
| Fixture corpus | `JACKETT_API_KEY=... mise run fixtures:jackett`, then `mise run verify` | local, against a live local Jackett; see [Fixtures](./CONTRIBUTING.md#fixtures) | rewritten `testdata/jackett/*.json` |
| Workflows | `mise run actions` and `go test ./internal/workflowpolicy` (in `verify`) | local, [CI](./.github/workflows/ci.yml) `verify` | exit status |
| Pushed workflow changes | [shared scan](https://github.com/chill-institute/.github/tree/main/.github/actions/scan), last step of the [CI](./.github/workflows/ci.yml) `verify` job: Actionlint and Zizmor when the pushed range touches workflows; secrets rely on GitHub secret scanning | CI on push to `main` (pushed range) and dispatch (full history) | failed run |
| Release | [CI](./.github/workflows/ci.yml) `release` after all jobs on `main` | CI only, `release` Environment | tag and GitHub release for `feat`, `fix`, `perf`, `revert`, or breaking commits; see [Delivery](./docs/DELIVERY.md) |

`mise run verify` covers table tests in `parser_test.go`, golden cases in
`fixture_accuracy_test.go`, the `testdata/jackett` and `testdata/synthetic`
corpora, and `corpus:metrics` field-presence floors.

Gaps:

- No property tests; `FuzzParse` asserts only that `Parse` returns without error or panic. Owner: chill-institute/torrentname.
- No benchmark regression gate; comparison is manual. Owner: chill-institute/torrentname.
- No Markdown or link check for docs. Owner: chill-institute/torrentname.

## Read More

- [Parser contract](./docs/SPEC.md)
- [Contributor checks](./CONTRIBUTING.md)
