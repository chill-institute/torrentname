package torrentname

import "testing"

func FuzzParse(f *testing.F) {
	seeds := []string{
		"Sample Series S05E03 720p HDTV x264-GRP",
		"Open.Feature.2014.EXTENDED.1080p.WEB-DL.DD5.1.H264-RARBG",
		"[SafeSource] Sample.Series.S05E10.480p.BluRay.x264-GAnGSteR",
		"[HorribleSubs] Sample Planet - 10 [480p].mkv",
		"Sample.Series.S01E01E02E03.720p.WEB-DL.x264-GRP",
		"Sample Series 2x01-2x02 HDTV x264-GRP",
		"[Hi-Res] Symphonic Suite AKIRA 2016 ハイパーハイレゾエディション／芸能山城組 (diff DSD256 11.2MHz タグ付き)",
		"",
		"[] -",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, filename string) {
		if _, err := Parse(filename); err != nil {
			t.Fatalf("Parse(%q) returned error: %v", filename, err)
		}
	})
}
