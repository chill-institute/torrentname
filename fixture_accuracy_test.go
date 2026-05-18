package torrentname

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

var fixtureGoldenCases = []struct {
	name   string
	title  string
	want   TorrentInfo
	fields []string
}{
	{
		name:  "movie bluray with audio and group",
		title: "Delivery Run 2024 BluRay 1080p DDP 5 1 x264 hallowed",
		want: TorrentInfo{
			Title:      "Delivery Run",
			Year:       2024,
			Resolution: "1080p",
			Quality:    "BluRay",
			Codec:      "x264",
			Audio:      "DDP5.1",
			Group:      "hallowed",
		},
		fields: []string{"title", "year", "resolution", "quality", "codec", "audio", "group"},
	},
	{
		name:  "anime dual audio multisub",
		title: "[Anime Time] Sentenced To Be A Hero (Yuusha-kei ni Shosu: Choubatsu Yuusha 9004 Tai Keimu Kiroku) - S01E05 [1080p][HEVC 10bit x265][AAC][Dual-Audio] [Multi Sub] [Weekly]",
		want: TorrentInfo{
			Title:      "Sentenced To Be A Hero",
			Season:     1,
			Episode:    5,
			Resolution: "1080p",
			Codec:      "H265",
			Audio:      "AAC",
			Website:    "Anime Time",
			BitDepth:   "10-bit",
			Edition:    "Dual Audio Multi Subs",
		},
		fields: []string{"title", "season", "episode", "resolution", "codec", "audio", "website", "bit_depth", "edition"},
	},
	{
		name:  "complete series bracket group",
		title: "Crossing Jordan 2001 Complete Series Seasons 1 to 6 1080p WEB x264 [i_c]",
		want: TorrentInfo{
			Title:      "Crossing Jordan",
			Year:       2001,
			Resolution: "1080p",
			Quality:    "WEB",
			Codec:      "x264",
			Group:      "i_c",
			Complete:   true,
		},
		fields: []string{"title", "year", "resolution", "quality", "codec", "group", "complete"},
	},
	{
		name:  "cr web dl source",
		title: "Kaya chan Isnt Scary S01E11 1080p CR WEB DL DUAL AAC2 0 H 264 VARYG (Kaya chan wa Kowakunai  Dual Audio  Multi Subs)",
		want: TorrentInfo{
			Title:      "Kaya chan Isnt Scary",
			Season:     1,
			Episode:    11,
			Resolution: "1080p",
			Quality:    "WEB-DL",
			Codec:      "H264",
			Audio:      "AAC2.0",
			Source:     "CR",
			Edition:    "Dual Audio Multi Subs",
		},
		fields: []string{"title", "season", "episode", "resolution", "quality", "codec", "audio", "source", "edition"},
	},
	{
		name:  "multi hdr hevc",
		title: "Monarch Legacy of Monsters S02E04 2160p HDR10Plus DV WEBRip 6CH x265 HEVC-P",
		want: TorrentInfo{
			Title:      "Monarch Legacy of Monsters",
			Season:     2,
			Episode:    4,
			Resolution: "2160p",
			Quality:    "WEBRip",
			Codec:      "H265",
			HDR:        "HDR10+ DV",
			Group:      "P",
		},
		fields: []string{"title", "season", "episode", "resolution", "quality", "codec", "hdr", "group"},
	},
	{
		name:  "uhd remux edition",
		title: "Nightmare Alley 2021 BW Directors Cut 2160p UHD Blu ray Remux DV HDR HEVC DTS HD MA 5 1 CiNEPHiLES",
		want: TorrentInfo{
			Title:      "Nightmare Alley",
			Year:       2021,
			Resolution: "2160p",
			Quality:    "REMUX",
			Codec:      "H265",
			HDR:        "DV HDR",
			Audio:      "DTS-HD MA 5.1",
			Group:      "CiNEPHiLES",
			Edition:    "Black and White Director's Cut",
		},
		fields: []string{"title", "year", "resolution", "quality", "codec", "hdr", "audio", "group", "edition"},
	},
	{
		name:  "hdrip xvid ac3",
		title: "On the Brain 2018 HDRip XviD AC3-EVO",
		want: TorrentInfo{
			Title:   "On the Brain",
			Year:    2018,
			Quality: "HDRip",
			Codec:   "XViD",
			Audio:   "AC3",
			Group:   "EVO",
		},
		fields: []string{"title", "year", "quality", "codec", "audio", "group"},
	},
	{
		name:  "episode part",
		title: "You and I Are Polar Opposites S01E10 Class Trip Part 1 1080p CR WEB-DL DUAL AAC2.0 H 264-VARYG (Seihantai na Kimi to Boku, Dual-Audio, Multi-Subs)",
		want: TorrentInfo{
			Title:      "You and I Are Polar Opposites",
			Season:     1,
			Episode:    10,
			Part:       1,
			Resolution: "1080p",
			Quality:    "WEB-DL",
			Codec:      "H264",
			Audio:      "AAC2.0",
			Source:     "CR",
			Group:      "VARYG",
			Edition:    "Dual Audio Multi Subs",
		},
		fields: []string{"title", "season", "episode", "part", "resolution", "quality", "codec", "audio", "source", "group", "edition"},
	},
	{
		name:  "remux truehd atmos",
		title: "[hchcsen] Gurren Lagann the Movie: The Lights in the Sky Are Stars (2009) v2 (BD Remux 1080p x264 8-bit TrueHD Atmos) [Dual Audio] | Tengen Toppa Gurren Lagann: Lagann-hen",
		want: TorrentInfo{
			Title:      "Gurren Lagann the Movie: The Lights in the Sky Are Stars",
			Year:       2009,
			Resolution: "1080p",
			Quality:    "REMUX",
			Codec:      "x264",
			Audio:      "TrueHD Atmos",
			Website:    "hchcsen",
			BitDepth:   "8-bit",
			Edition:    "Dual Audio",
		},
		fields: []string{"title", "year", "resolution", "quality", "codec", "audio", "website", "bit_depth", "edition"},
	},
	{
		name:  "hdtv h264 group",
		title: "Roman Empire By Train With Alice Roberts S01E06 720p HDTV H264-JFF",
		want: TorrentInfo{
			Title:      "Roman Empire By Train With Alice Roberts",
			Season:     1,
			Episode:    6,
			Resolution: "720p",
			Quality:    "HDTV",
			Codec:      "H264",
			Group:      "JFF",
		},
		fields: []string{"title", "season", "episode", "resolution", "quality", "codec", "group"},
	},
	{
		name:  "nf web dl",
		title: "Peaky Blinders S03E01 Episode 1 1080p NF WEB-DL DDP5 1 H 264-Kitsune",
		want: TorrentInfo{
			Title:      "Peaky Blinders",
			Season:     3,
			Episode:    1,
			Resolution: "1080p",
			Quality:    "WEB-DL",
			Codec:      "H264",
			Audio:      "DDP5.1",
			Source:     "NF",
			Group:      "Kitsune",
		},
		fields: []string{"title", "season", "episode", "resolution", "quality", "codec", "audio", "source", "group"},
	},
	{
		name:  "season only multilingual",
		title: "Young.Sherlock.2026.S01.720p.ITA-ENG.MULTI.WEBRip.x265.AAC-V3SP4EV3R",
		want: TorrentInfo{
			Title:      "Young Sherlock",
			Season:     1,
			Year:       2026,
			Resolution: "720p",
			Quality:    "WEBRip",
			Codec:      "x265",
			Audio:      "AAC",
			Group:      "V3SP4EV3R",
		},
		fields: []string{"title", "season", "year", "resolution", "quality", "codec", "audio", "group"},
	},
	{
		name:   "metadata bracket multisubs apostrophe",
		title:  "Widow's Bay (2026) s01e04 [Mkv - 1080p H264 - MultiLang Aac 2 0 - MultiSubs]",
		want:   TorrentInfo{},
		fields: []string{"group"},
	},
	{
		name:   "metadata bracket multisubs collapsed punctuation",
		title:  "Widow s Bay (2026) s01e04 [Mkv  1080p H264  MultiLang Aac 2 0  MultiSubs]",
		want:   TorrentInfo{},
		fields: []string{"group"},
	},
	{
		name:  "real bracket group after metadata",
		title: "WidowS Bay S01e04 [1080p Ita Eng Spa HEVC10 SubS] byMe7alh [MIRCrew]",
		want: TorrentInfo{
			Group: "MIRCrew",
		},
		fields: []string{"group"},
	},
	{
		name:  "real dash group after metadata",
		title: "Widows Bay S01E04 Beach Reads 2160p ATVP WEB-DL DD 5 1 Atmos H 265-playWEB",
		want: TorrentInfo{
			Group: "playWEB",
		},
		fields: []string{"group"},
	},
}

