# torrentname

![chill.institute library](https://chill.institute/banner.png)

Modern fork of [middlelink/go-parse-torrent-name](https://github.com/middelink/go-parse-torrent-name) / [jzjzjzj/parse-torrent-name](https://github.com/jzjzjzj/parse-torrent-name)

This repo started as a fork of [middelink/go-parse-torrent-name](https://github.com/middelink/go-parse-torrent-name), and it keeps the original MIT-licensed lineage visible while we iterate toward a cleaner, faster, more battle-tested parser.

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
JACKETT_API_KEY=... mise run fixtures:jackett
mise run verify
```

## Current Focus

- modernize the repo shell and developer workflow first
- preserve the original library credit and MIT license lineage
- iterate on parser accuracy, performance, and API ergonomics in small, verified steps

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md)

## License

MIT

Original library credit:

- [middelink/go-parse-torrent-name](https://github.com/middelink/go-parse-torrent-name)
