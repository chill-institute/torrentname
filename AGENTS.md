# Torrentname

`torrentname` is the `chill.institute` Go library for parsing torrent-style release names into structured metadata.

## Work

- `mise install`
- `mise run test`
- `mise run bench`
- `mise run verify`

## Conventions

- Keep parsing behavior deterministic and dependency-free.
- Keep release-token coverage catalog-driven: add common aliases to the relevant catalog in `catalog_tokens.go`, normalize through lookup helpers, and avoid one-off app-specific parser heuristics.
- Treat benchmarks and fuzz tests as part of parser safety, not optional extras.
- Preserve visible attribution to the original upstream library in public docs and package framing.

## Read More

- [Parser contract](./docs/SPEC.md)
- [Contributor checks](./CONTRIBUTING.md)
