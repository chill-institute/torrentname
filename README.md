# torrentname

![chill.institute library](https://chill.institute/banner.png)

[![CI](https://github.com/chill-institute/torrentname/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/chill-institute/torrentname/actions/workflows/ci.yml?query=branch%3Amain)
[![Go Reference](https://pkg.go.dev/badge/github.com/chill-institute/torrentname.svg)](https://pkg.go.dev/github.com/chill-institute/torrentname)

`torrentname` is a zero-dependency Go parser for torrent-style release names.
It turns noisy filenames into structured metadata such as title, year, season,
episode, quality, codec, audio, HDR, source, and release group.

Modern fork of [middelink/go-parse-torrent-name](https://github.com/middelink/go-parse-torrent-name) / [jzjzjzj/parse-torrent-name](https://github.com/jzjzjzj/parse-torrent-name).

## Install

```bash
go get github.com/chill-institute/torrentname
```

## Quick Start

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

	fmt.Printf("%s S%02dE%02d %s %s %s\n",
		info.Title,
		info.Season,
		info.Episode,
		info.Resolution,
		info.Quality,
		info.Group,
	)
}
```

Output:

```text
Sample Series S05E03 720p HDTV GRP
```

## Examples

TV episode:

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

Movie release:

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

Complete season:

```text
Sample Series S01 COMPLETE 720p WEBRip x264-GRP
```

```go
&torrentname.TorrentInfo{
	Title:      "Sample Series",
	Season:     1,
	Resolution: "720p",
	Quality:    "WEBRip",
	Codec:      "x264",
	Group:      "GRP",
	Complete:   true,
}
```

## Supported Metadata

`Parse` fills the fields it can prove from the release name:

- identity: `Title`, `Year`, `Season`, `Episode`, `EpisodeEnd`, `Part`
- release traits: `Resolution`, `Quality`, `Codec`, `HDR`, `Audio`, `BitDepth`
- source traits: `Source`, `Group`, `Website`, `Language`, `Region`
- flags: `Extended`, `Hardcoded`, `Proper`, `Repack`, `Remastered`, `Widescreen`, `Unrated`, `ThreeD`, `IMAX`, `Complete`
- file traits: `Container`, `Sbs`, `Size`, `Excess`

See [Parser Spec](./docs/SPEC.md) for the full contract and normalization rules.

## Behavior

- Parsing is deterministic and best-effort.
- Missing or unrecognized fields stay at their Go zero value.
- The parser does not call external services or validate titles against a media catalog.
- Runtime code has no third-party dependencies.

## Docs

- [Parser Spec](./docs/SPEC.md)
- [Contributing](./CONTRIBUTING.md)
- [Release workflow](./docs/DELIVERY.md)

## Development

```bash
mise install
mise run verify
```

Contributor checks, fixture refreshes, corpus metrics, and benchmark workflows
live in [Contributing](./CONTRIBUTING.md).

## License

MIT. See [LICENSE](./LICENSE)
