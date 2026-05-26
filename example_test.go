package torrentname_test

import (
	"fmt"

	"github.com/chill-institute/torrentname"
)

func ExampleParse_episode() {
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

	// Output:
	// Sample Series S05E03 720p HDTV GRP
}

func ExampleParse_movie() {
	info, err := torrentname.Parse("Open.Feature.2014.EXTENDED.1080p.WEB-DL.DD5.1.H264-RARBG")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s (%d) %s %s %s extended=%t\n",
		info.Title,
		info.Year,
		info.Resolution,
		info.Quality,
		info.Audio,
		info.Extended,
	)

	// Output:
	// Open Feature (2014) 1080p WEB-DL DD5.1 extended=true
}
