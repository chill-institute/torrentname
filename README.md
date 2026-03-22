# torrentname

![chill.institute library](https://chill.institute/banner.png)

Modern fork of [middelink/go-parse-torrent-name](https://github.com/middelink/go-parse-torrent-name) / [jzjzjzj/parse-torrent-name](https://github.com/jzjzjzj/parse-torrent-name)

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