func TestFixtureGoldenCases(t *testing.T) {
	t.Parallel()

	fixtureTitles := loadFixtureTitleSet(t)
	for _, tc := range fixtureGoldenCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if !fixtureTitles[tc.title] {
				t.Fatalf("golden title is not present in testdata/jackett: %q", tc.title)
			}
			got, err := Parse(tc.title)
			if err != nil {
				t.Fatalf("Parse(%q) error = %v", tc.title, err)
			}
			assertTorrentInfoFields(t, tc.title, *got, tc.want, tc.fields)
		})
	}
}

func loadFixtureTitleSet(t *testing.T) map[string]bool {
	t.Helper()

	files, err := filepath.Glob(filepath.Join("testdata", "jackett", "*.json"))
	if err != nil {
		t.Fatalf("glob jackett fixtures: %v", err)
	}

	titles := make(map[string]bool)
	for _, path := range files {
		payload, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read fixture %s: %v", path, err)
		}
		var fixture jackettFixture
		if err := json.Unmarshal(payload, &fixture); err != nil {
			t.Fatalf("decode fixture %s: %v", path, err)
		}
		for _, result := range fixture.Results {
			titles[result.Title] = true
		}
	}
	return titles
}

func assertTorrentInfoFields(t *testing.T, title string, got TorrentInfo, want TorrentInfo, fields []string) {
	t.Helper()

	for _, field := range fields {
		switch field {
		case "title":
			assertEqual(t, title, field, got.Title, want.Title)
		case "season":
			assertEqual(t, title, field, got.Season, want.Season)
		case "episode":
			assertEqual(t, title, field, got.Episode, want.Episode)
		case "episode_end":
			assertEqual(t, title, field, got.EpisodeEnd, want.EpisodeEnd)
		case "part":
			assertEqual(t, title, field, got.Part, want.Part)
		case "year":
			assertEqual(t, title, field, got.Year, want.Year)
		case "resolution":
			assertEqual(t, title, field, got.Resolution, want.Resolution)
		case "quality":
			assertEqual(t, title, field, got.Quality, want.Quality)
		case "codec":
			assertEqual(t, title, field, got.Codec, want.Codec)
		case "hdr":
			assertEqual(t, title, field, got.HDR, want.HDR)
		case "audio":
			assertEqual(t, title, field, got.Audio, want.Audio)
		case "source":
			assertEqual(t, title, field, got.Source, want.Source)
		case "group":
			assertEqual(t, title, field, got.Group, want.Group)
		case "website":
			assertEqual(t, title, field, got.Website, want.Website)
		case "bit_depth":
			assertEqual(t, title, field, got.BitDepth, want.BitDepth)
		case "edition":
			assertEqual(t, title, field, got.Edition, want.Edition)
		case "complete":
			assertEqual(t, title, field, got.Complete, want.Complete)
		default:
			t.Fatalf("unsupported golden field %q", field)
		}
	}
}

func assertEqual[T comparable](t *testing.T, title string, field string, got T, want T) {
	t.Helper()
	if got != want {
		t.Fatalf("Parse(%q) field %s = %#v, want %#v", title, field, got, want)
	}
}
