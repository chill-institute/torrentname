# Torrentname

`torrentname` is the `chill.institute` Go library for parsing torrent-style release names into structured metadata.

## Stack

- Go library with zero runtime dependencies
- `mise` for toolchain and repo tasks

## Commands

- `mise install`
- `mise run fmt`
- `mise run test`
- `mise run bench`
- `mise run verify`

## Conventions

- Keep parsing behavior deterministic and dependency-free.
- Keep release-token coverage catalog-driven: add common aliases to the relevant catalog in `catalog_tokens.go`, normalize through lookup helpers, and avoid one-off app-specific parser heuristics.
- Prefer typed helpers and explicit data flow when modernizing internals.
- Treat benchmarks and fuzz tests as part of parser safety, not optional extras.
- Preserve visible attribution to the original upstream library in public docs and package framing.

## Read More

- contributor workflow: [CONTRIBUTING.md](./CONTRIBUTING.md)
