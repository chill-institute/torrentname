# torrentname

![chill.institute library](https://chill.institute/banner.png)

Modern fork of [middelink/go-parse-torrent-name](https://github.com/middelink/go-parse-torrent-name) / [jzjzjzj/parse-torrent-name](https://github.com/jzjzjzj/parse-torrent-name)

`torrentname` parses torrent-style release names into structured metadata such as title, year, season, episode, quality, codec, audio tags, and release group.

## Install

Add the module:

```bash
go get github.com/chill-institute/torrentname
```

## Quickstart

Install the pinned toolchain:

```bash
mise install
```

Run the local verification baseline:

```bash
mise run verify
```

Basic usage:

```go
package main

import (
	"fmt"

	"github.com/chill-institute/torrentname"
)

func main() {
	info, err := torrentname.Parse("Hercules.2014.EXTENDED.1080p.WEB-DL.DD5.1.H264-RARBG")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", info)
}
```

## What It Does

Turn noisy release names into structured fields you can search, filter, deduplicate, or display.

TV example:

```text
Sample Series S05E03 720p HDTV x264-GRP
```

```go
&torrentname.TorrentInfo{
	Title:      "Sample Series",
	Season:     5,
	Episode:    3,
	Resolution: "720p",
	Quality:    "HDTV",
	Codec:      "x264",
	Group:      "GRP",
}
```

Movie example:

```text
Open.Feature.2014.EXTENDED.1080p.WEB-DL.DD5.1.H264-RARBG
```

```go
&torrentname.TorrentInfo{
	Title:      "Open Feature",
	Year:       2014,
	Resolution: "1080p",
	Quality:    "WEB-DL",
	Codec:      "H264",
	Audio:      "DD5.1",
	Group:      "RARBG",
	Extended:   true,
}
```

## Repo Tasks

Common local commands:

```bash
mise run fmt
mise run test
mise run bench
mise run test:fuzz
JACKETT_API_KEY=... mise run fixtures:jackett
mise run verify
```

## Benchmarking

We track parser performance in two ways:

- low-level microbenchmarks for `ns/op`, `B/op`, and `allocs/op`
- batch-style overhead translated into `ms per 1k rows`

Use `mise run bench` for the local benchmark suite.

See [CONTRIBUTING.md](./CONTRIBUTING.md) for the repeatable before/after benchmark workflow.

## Current Focus

- improve parser accuracy and coverage against the real-world fixture corpus
- keep tightening parser performance and allocations with measured benchmark passes
- evolve the public parser contract toward the full library spec in small, verified steps

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md)

## License

MIT. See [LICENSE](./LICENSE)
