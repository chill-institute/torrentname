package torrentname

import "testing"

func assertTorrentInfo(t *testing.T, filename string, got TorrentInfo, want TorrentInfo) {
	t.Helper()
	if got != want {
		t.Fatalf("Parse(%q)\nwant: %#v\ngot:  %#v", filename, want, got)
	}
}

func TestParse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		filename string
		want     TorrentInfo
	}{
		{
			name:     "tv basic",
			filename: "Sample Series S05E03 720p HDTV x264-GRP",
			want: TorrentInfo{
				Title:      "Sample Series",
				Season:     5,
				Episode:    3,
				Resolution: "720p",
				Quality:    "HDTV",
				Codec:      "x264",
				Group:      "GRP",
			},
		},
		{
			name:     "movie bluray",
			filename: "Open Feature (2014) 1080p BrRip H264 - YIFY",
			want: TorrentInfo{
				Title:      "Open Feature",
				Year:       2014,
				Resolution: "1080p",
				Quality:    "BrRip",
				Codec:      "H264",
				Group:      "YIFY",
			},
		},
		{
			name:     "movie hdrip xvid",
			filename: "Open.Feature.2014.HDRip.XViD-EVO",
			want: TorrentInfo{
				Title:   "Open Feature",
				Year:    2014,
				Quality: "HDRip",
				Codec:   "XViD",
				Group:   "EVO",
			},
		},
		{
			name:     "extended web dl",
			filename: "Open.Feature.2014.EXTENDED.1080p.WEB-DL.DD5.1.H264-RARBG",
			want: TorrentInfo{
				Title:      "Open Feature",
				Year:       2014,
				Resolution: "1080p",
				Quality:    "WEB-DL",
				Codec:      "H264",
				Audio:      "DD5.1",
				Group:      "RARBG",
				Extended:   true,
			},
		},
		{
			name:     "x style episode",
			filename: "Sample Drama 5x06 HDTV x264-FoV",
			want: TorrentInfo{
				Title:   "Sample Drama",
				Season:  5,
				Episode: 6,
				Quality: "HDTV",
				Codec:   "x264",
				Group:   "FoV",
			},
		},
		{
			name:     "hardcoded audio",
			filename: "Archive.2014.HC.HDRip.XViD.AC3-SAFE",
			want: TorrentInfo{
				Title:     "Archive",
				Year:      2014,
				Quality:   "HDRip",
				Codec:     "XViD",
				Audio:     "AC3",
				Group:     "SAFE",
				Hardcoded: true,
			},
		},
		{
			name:     "dual audio size",
			filename: "Sample 2014 Dual-Audio WEBRip 1400Mb",
			want: TorrentInfo{
				Title:   "Sample",
				Year:    2014,
				Quality: "WEBRip",
				Audio:   "Dual-Audio",
				Size:    "1400Mb",
			},
		},
		{
			name:     "language tag",
			filename: "Example.Series.S02E20.rus.eng.720p.Kybik.v.Kybe",
			want: TorrentInfo{
				Title:      "Example Series",
				Season:     2,
				Episode:    20,
				Resolution: "720p",
				Language:   "rus.eng",
			},
		},
		{
			name:     "3d sbs",
			filename: "Open.Movie.2016.3D.1080p.BRRip.SBS.x264.AAC-ETRG",
			want: TorrentInfo{
				Title:      "Open Movie",
				Year:       2016,
				Resolution: "1080p",
				Quality:    "BRRip",
				Codec:      "x264",
				Audio:      "AAC",
				Sbs:        "SBS",
				Group:      "ETRG",
				ThreeD:     true,
			},
		},
		{
			name:     "half sbs",
			filename: "Open.Movie.2015.3D.1080p.BRRip.Half-SBS.x264.AAC-m2g",
			want: TorrentInfo{
				Title:      "Open Movie",
				Year:       2015,
				Resolution: "1080p",
				Quality:    "BRRip",
				Codec:      "x264",
				Audio:      "AAC",
				Sbs:        "Half-SBS",
				Group:      "m2g",
				ThreeD:     true,
			},
		},
		{
			name:     "unrated",
			filename: "Sample.Feature.2016.UNRATED.720p.BRRip.x264.AAC-ETRG",
			want: TorrentInfo{
				Title:      "Sample Feature",
				Year:       2016,
				Resolution: "720p",
				Quality:    "BRRip",
				Codec:      "x264",
				Audio:      "AAC",
				Group:      "ETRG",
				Unrated:    true,
			},
		},
		{
			name:     "website prefix",
			filename: "[SafeSource] Sample.Series.S05E10.480p.BluRay.x264-GAnGSteR",
			want: TorrentInfo{
				Title:      "Sample Series",
				Season:     5,
				Episode:    10,
				Resolution: "480p",
				Quality:    "BluRay",
				Codec:      "x264",
				Group:      "GAnGSteR",
				Website:    "SafeSource",
			},
		},
		{
			name:     "website prefix with dash",
			filename: "[ sample.source ] -Sample.Series.S07E07.720p.HDTV.x264-DIMENSION",
			want: TorrentInfo{
				Title:      "Sample Series",
				Season:     7,
				Episode:    7,
				Resolution: "720p",
				Quality:    "HDTV",
				Codec:      "x264",
				Group:      "DIMENSION",
				Website:    "sample.source",
			},
		},
		{
			name:     "anime style with container",
			filename: "[HorribleSubs] Sample Planet - 10 [480p].mkv",
			want: TorrentInfo{
				Title:      "Sample Planet",
				Episode:    10,
				Resolution: "480p",
				Container:  "mkv",
				Website:    "HorribleSubs",
			},
		},
		{
			name:     "high episode anime style",
			filename: "[HorribleSubs] Sample Conan - 862 [1080p].mkv",
			want: TorrentInfo{
				Title:      "Sample Conan",
				Episode:    862,
				Resolution: "1080p",
				Container:  "mkv",
				Website:    "HorribleSubs",
			},
		},
		{
			name:     "proper tag",
			filename: "Sample.Toon.S01E05.HDTV.x264.PROPER-LOL",
			want: TorrentInfo{
				Title:   "Sample Toon",
				Season:  1,
				Episode: 5,
				Quality: "HDTV",
				Codec:   "x264",
				Group:   "LOL",
				Proper:  true,
			},
		},
		{
			name:     "repack tag",
			filename: "Sample.Toon.S01E06.HDTV.x264.REPACK-LOL",
			want: TorrentInfo{
				Title:   "Sample Toon",
				Season:  1,
				Episode: 6,
				Quality: "HDTV",
				Codec:   "x264",
				Group:   "LOL",
				Repack:  true,
			},
		},
		{
			name:     "region and line audio",
			filename: "Classic.Feature.2012.R5.DVDRip.XViD.LiNE-UNiQUE",
			want: TorrentInfo{
				Title:   "Classic Feature",
				Year:    2012,
				Quality: "DVDRip",
				Codec:   "XViD",
				Audio:   "LiNE",
				Group:   "UNiQUE",
				Region:  "5",
			},
		},
		{
			name:     "widescreen",
			filename: "Demo.News.2014.WS.PDTV.x264-RKOFAN1990",
			want: TorrentInfo{
				Title:      "Demo News",
				Year:       2014,
				Quality:    "PDTV",
				Codec:      "x264",
				Group:      "RKOFAN1990",
				Widescreen: true,
			},
		},
		{
			name:     "container token",
			filename: "Sample Archive 2014 WEB-DL x264 MKV",
			want: TorrentInfo{
				Title:     "Sample Archive",
				Year:      2014,
				Quality:   "WEB-DL",
				Codec:     "x264",
				Container: "MKV",
			},
		},
		{
			name:     "dual audio with mixed languages",
			filename: "Demo.(2009).1080p.Dual Audio(Hindi+English) 5.1 Audios",
			want: TorrentInfo{
				Title:      "Demo",
				Year:       2009,
				Resolution: "1080p",
				Audio:      "Dual Audio",
			},
		},
		{
			name:     "dual audio bluray",
			filename: "Demo 2009 x264 720p Esub BluRay 6.0 Dual Audio English Hindi GOPISAHI",
			want: TorrentInfo{
				Title:      "Demo",
				Year:       2009,
				Resolution: "720p",
				Quality:    "BluRay",
				Codec:      "x264",
				Audio:      "Dual Audio",
			},
		},
		{
			name:     "complete season with standalone season marker",
			filename: "Sample Series S01 COMPLETE 720p WEBRip x264-GRP",
			want: TorrentInfo{
				Title:      "Sample Series",
				Season:     1,
				Resolution: "720p",
				Quality:    "WEBRip",
				Codec:      "x264",
				Group:      "GRP",
			},
		},
		{
			name:     "complete season with spaced codec and web dl",
			filename: "Sample Series S01 COMPLETE 1080p WEB DL DDP5 1 Atmos H 264-GRP",
			want: TorrentInfo{
				Title:      "Sample Series",
				Season:     1,
				Resolution: "1080p",
				Quality:    "WEB-DL",
				Codec:      "H264",
				Audio:      "DDP5.1",
				Group:      "GRP",
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := Parse(tc.filename)
			if err != nil {
				t.Fatalf("Parse(%q) error = %v", tc.filename, err)
			}
			assertTorrentInfo(t, tc.filename, *got, tc.want)
		})
	}
}

func TestParserDoesNotPanicOnJackettFixtureTitle(t *testing.T) {
	t.Parallel()

	_, err := Parse("[Hi-Res] Symphonic Suite AKIRA 2016 ハイパーハイレゾエディション／芸能山城組 (diff DSD256 11.2MHz タグ付き)")
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
}
