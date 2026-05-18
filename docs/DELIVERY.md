# Delivery

`torrentname` is delivered as a Go module. There is no deploy pipeline because the repository does not own a running service, environment, image, or hosted artifact to promote.

## CI

Pull requests and pushes to `main` run the same guardrails:

```bash
mise run verify
mise run test:fuzz
go test . -run=^$ -bench=BenchmarkParse -benchmem -count=1
```

The GitHub Actions workflow keeps default permissions empty, grants each verification job read-only repository access, pins high-trust actions to full commit SHAs, and uses the repo's `mise.toml` as the toolchain source of truth.

## Releases

Pushes to `main` run release only after verify, fuzz, and benchmark smoke jobs pass. Releases use semantic-release to analyze Conventional Commits, create a `vX.Y.Z` Git tag, and publish GitHub release notes. Go consumers then resolve the module through the tagged Git history.

The release job uses the `release` GitHub Environment as the release trust boundary with deployment records disabled. That Environment is branch-restricted to `main`. The job requests `contents: write` only for the release lane so the default repository token posture can stay read-only for non-release jobs.

The release lane does not publish to a package registry, build binaries, or push a version-bump commit because this is a tag-delivered Go library with no source version file.

## Operator Checklist

- Keep `main` as the release branch.
- Use Conventional Commits such as `fix: ...`, `feat: ...`, and `perf: ...`.
- Keep the `release` Environment limited to `main`.
- Keep the release job's `contents: write` permission allowed to create tags and GitHub Releases, or replace it with a narrowly scoped release bot if branch/tag rules require an allowlisted actor.
- Do not add deploy credentials to this repository unless it grows a real service artifact and environment smoke target.
