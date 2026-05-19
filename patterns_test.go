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

func TestParseDoesNotMisreadMetadataBracketAsGroup(t *testing.T) {
	t.Parallel()

	titles := []string{
		"Widow's Bay (2026) s01e04 [Mkv - 1080p H264 - MultiLang Aac 2.0 - MultiSubs]",
		"Widow\u2019s Bay (2026) s01e04 [Mkv - 1080p H264 - MultiLang Aac 2.0 - MultiSubs]",
	}
	for _, title := range titles {
		info, err := Parse(title)
		if err != nil {
			t.Fatalf("Parse(%q) error: %v", title, err)
		}
		if info.Group != "" {
			t.Errorf("Parse(%q).Group = %q, want \"\"", title, info.Group)
		}
	}
}

func TestNormalizeModernAudioVariants(t *testing.T) {
	t.Parallel()

	cases := map[string]string{
		"EAC3.Atmos.5.1":   "EAC3 Atmos 5.1",
		"DTS.X.7.1":        "DTS X 7.1",
		"DTS.6CH":          "DTS 5.1",
		"DTS-HD.HRA.5.1":   "DTS-HD HRA 5.1",
		"EAC3.6CH":         "EAC3 5.1",
		"TrueHD.7.1.Atmos": "TrueHD Atmos 7.1",
		"DDPlus.Atmos":     "DD+ Atmos",
		"DDPlus.5.1":       "DD+ 5.1",
		"DD+.Atmos.5.1":    "DD+ Atmos 5.1",
		"DDPlus.5.1.Atmos": "DD+ Atmos 5.1",
		"6CH":              "5.1",
		"PCM.2.0":          "PCM 2.0",
		"LPCM.5.1":         "LPCM 5.1",
		"Opus.5.1":         "Opus 5.1",
	}
	for input, want := range cases {
		if got := normalizeAudioRich(input); got != want {
			t.Errorf("normalizeAudioRich(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestNormalizeLanguageVariants(t *testing.T) {
	t.Parallel()

	cases := map[string]string{
		"Eng":       "ENG",
		"English":   "ENG",
		"Jps":       "JPN",
		"MultiLang": "MULTI",
		"VOSTFR":    "VOSTFR",
	}
	for input, want := range cases {
		if got := normalizeLanguage(input); got != want {
			t.Errorf("normalizeLanguage(%q) = %q, want %q", input, got, want)
		}
	}
}
