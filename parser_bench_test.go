package torrentname

import "testing"

var benchmarkSink *TorrentInfo

var benchmarkCases = []struct {
	name     string
	filename string
}{
	{name: "tv_basic", filename: "Sample Series S05E03 720p HDTV x264-GRP"},
	{name: "movie_webdl", filename: "Open.Feature.2014.EXTENDED.1080p.WEB-DL.DD5.1.H264-RARBG"},
	{name: "anime_style", filename: "[HorribleSubs] Sample Planet - 10 [480p].mkv"},
	{name: "long_noisy", filename: "[Hi-Res] Symphonic Suite AKIRA 2016 ハイパーハイレゾエディション／芸能山城組 (diff DSD256 11.2MHz タグ付き)"},
}

func BenchmarkParse(b *testing.B) {
	b.ReportAllocs()

	for _, tc := range benchmarkCases {
		b.Run(tc.name, func(b *testing.B) {
			for range b.N {
				info, err := Parse(tc.filename)
				if err != nil {
					b.Fatalf("Parse(%q) returned error: %v", tc.filename, err)
				}
				benchmarkSink = info
			}
		})
	}
}

func BenchmarkParseBatch1000(b *testing.B) {
	b.ReportAllocs()

	const batchSize = 1000
	batch := make([]string, batchSize)
	for i := range batch {
		batch[i] = benchmarkCases[i%len(benchmarkCases)].filename
	}

	for range b.N {
		for _, filename := range batch {
			info, err := Parse(filename)
			if err != nil {
				b.Fatalf("Parse(%q) returned error: %v", filename, err)
			}
			benchmarkSink = info
		}
	}

	rowsPerSecond := float64(batchSize*b.N) / b.Elapsed().Seconds()
	b.ReportMetric(rowsPerSecond, "rows/sec")
	b.ReportMetric(float64(b.Elapsed().Microseconds())/1000/float64(b.N), "ms/1krows")
}
