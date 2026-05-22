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
				Edition: "Dual Audio",
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
				Language:   "HIN ENG",
				Edition:    "Dual Audio",
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
				Language:   "ENG HIN",
				Edition:    "Dual Audio",
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
				Complete:   true,
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
				Audio:      "DDP5.1 Atmos",
				Group:      "GRP",
				Complete:   true,
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

func TestReleaseInfoExamples(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		filename string
		want     TorrentInfo
	}{
		{
			name:     "cr anime web dl",
			filename: "Kaya-chan.Isnt.Scary.S01E11.1080p.CR.WEB-DL.DUAL.AAC2.0.H.264-VARYG",
			want: TorrentInfo{
				Title:      "Kaya-chan Isnt Scary",
				Season:     1,
				Episode:    11,
				Resolution: "1080p",
				Quality:    "WEB-DL",
				Codec:      "H264",
				Audio:      "AAC2.0",
				Source:     "CR",
				Group:      "VARYG",
				Edition:    "Dual Audio",
			},
		},
		{
			name:     "multi hdr hevc",
			filename: "Monarch Legacy of Monsters S02E04 2160p HDR10Plus DV WEBRip 6CH x265 HEVC-P",
			want: TorrentInfo{
				Title:      "Monarch Legacy of Monsters",
				Season:     2,
				Episode:    4,
				Resolution: "2160p",
				Quality:    "WEBRip",
				Codec:      "H265",
				HDR:        "HDR10+ DV",
				Audio:      "5.1",
				Group:      "P",
			},
		},
		{
			name:     "bluray remux avc truehd",
			filename: "Zootopia 2 2025 1080p Blu ray Remux AVC TrueHD Atmos 7 1 CiNEPHiL",
			want: TorrentInfo{
				Title:      "Zootopia 2",
				Year:       2025,
				Resolution: "1080p",
				Quality:    "REMUX",
				Codec:      "H264",
				Audio:      "TrueHD Atmos 7.1",
				Group:      "CiNEPHiL",
			},
		},
		{
			name:     "uhd remux dv hdr",
			filename: "Nightmare Alley 2021 2160p UHD Blu ray Remux DV HDR HEVC TrueHD Atmos 7 1 CiNEPHiLES",
			want: TorrentInfo{
				Title:      "Nightmare Alley",
				Year:       2021,
				Resolution: "2160p",
				Quality:    "REMUX",
				Codec:      "H265",
				HDR:        "DV HDR",
				Audio:      "TrueHD Atmos 7.1",
				Group:      "CiNEPHiLES",
			},
		},
		{
			name:     "complete series bracket group",
			filename: "Crossing Jordan 2001 Complete Series Seasons 1 to 6 1080p WEB x264 [i_c]",
			want: TorrentInfo{
				Title:      "Crossing Jordan",
				Year:       2001,
				Resolution: "1080p",
				Quality:    "WEB",
				Codec:      "x264",
				Group:      "i_c",
				Complete:   true,
			},
		},
		{
			name:     "p2p scene modern eac3 atmos",
			filename: "The.Series.Title.2010.S01E01.ATVP.WEBDL-2160p.EAC3.Atmos.5.1.DV.HDR10Plus.h265-RlsGrp",
			want: TorrentInfo{
				Title:      "The Series Title",
				Year:       2010,
				Season:     1,
				Episode:    1,
				Resolution: "2160p",
				Quality:    "WEB-DL",
				Codec:      "H265",
				HDR:        "DV HDR10+",
				Audio:      "EAC3 Atmos 5.1",
				Source:     "ATVP",
				Group:      "RlsGrp",
			},
		},
		{
			name:     "bravia core dts x",
			filename: "Festival.Movie.2025.2160p.BCORE.WEB-DL.DTS.X.7.1.HDR10.HEVC-GRP",
			want: TorrentInfo{
				Title:      "Festival Movie",
				Year:       2025,
				Resolution: "2160p",
				Quality:    "WEB-DL",
				Codec:      "H265",
				HDR:        "HDR10",
				Audio:      "DTS X 7.1",
				Source:     "BCORE",
				Group:      "GRP",
			},
		},
		{
			name:     "yts style bare channel",
			filename: "Sample Action Chapter 3 - Final Protocol (2019) 2160p BRRip 5.1 10Bit x265 -YTS",
			want: TorrentInfo{
				Title:      "Sample Action Chapter 3 - Final Protocol",
				Year:       2019,
				Resolution: "2160p",
				Quality:    "BRRip",
				Codec:      "x265",
				Audio:      "5.1",
				Group:      "YTS",
				BitDepth:   "10-bit",
			},
		},
		{
			name:     "bracketed bare channel",
			filename: "Sample Crash (2026) [2160p] [WEBRip] [x265] [10bit] [5.1] [YTS.BZ]",
			want: TorrentInfo{
				Title:      "Sample Crash",
				Year:       2026,
				Resolution: "2160p",
				Quality:    "WEBRip",
				Codec:      "x265",
				Audio:      "5.1",
				Group:      "YTS.BZ",
				BitDepth:   "10-bit",
			},
		},
		{
			name:     "spaced dd atmos channel",
			filename: "Sample Substitute 1996 UHD BluRay 1080p DD Atmos 5 1 DoVi HDR10 x265 SM737",
			want: TorrentInfo{
				Title:      "Sample Substitute",
				Year:       1996,
				Resolution: "1080p",
				Quality:    "BluRay",
				Codec:      "x265",
				HDR:        "DV HDR10",
				Audio:      "DD Atmos 5.1",
				Group:      "SM737",
			},
		},
		{
			name:     "hdr10 plus punctuation and latin audio pair",
			filename: "Sample.Directors.The.Monster.2026.2160p.WEB-DL.DV.HDR10+.ENG.LAT.Atmos.H265.MP4-BTM",
			want: TorrentInfo{
				Title:      "Sample Directors The Monster",
				Year:       2026,
				Resolution: "2160p",
				Quality:    "WEB-DL",
				Codec:      "H265",
				HDR:        "DV HDR10+",
				Audio:      "Atmos",
				Group:      "BTM",
				Container:  "MP4",
				Language:   "ENG LAT",
			},
		},
		{
			name:     "numbered audio marker",
			filename: "Sample Heir S01E11 2026 2160p WEB-DL H 265 DDP2 0 2Audios-HHWEB",
			want: TorrentInfo{
				Title:      "Sample Heir",
				Season:     1,
				Episode:    11,
				Year:       2026,
				Resolution: "2160p",
				Quality:    "WEB-DL",
				Codec:      "H265",
				Audio:      "DDP2.0",
				Group:      "HHWEB",
				Edition:    "Dual Audio",
			},
		},
		{
			name:     "ddp atmos abbreviation",
			filename: "Sample Feature 2026 1080p NF WEB-DL DDPA 5 1 x264-GRP",
			want: TorrentInfo{
				Title:      "Sample Feature",
				Year:       2026,
				Resolution: "1080p",
				Quality:    "WEB-DL",
				Codec:      "x264",
				Audio:      "DDP Atmos 5.1",
				Source:     "NF",
				Group:      "GRP",
			},
		},
		{
			name:     "anime land single language and short multisub",
			filename: "[Anime Land] Sample Feature (2026) (UHD WEBRip 1080p HEVC HDR10 EAC3 M SUB) KOREAN [56EEAD89].mkv",
			want: TorrentInfo{
				Title:      "Sample Feature",
				Year:       2026,
				Resolution: "1080p",
				Quality:    "WEBRip",
				Codec:      "H265",
				HDR:        "HDR10",
				Audio:      "EAC3",
				Website:    "Anime Land",
				Language:   "KOR",
				BitDepth:   "",
				Edition:    "Multi Subs",
				Container:  "mkv",
			},
		},
		{
			name:     "anime land single japanese language",
			filename: "[Anime Land] Sample Feature (2026) (UHD WEBRip 1080p HEVC HDR10 EAC3 Atmos) JAPANESE [05F3E9AC].mkv",
			want: TorrentInfo{
				Title:      "Sample Feature",
				Year:       2026,
				Resolution: "1080p",
				Quality:    "WEBRip",
				Codec:      "H265",
				HDR:        "HDR10",
				Audio:      "EAC3 Atmos",
				Website:    "Anime Land",
				Language:   "JPN",
				Container:  "mkv",
			},
		},
		{
			name:     "bare dts with six channel audio",
			filename: "Sample.Movie.2024.1080p.BluRay.DTS.6CH.x264-GRP",
			want: TorrentInfo{
				Title:      "Sample Movie",
				Year:       2024,
				Resolution: "1080p",
				Quality:    "BluRay",
				Codec:      "x264",
				Audio:      "DTS 5.1",
				Group:      "GRP",
			},
		},
		{
			name:     "ddplus atmos",
			filename: "Sample.Series.S01E01.1080p.WEB-DL.DDPlus.Atmos.H264-GRP",
			want: TorrentInfo{
				Title:      "Sample Series",
				Season:     1,
				Episode:    1,
				Resolution: "1080p",
				Quality:    "WEB-DL",
				Codec:      "H264",
				Audio:      "DD+ Atmos",
				Group:      "GRP",
			},
		},
		{
			name:     "ddplus atmos channel",
			filename: "Sample.Series.S01E01.1080p.WEB-DL.DDPlus.Atmos.5.1.H264-GRP",
			want: TorrentInfo{
				Title:      "Sample Series",
				Season:     1,
				Episode:    1,
				Resolution: "1080p",
				Quality:    "WEB-DL",
				Codec:      "H264",
				Audio:      "DD+ Atmos 5.1",
				Group:      "GRP",
			},
		},
		{
			name:     "ddplus channel before atmos",
			filename: "Sample.Series.S01E01.1080p.WEB-DL.DDPlus.5.1.Atmos.H264-GRP",
			want: TorrentInfo{
				Title:      "Sample Series",
				Season:     1,
				Episode:    1,
				Resolution: "1080p",
				Quality:    "WEB-DL",
				Codec:      "H264",
				Audio:      "DD+ Atmos 5.1",
				Group:      "GRP",
			},
		},
		{
			name:     "roku pcm av1",
			filename: "Indie.Feature.2024.1080p.ROKU.WEBRip.PCM.2.0.AV1-GRP",
			want: TorrentInfo{
				Title:      "Indie Feature",
				Year:       2024,
				Resolution: "1080p",
				Quality:    "WEBRip",
				Codec:      "AV1",
				Audio:      "PCM 2.0",
				Source:     "ROKU",
				Group:      "GRP",
			},
		},
		{
			name:     "edition and newer source tokens",
			filename: "Sample.Movie.2024.Special.Edition.Open.Matte.1080p.NOW.WEB-DL.Opus.5.1.x265-GRP",
			want: TorrentInfo{
				Title:      "Sample Movie",
				Year:       2024,
				Resolution: "1080p",
				Quality:    "WEB-DL",
				Codec:      "x265",
				Audio:      "Opus 5.1",
				Source:     "NOW",
				Group:      "GRP",
				Edition:    "Special Edition Open Matte",
			},
		},
		{
			name:     "abema source token",
			filename: "Sample.Series.S01E01.1080p.ABEMA.WEB-DL.AAC2.0.H264-GRP",
			want: TorrentInfo{
				Title:      "Sample Series",
				Season:     1,
				Episode:    1,
				Resolution: "1080p",
				Quality:    "WEB-DL",
				Codec:      "H264",
				Audio:      "AAC2.0",
				Source:     "ABEMA",
				Group:      "GRP",
			},
		},
		{
			name:     "bili source token",
			filename: "Sample.Anime.S01E01.1080p.BILI.WEB-DL.AAC2.0.H264-GRP",
			want: TorrentInfo{
				Title:      "Sample Anime",
				Season:     1,
				Episode:    1,
				Resolution: "1080p",
				Quality:    "WEB-DL",
				Codec:      "H264",
				Audio:      "AAC2.0",
				Source:     "BILI",
				Group:      "GRP",
			},
		},
		{
			name:     "bare opus stays in title",
			filename: "Mr.Hollands.Opus.1995.1080p.BluRay.x264-GRP",
			want: TorrentInfo{
				Title:      "Mr Hollands Opus",
				Year:       1995,
				Resolution: "1080p",
				Quality:    "BluRay",
				Codec:      "x264",
				Group:      "GRP",
			},
		},
		{
			name:     "bare pcm stays in title",
			filename: "PCM.Movie.2024.1080p.WEB-DL.x264-GRP",
			want: TorrentInfo{
				Title:      "PCM Movie",
				Year:       2024,
				Resolution: "1080p",
				Quality:    "WEB-DL",
				Codec:      "x264",
				Group:      "GRP",
			},
		},
		{
			name:     "season range complete with six channel audio",
			filename: "Sample.Series.2025.Season.1-2.COMPLETE.SERIES.1080p.H265.EAC3.6CH-GRP",
			want: TorrentInfo{
				Title:      "Sample Series",
				Year:       2025,
				Season:     1,
				Resolution: "1080p",
				Codec:      "H265",
				Audio:      "EAC3 5.1",
				Group:      "GRP",
				Complete:   true,
			},
		},
		{
			name:     "season word source before quality",
			filename: "Sample.Series.Season.1.NF.WEB-DL.x264-GRP",
			want: TorrentInfo{
				Title:   "Sample Series",
				Season:  1,
				Quality: "WEB-DL",
				Codec:   "x264",
				Source:  "NF",
				Group:   "GRP",
			},
		},
		{
			name:     "standalone max source",
			filename: "Sample.Movie.2024.1080p.MAX.WEB-DL.x264-GRP",
			want: TorrentInfo{
				Title:      "Sample Movie",
				Year:       2024,
				Resolution: "1080p",
				Quality:    "WEB-DL",
				Codec:      "x264",
				Source:     "MAX",
				Group:      "GRP",
			},
		},
		{
			name:     "terminal ambiguous provider remains source",
			filename: "Sample.Movie.2024.1080p.WEB-DL.ROKU",
			want: TorrentInfo{
				Title:      "Sample Movie",
				Year:       2024,
				Resolution: "1080p",
				Quality:    "WEB-DL",
				Source:     "ROKU",
			},
		},
		{
			name:     "common provider before resolution and quality remains source",
			filename: "Sample.Series.S01E01.NOW.1080p.WEB-DL.x264-GRP",
			want: TorrentInfo{
				Title:      "Sample Series",
				Season:     1,
				Episode:    1,
				Resolution: "1080p",
				Quality:    "WEB-DL",
				Codec:      "x264",
				Source:     "NOW",
				Group:      "GRP",
			},
		},
		{
			name:     "multi language subtitles",
			filename: "Sample.Collection.2003.Eng.Rus.Multi-Subs.1080p.H264-mp4-GRP",
			want: TorrentInfo{
				Title:      "Sample Collection",
				Year:       2003,
				Resolution: "1080p",
				Codec:      "H264",
				Container:  "mp4",
				Group:      "GRP",
				Language:   "ENG RUS MULTI",
				Edition:    "Multi Subs",
			},
		},
		{
			name:     "multi subtitles without language",
			filename: "Sample.Movie.2024.1080p.WEB-DL.H264.Multi-Subs-GRP",
			want: TorrentInfo{
				Title:      "Sample Movie",
				Year:       2024,
				Resolution: "1080p",
				Quality:    "WEB-DL",
				Codec:      "H264",
				Group:      "GRP",
				Edition:    "Multi Subs",
			},
		},
		{
			name:     "multilang with multi subtitles keeps language",
			filename: "Sample.Movie.2024.MultiLang.Multi-Subs.1080p.WEB-DL.H264-GRP",
			want: TorrentInfo{
				Title:      "Sample Movie",
				Year:       2024,
				Resolution: "1080p",
				Quality:    "WEB-DL",
				Codec:      "H264",
				Group:      "GRP",
				Language:   "MULTI",
				Edition:    "Multi Subs",
			},
		},
		{
			name:     "music style hi res flac",
			filename: "Sample.Artist.Sample.Album.2026.Hi-Res.FLAC.24Bit.96kHz-GRP",
			want: TorrentInfo{
				Title:    "Sample Artist Sample Album",
				Year:     2026,
				Audio:    "FLAC",
				Group:    "GRP",
				BitDepth: "24-bit",
			},
		},
		{
			name:     "source-like final group is not source",
			filename: "Sample.Movie.2024.1080p.WEB-DL.H264-PLAY",
			want: TorrentInfo{
				Title:      "Sample Movie",
				Year:       2024,
				Resolution: "1080p",
				Quality:    "WEB-DL",
				Codec:      "H264",
				Group:      "PLAY",
			},
		},
		{
			name:     "source-like dot suffix group is not source",
			filename: "Sample.Series.S01E01.720p.HDTV.x264.PLAY",
			want: TorrentInfo{
				Title:      "Sample Series",
				Season:     1,
				Episode:    1,
				Resolution: "720p",
				Quality:    "HDTV",
				Codec:      "x264",
				Group:      "PLAY",
			},
		},
		{
			name:     "source-like final group before container is not source",
			filename: "Sample.Movie.2024.1080p.WEB-DL.H264-PLAY.mkv",
			want: TorrentInfo{
				Title:      "Sample Movie",
				Year:       2024,
				Resolution: "1080p",
				Quality:    "WEB-DL",
				Codec:      "H264",
				Group:      "PLAY",
				Container:  "mkv",
			},
		},
		{
			name:     "source-like provider group before container is not source",
			filename: "Sample.Movie.2024.1080p.WEB-DL.H264-NF.mkv",
			want: TorrentInfo{
				Title:      "Sample Movie",
				Year:       2024,
				Resolution: "1080p",
				Quality:    "WEB-DL",
				Codec:      "H264",
				Group:      "NF",
				Container:  "mkv",
			},
		},
		{
			name:     "source-like dot suffix group before container is not source",
			filename: "Sample.Series.S01E01.720p.HDTV.x264.PLAY.mkv",
			want: TorrentInfo{
				Title:      "Sample Series",
				Season:     1,
				Episode:    1,
				Resolution: "720p",
				Quality:    "HDTV",
				Codec:      "x264",
				Group:      "PLAY",
				Container:  "mkv",
			},
		},
		{
			name:     "source-like dot provider group before container is not source",
			filename: "Sample.Movie.2024.1080p.WEB-DL.H264.NF.mkv",
			want: TorrentInfo{
				Title:      "Sample Movie",
				Year:       2024,
				Resolution: "1080p",
				Quality:    "WEB-DL",
				Codec:      "H264",
				Group:      "NF",
				Container:  "mkv",
			},
		},
		{
			name:     "source-like bracket group is not source",
			filename: "Sample.Movie.2024.1080p.WEB-DL.H264.[PLAY]",
			want: TorrentInfo{
				Title:      "Sample Movie",
				Year:       2024,
				Resolution: "1080p",
				Quality:    "WEB-DL",
				Codec:      "H264",
				Group:      "PLAY",
			},
		},
		{
			name:     "source-like bracket provider group before container is not source",
			filename: "Sample.Movie.2024.1080p.WEB-DL.H264.[NF].mkv",
			want: TorrentInfo{
				Title:      "Sample Movie",
				Year:       2024,
				Resolution: "1080p",
				Quality:    "WEB-DL",
				Codec:      "H264",
				Group:      "NF",
				Container:  "mkv",
			},
		},
		{
			name:     "rm4k remastered marker",
			filename: "Sample.Feature.1999.1080p.RM4K.BluRay.x265.EAC3.6CH-GRP",
			want: TorrentInfo{
				Title:      "Sample Feature",
				Year:       1999,
				Resolution: "1080p",
				Quality:    "BluRay",
				Codec:      "x265",
				Audio:      "EAC3 5.1",
				Group:      "GRP",
				Remastered: true,
			},
		},
		{
			name:     "truehd channel before atmos",
			filename: "Sample.Feature.2025.2160p.BluRay.Remux.HEVC.TrueHD.7.1.Atmos-GRP",
			want: TorrentInfo{
				Title:      "Sample Feature",
				Year:       2025,
				Resolution: "2160p",
				Quality:    "REMUX",
				Codec:      "H265",
				Audio:      "TrueHD Atmos 7.1",
				Group:      "GRP",
			},
		},
		{
			name:     "legacy dot suffix group",
			filename: "Sample.Series.S01E01.720p.HDTV.x264.GRP",
			want: TorrentInfo{
				Title:      "Sample Series",
				Season:     1,
				Episode:    1,
				Resolution: "720p",
				Quality:    "HDTV",
				Codec:      "x264",
				Group:      "GRP",
			},
		},
		{
			name:     "terminal codec metadata is not group",
			filename: "Sample.Movie.2024.1080p.WEB-DL.AV1",
			want: TorrentInfo{
				Title:      "Sample Movie",
				Year:       2024,
				Resolution: "1080p",
				Quality:    "WEB-DL",
				Codec:      "AV1",
			},
		},
		{
			name:     "terminal audio metadata is not group",
			filename: "Sample.Movie.2024.1080p.WEB-DL.DTS",
			want: TorrentInfo{
				Title:      "Sample Movie",
				Year:       2024,
				Resolution: "1080p",
				Quality:    "WEB-DL",
				Audio:      "DTS",
			},
		},
		{
			name:     "terminal rich audio metadata is not group",
			filename: "Sample.Movie.2024.1080p.WEB-DL.FLAC",
			want: TorrentInfo{
				Title:      "Sample Movie",
				Year:       2024,
				Resolution: "1080p",
				Quality:    "WEB-DL",
				Audio:      "FLAC",
			},
		},
		{
			name:     "terminal flag metadata is not group",
			filename: "Sample.Movie.2024.1080p.WEB-DL.PROPER",
			want: TorrentInfo{
				Title:      "Sample Movie",
				Year:       2024,
				Resolution: "1080p",
				Quality:    "WEB-DL",
				Proper:     true,
			},
		},
		{
			name:     "korsub hardcoded marker",
			filename: "Sample.Series.S01E01.KORSUB.1080p.WEBRip.x264-GRP",
			want: TorrentInfo{
				Title:      "Sample Series",
				Season:     1,
				Episode:    1,
				Resolution: "1080p",
				Quality:    "WEBRip",
				Codec:      "x264",
				Group:      "GRP",
				Hardcoded:  true,
			},
		},
		{
			name:     "standalone channel count",
			filename: "Sample.Series.S01E01.2160p.HDR10Plus.DV.WEBRip.6CH.x265.HEVC-GRP",
			want: TorrentInfo{
				Title:      "Sample Series",
				Season:     1,
				Episode:    1,
				Resolution: "2160p",
				Quality:    "WEBRip",
				Codec:      "H265",
				HDR:        "HDR10+ DV",
				Audio:      "5.1",
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

func TestReleaseInfoEdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		filename string
		check    func(*testing.T, TorrentInfo)
	}{
		{
			name:     "it title is not source",
			filename: "It 2017 1080p BluRay x264",
			check: func(t *testing.T, got TorrentInfo) {
				t.Helper()
				if got.Source != "" {
					t.Fatalf("Source = %q, want empty", got.Source)
				}
			},
		},
		{
			name:     "critical title is not cr source",
			filename: "Critical Role S01E01 1080p WEB x264",
			check: func(t *testing.T, got TorrentInfo) {
				t.Helper()
				if got.Source != "" {
					t.Fatalf("Source = %q, want empty", got.Source)
				}
			},
		},
		{
			name:     "hdrip is not hdr",
			filename: "Movie 2025 HDRip XViD AC3",
			check: func(t *testing.T, got TorrentInfo) {
				t.Helper()
				if got.Quality != "HDRip" || got.HDR != "" {
					t.Fatalf("Quality/HDR = %q/%q, want HDRip/empty", got.Quality, got.HDR)
				}
			},
		},
		{
			name:     "dvdrip is not dolby vision",
			filename: "Movie 2025 DVDRip XViD",
			check: func(t *testing.T, got TorrentInfo) {
				t.Helper()
				if got.Quality != "DVDRip" || got.HDR != "" {
					t.Fatalf("Quality/HDR = %q/%q, want DVDRip/empty", got.Quality, got.HDR)
				}
			},
		},
		{
			name:     "episode range",
			filename: "Show.Name.S01E01-E03.1080p.WEB-DL.x264-GRP",
			check: func(t *testing.T, got TorrentInfo) {
				t.Helper()
				if got.Season != 1 || got.Episode != 1 || got.EpisodeEnd != 3 {
					t.Fatalf("season/episode range = S%dE%d-E%d, want S1E1-E3", got.Season, got.Episode, got.EpisodeEnd)
				}
			},
		},
		{
			name:     "edition and bit depth",
			filename: "Movie 2025 Directors Cut 2160p REMUX HEVC 10bit TrueHD Atmos 7 1-GRP",
			check: func(t *testing.T, got TorrentInfo) {
				t.Helper()
				if got.Edition != "Director's Cut" || got.BitDepth != "10-bit" {
					t.Fatalf("edition/bit_depth = %q/%q, want Director's Cut/10-bit", got.Edition, got.BitDepth)
				}
			},
		},
		{
			name:     "anime crc is not release group",
			filename: "[Group] Sample Anime - 001v2 [1080p][HEVC 10bit][AAC][MultiSub][A1B2C3D4].mkv",
			check: func(t *testing.T, got TorrentInfo) {
				t.Helper()
				if got.Title != "Sample Anime" || got.Episode != 1 {
					t.Fatalf("title/episode = %q/%d, want Sample Anime/1", got.Title, got.Episode)
				}
				if got.Group != "" {
					t.Fatalf("Group = %q, want empty for CRC-style bracket", got.Group)
				}
			},
		},
		{
			name:     "season range does not become group",
			filename: "Sample.Series.2025.Season.1-3.Complete.1080p.WEB-DL.x265",
			check: func(t *testing.T, got TorrentInfo) {
				t.Helper()
				if got.Season != 1 || !got.Complete {
					t.Fatalf("season/complete = %d/%v, want 1/true", got.Season, got.Complete)
				}
				if got.Group != "" {
					t.Fatalf("Group = %q, want empty for season range", got.Group)
				}
			},
		},
		{
			name:     "terminal provider remains source",
			filename: "Sample.Movie.2024.1080p.WEB-DL.NF",
			check: func(t *testing.T, got TorrentInfo) {
				t.Helper()
				if got.Source != "NF" {
					t.Fatalf("Source = %q, want NF", got.Source)
				}
				if got.Group != "" {
					t.Fatalf("Group = %q, want empty", got.Group)
				}
			},
		},
		{
			name:     "terminal hdtv quality is not group",
			filename: "Sample.Series.S01E01.720p.HDTV",
			check: func(t *testing.T, got TorrentInfo) {
				t.Helper()
				if got.Quality != "HDTV" {
					t.Fatalf("Quality = %q, want HDTV", got.Quality)
				}
				if got.Group != "" {
					t.Fatalf("Group = %q, want empty", got.Group)
				}
			},
		},
		{
			name:     "terminal hdrip quality is not group",
			filename: "Sample.Feature.2024.HDRip",
			check: func(t *testing.T, got TorrentInfo) {
				t.Helper()
				if got.Quality != "HDRip" {
					t.Fatalf("Quality = %q, want HDRip", got.Quality)
				}
				if got.Group != "" {
					t.Fatalf("Group = %q, want empty", got.Group)
				}
			},
		},
		{
			name:     "episode title play is not source",
			filename: "Sample.Series.S01E01.Play.Time.1080p.WEB-DL.x264-GRP",
			check: func(t *testing.T, got TorrentInfo) {
				t.Helper()
				if got.Source != "" {
					t.Fatalf("Source = %q, want empty", got.Source)
				}
				if got.Group != "GRP" {
					t.Fatalf("Group = %q, want GRP", got.Group)
				}
			},
		},
		{
			name:     "episode title now is not source",
			filename: "Sample.Series.S01E01.Now.What.1080p.WEB-DL.x264-GRP",
			check: func(t *testing.T, got TorrentInfo) {
				t.Helper()
				if got.Source != "" {
					t.Fatalf("Source = %q, want empty", got.Source)
				}
				if got.Group != "GRP" {
					t.Fatalf("Group = %q, want GRP", got.Group)
				}
			},
		},
		{
			name:     "rejected source word does not hide following provider",
			filename: "Sample.Series.S01E01.Now.NF.WEB-DL.x264-GRP",
			check: func(t *testing.T, got TorrentInfo) {
				t.Helper()
				if got.Source != "NF" {
					t.Fatalf("Source = %q, want NF", got.Source)
				}
				if got.Group != "GRP" {
					t.Fatalf("Group = %q, want GRP", got.Group)
				}
			},
		},
		{
			name:     "episode title stan is not source",
			filename: "Sample.Series.S01E01.Stan.Knows.Best.1080p.WEB-DL.x264-GRP",
			check: func(t *testing.T, got TorrentInfo) {
				t.Helper()
				if got.Source != "" {
					t.Fatalf("Source = %q, want empty", got.Source)
				}
				if got.Group != "GRP" {
					t.Fatalf("Group = %q, want GRP", got.Group)
				}
			},
		},
		{
			name:     "common stan provider before resolution and quality remains source",
			filename: "Sample.Series.S01E01.STAN.1080p.WEB-DL.x264-GRP",
			check: func(t *testing.T, got TorrentInfo) {
				t.Helper()
				if got.Source != "STAN" {
					t.Fatalf("Source = %q, want STAN", got.Source)
				}
				if got.Group != "GRP" {
					t.Fatalf("Group = %q, want GRP", got.Group)
				}
			},
		},
		{
			name:     "ambiguous provider before resolution hdr and quality remains source",
			filename: "Sample.Series.S01E01.STAN.2160p.HDR10.WEB-DL.x264-GRP",
			check: func(t *testing.T, got TorrentInfo) {
				t.Helper()
				if got.Source != "STAN" {
					t.Fatalf("Source = %q, want STAN", got.Source)
				}
				if got.HDR != "HDR10" || got.Quality != "WEB-DL" {
					t.Fatalf("HDR/Quality = %q/%q, want HDR10/WEB-DL", got.HDR, got.Quality)
				}
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
			tc.check(t, *got)
		})
	}
}
