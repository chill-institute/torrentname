# Contributing

Thanks for contributing to `torrentname`

## Setup

Install the pinned toolchain:

```bash
mise install
```

Run the test suite once to confirm the repo is healthy:

```bash
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
mise run test:fuzz
JACKETT_API_KEY=... mise run fixtures:jackett
```

For benchmark comparisons, capture before and after runs with:

```bash
go test . -run=^$ -bench=BenchmarkParse -benchmem -count=10 > /tmp/before.txt
go test . -run=^$ -bench=BenchmarkParse -benchmem -count=10 > /tmp/after.txt
benchstat /tmp/before.txt /tmp/after.txt
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
- Tooling and contributor workflow belong in this repo. Workspace-wide operator docs belong in the workspace repo.

## Pull Requests

- Keep instructions current when setup or validation changes.
- Update docs when the public package surface or developer workflow changes.
- Add or update tests and benchmarks when parser behavior changes.
