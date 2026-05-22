# Contributing

Thanks for contributing to `torrentname`

## Setup

Install the pinned toolchain and run the test suite once:

```bash
mise install
mise run test
```

## Validation

Run the full local verification baseline before opening or updating a pull request:

```bash
mise run verify
```

Optional parser-focused checks:

```bash
mise run bench
BENCH_OUT=tmp/bench/baseline.txt mise run bench:record
BENCH_OUT=tmp/bench/current.txt mise run bench:record
mise run bench:compare
mise run corpus:metrics
mise run test:cover
mise run test:fuzz
JACKETT_API_KEY=... mise run fixtures:jackett
```

Use the [Parser Spec](./docs/SPEC.md) to decide whether a parser behavior change is a contract expansion, a normalization fix, or an unsupported inference.

For benchmark comparisons, capture before and after runs with:

```bash
BENCH_OUT=tmp/bench/baseline.txt mise run bench:record
# change parser code
BENCH_OUT=tmp/bench/current.txt mise run bench:record
mise run bench:compare
```

For the search-style `1k rows` view, run:

```bash
go test . -run=^$ -bench=BenchmarkParseBatch1000 -benchmem -count=10
```

When discussing search-style workloads, translate parser cost into `ms per 1k rows`:

- `ns/op` is the per-title parse cost
- `ms per 1k rows = ns/op / 1_000_000`
- example: `35,000 ns/op` is about `35 ms` of parser overhead for `1,000` results

## Development Notes

- This repo owns a standalone Go parsing library with zero runtime dependencies.
- Keep parser behavior deterministic and fast.
- Prefer small, explicit refactors over broad speculative rewrites.
- Keep visible credit to the original `middelink/go-parse-torrent-name` project when updating public docs or package framing.
- Treat live Jackett fixtures as sanitized test inputs. Never check in API keys or raw Jackett download URLs.
- Treat generated fuzz artifacts as opt-in. Only commit minimized repro cases you intentionally want to keep under `testdata/fuzz/`
- Tooling and contributor workflow belong in this repo. Workspace-wide operator docs belong in the workspace repo.

## Pull Requests

- Keep instructions current when setup or validation changes.
- Update docs when the public package surface or developer workflow changes.
- Add or update tests and benchmarks when parser behavior changes.
