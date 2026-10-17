# Contributing

Back parser changes with representative release names; update the [spec](./docs/spec.md).

```bash
mise install
mise run verify          # format, vet, tests, corpus field-presence floors
mise run test:fuzz
BENCH_OUT=tmp/bench/baseline.txt mise run bench:record   # before the change
BENCH_OUT=tmp/bench/current.txt mise run bench:record    # after
mise run bench:compare
```

Refresh fixtures from a local Jackett with `JACKETT_API_KEY=... mise run
fixtures:jackett` (`JACKETT_BASE_URL` overrides the address). The generator
strips URLs and credentials; commit only representative fixtures and minimized
fuzz repros.

Merges to `main` release from Conventional Commits once [CI](./.github/workflows/ci.yml)
passes; Go resolves the `vX.Y.Z` tag directly.
