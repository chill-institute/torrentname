package torrentname

import "testing"

func TestPatternsHaveExpectedCaptureGroups(t *testing.T) {
	t.Parallel()

	for _, pattern := range patterns {
		if got := pattern.re.NumSubexp(); got != 2 {
			t.Fatalf("pattern %q capture groups = %d, want 2", pattern.name, got)
		}
	}
}

func TestNormalizeCodecVariants(t *testing.T) {
	t.Parallel()

	cases := map[string]string{
		"x264":  "x264",
		"x 264": "x264",
		"x.264": "x264",
		"x-264": "x264",
		"x_264": "x264",
		"X264":  "x264",
		"x265":  "x265",
		"x 265": "x265",
		"x.265": "x265",
		"x-265": "x265",
		"x_265": "x265",
		"H264":  "H264",
		"h 264": "H264",
		"H.264": "H264",
		"H-264": "H264",
		"H_264": "H264",
		"AVC":   "H264",
		"avc":   "H264",
		"H265":  "H265",
		"h 265": "H265",
		"H.265": "H265",
		"H-265": "H265",
		"H_265": "H265",
		"HEVC":  "H265",
		"hevc":  "H265",
		"AV1":   "AV1",
		"XviD":  "XViD",
	}
	for input, want := range cases {
		if got := normalizeCodec(input); got != want {
			t.Errorf("normalizeCodec(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestParseCodecVariantsInTitles(t *testing.T) {
	t.Parallel()

	cases := map[string]string{
		"Sample Movie 2020 1080p H-264 GRP": "H264",
		"Sample Movie 2020 1080p H_264 GRP": "H264",
		"Sample Movie 2020 1080p x-264 GRP": "x264",
		"Sample Movie 2020 1080p x-265 GRP": "x265",
		"Sample Movie 2020 1080p H-265 GRP": "H265",
	}
	for title, want := range cases {
		info, err := Parse(title)
		if err != nil {
			t.Errorf("Parse(%q) error: %v", title, err)
			continue
		}
		if info.Codec != want {
			t.Errorf("Parse(%q).Codec = %q, want %q", title, info.Codec, want)
		}
	}
}
