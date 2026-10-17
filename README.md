# torrentname

![chill.institute library](https://chill.institute/banner.png)

[![CI](https://github.com/chill-institute/torrentname/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/chill-institute/torrentname/actions/workflows/ci.yml?query=branch%3Amain)
[![Go Reference](https://pkg.go.dev/badge/github.com/chill-institute/torrentname.svg)](https://pkg.go.dev/github.com/chill-institute/torrentname)

A zero-dependency Go parser for torrent-style release names. It extracts title,
year, season, episode, quality, codec, audio, source, and release-group metadata
without calling external services.

Forked from
[middelink/go-parse-torrent-name](https://github.com/middelink/go-parse-torrent-name)
and [jzjzjzj/parse-torrent-name](https://github.com/jzjzjzj/parse-torrent-name).

## Install

```bash
go get github.com/chill-institute/torrentname
```

## Use

```go
package main

import (
	"fmt"

	"github.com/chill-institute/torrentname"
)

func main() {
	info, err := torrentname.Parse("Sample.Series.S05E03.720p.HDTV.x264-GRP")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s S%02dE%02d %s\n", info.Title, info.Season, info.Episode, info.Group)
}
```

`Parse` is deterministic and best-effort. Unknown fields keep their Go zero
value; recognized aliases normalize to stable values.

## Contract

[`TorrentInfo`](https://pkg.go.dev/github.com/chill-institute/torrentname#TorrentInfo)
lists the parsed fields. The [parser specification](./docs/SPEC.md) defines
accepted tokens and normalization rules.

## Develop

```bash
mise install
mise run verify
```

[Contributing](./CONTRIBUTING.md) · [Security](./SECURITY.md) ·
[MIT License](./LICENSE)
