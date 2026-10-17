# torrentname

![chill.institute library](https://chill.institute/banner.png)

[![CI](https://github.com/chill-institute/torrentname/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/chill-institute/torrentname/actions/workflows/ci.yml?query=branch%3Amain)
[![Go Reference](https://pkg.go.dev/badge/github.com/chill-institute/torrentname.svg)](https://pkg.go.dev/github.com/chill-institute/torrentname)

A zero-dependency Go parser for torrent-style release names: title, year,
season, episode, quality, codec, audio, source, and release group. Forked from
[middelink/go-parse-torrent-name](https://github.com/middelink/go-parse-torrent-name)
and [jzjzjzj/parse-torrent-name](https://github.com/jzjzjzj/parse-torrent-name).

```bash
go get github.com/chill-institute/torrentname
```

```go
info, err := torrentname.Parse("Sample.Series.S05E03.720p.HDTV.x264-GRP")
// info.Title == "Sample Series", info.Season == 5, info.Episode == 3
```

Unknown fields keep their zero value. [`TorrentInfo`](https://pkg.go.dev/github.com/chill-institute/torrentname#TorrentInfo)
lists the fields; the [parser spec](./docs/spec.md) defines tokens and
normalization.

[Contributing](./CONTRIBUTING.md) · [Security](./SECURITY.md) · [MIT License](./LICENSE)
