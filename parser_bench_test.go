package torrentname

import "testing"

func BenchmarkParse(b *testing.B) {
	b.ReportAllocs()

	cases := []struct {
		name     string
		filename string
	}{
		{name: "tv_basic", filename: "Sample Series S05E03 720p HDTV x264-GRP"},
		{name: "movie_webdl", filename: "Open.Feature.2014.EXTENDED.1080p.WEB-DL.DD5.1.H264-RARBG"},
		{name: "anime_style", filename: "[HorribleSubs] Sample Planet - 10 [480p].mkv"},
		{name: "long_noisy", filename: "[Hi-Res] Symphonic Suite AKIRA 2016 ハイパーハイレゾエディション／芸能山城組 (diff DSD256 11.2MHz タグ付き)"},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			for range b.N {
				if _, err := Parse(tc.filename); err != nil {
					b.Fatalf("Parse(%q) returned error: %v", tc.filename, err)
				}
			}
		})
	}
}
