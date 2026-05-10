package torrentname

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type pattern struct {
	name  string
	last  bool
	re    *regexp.Regexp
	apply func(*TorrentInfo, string)
}

var patterns = []pattern{
	{name: "season", re: regexp.MustCompile(`(?i)(s?([0-9]{1,2}))[ex]`), apply: setInt(func(t *TorrentInfo, value int) { t.Season = value })},
	{name: "season", re: regexp.MustCompile(`(?i)\b((?:s([0-9]{1,2}))(?:\b|[^a-z0-9]))`), apply: setInt(func(t *TorrentInfo, value int) { t.Season = value })},
	{name: "episode", re: regexp.MustCompile(`(?i)([ex]([0-9]{2})(?:[^0-9]|$))`), apply: setInt(func(t *TorrentInfo, value int) { t.Episode = value })},
	{name: "episode", re: regexp.MustCompile(`(-\s+([0-9]{1,})(?:[^0-9]|$))`), apply: setInt(func(t *TorrentInfo, value int) { t.Episode = value })},
	{name: "year", last: true, re: regexp.MustCompile(`\b(((?:19[0-9]|20[0-9])[0-9]))\b`), apply: setInt(func(t *TorrentInfo, value int) { t.Year = value })},
	{name: "resolution", re: regexp.MustCompile(`\b(([0-9]{3,4}p))\b`), apply: func(t *TorrentInfo, value string) { t.Resolution = value }},
	{name: "quality", re: regexp.MustCompile(`(?i)\b(((?:PPV\.)?[HP]DTV|(?:HD)?CAM|B[DR]Rip|(?:HD-?)?TS|(?:PPV )?WEB[ .-]?DL(?: DVDRip)?|HDRip|DVDRip|DVDRIP|CamRip|WEB[ .-]?Rip|WBRip|Blu[ .-]?Ray(?:[ .-]?Remux)?|BDRemux|REMUX|DvDScr|telesync))\b`), apply: func(t *TorrentInfo, value string) { t.Quality = normalizeQuality(value) }},
	{name: "codec", re: regexp.MustCompile(`(?i)\b((xvid|x[ .]?26[45]|[h]\.? ?26[45]|HEVC|AVC|AV1))\b`), apply: func(t *TorrentInfo, value string) { t.Codec = normalizeCodec(value) }},
	{name: "audio", re: regexp.MustCompile(`(?i)\b((MP3|TrueHD|Atmos|DTS[ .-]?HD(?:[ .-]?MA)?|E-?AC-?3|DDP[ .]?[257][ .]?[01]|DD[ .]?5[ .]?1|DD\+|Dual[\- ]Audio|LiNE|DTS|AAC[.-]?LC|AAC(?:\.?[257]\.?[01])?|AC3(?:\.5\.1)?|FLAC))\b`), apply: func(t *TorrentInfo, value string) { t.Audio = normalizeAudio(value) }},
	{name: "region", re: regexp.MustCompile(`(?i)\b(R([0-9]))\b`), apply: func(t *TorrentInfo, value string) { t.Region = value }},
	{name: "size", re: regexp.MustCompile(`(?i)\b((\d+(?:\.\d+)?(?:GB|MB)))\b`), apply: func(t *TorrentInfo, value string) { t.Size = value }},
	{name: "website", re: regexp.MustCompile(`^(\[ ?([^\]]+?) ?\])`), apply: func(t *TorrentInfo, value string) { t.Website = value }},
	{name: "language", re: regexp.MustCompile(`(?i)\b((rus\.eng|ita\.eng))\b`), apply: func(t *TorrentInfo, value string) { t.Language = value }},
	{name: "sbs", re: regexp.MustCompile(`(?i)\b(((?:Half-)?SBS))\b`), apply: func(t *TorrentInfo, value string) { t.Sbs = value }},
	{name: "container", re: regexp.MustCompile(`(?i)\b((MKV|AVI|MP4))\b`), apply: func(t *TorrentInfo, value string) { t.Container = value }},
	{name: "group", re: regexp.MustCompile(`(?i)(- ?([A-Za-z0-9\[{][A-Za-z0-9\[\]{}_=+-]*))$`), apply: func(t *TorrentInfo, value string) { t.Group = value }},
	{name: "extended", re: regexp.MustCompile(`(?i)\b(EXTENDED(:?.CUT)?)\b`), apply: func(t *TorrentInfo, _ string) { t.Extended = true }},
	{name: "hardcoded", re: regexp.MustCompile(`(?i)\b((HC))\b`), apply: func(t *TorrentInfo, _ string) { t.Hardcoded = true }},
	{name: "proper", re: regexp.MustCompile(`(?i)\b((PROPER))\b`), apply: func(t *TorrentInfo, _ string) { t.Proper = true }},
	{name: "repack", re: regexp.MustCompile(`(?i)\b((REPACK))\b`), apply: func(t *TorrentInfo, _ string) { t.Repack = true }},
	{name: "remastered", re: regexp.MustCompile(`(?i)\b((REMASTERED))\b`), apply: func(t *TorrentInfo, _ string) { t.Remastered = true }},
	{name: "widescreen", re: regexp.MustCompile(`(?i)\b((WS))\b`), apply: func(t *TorrentInfo, _ string) { t.Widescreen = true }},
	{name: "unrated", re: regexp.MustCompile(`(?i)\b((UNRATED))\b`), apply: func(t *TorrentInfo, _ string) { t.Unrated = true }},
	{name: "threeD", re: regexp.MustCompile(`(?i)\b((3D))\b`), apply: func(t *TorrentInfo, _ string) { t.ThreeD = true }},
	{name: "imax", re: regexp.MustCompile(`(?i)\b((IMAX))\b`), apply: func(t *TorrentInfo, _ string) { t.IMAX = true }},
	{name: "complete", re: regexp.MustCompile(`(?i)\b((COMPLETE(?:[ .-]?SEASON|[ .-]?SERIES)?|SEASON[ .-]?[0-9]{1,2}[ .-]?COMPLETE))\b`), apply: func(t *TorrentInfo, _ string) { t.Complete = true }},
}

func setInt(assign func(*TorrentInfo, int)) func(*TorrentInfo, string) {
	return func(t *TorrentInfo, raw string) {
		value, err := strconv.Atoi(raw)
		if err != nil {
			return
		}
		assign(t, value)
	}
}

func normalizeQuality(value string) string {
	collapsed := strings.ToLower(strings.NewReplacer(" ", "", "-", "", ".", "").Replace(value))
	switch collapsed {
	case "webdl", "ppvwebdl", "hdwebdl":
		return "WEB-DL"
	case "webrip", "wbrip":
		return "WEBRip"
	case "web":
		return "WEB"
	case "bluray":
		return "BluRay"
	case "blurayremux", "bdremux", "remux", "uhdblurayremux":
		return "REMUX"
	default:
		return value
	}
}

func normalizeCodec(value string) string {
	collapsed := strings.ToLower(strings.NewReplacer(" ", "", ".", "").Replace(value))
	switch collapsed {
	case "x264":
		return "x264"
	case "x265":
		return "x265"
	case "h264", "avc":
		return "H264"
	case "h265", "hevc":
		return "H265"
	case "av1":
		return "AV1"
	case "xvid":
		return "XViD"
	default:
		return value
	}
}

func normalizeAudio(value string) string {
	collapsed := strings.ToLower(strings.NewReplacer(" ", "", ".", "", "-", "").Replace(value))
	switch collapsed {
	case "ddp51":
		return "DDP5.1"
	case "ddp20":
		return "DDP2.0"
	case "dd51":
		return "DD5.1"
	case "eac3", "eac320":
		if strings.Contains(collapsed, "20") {
			return "EAC3 2.0"
		}
		return "EAC3"
	case "dd+":
		return "DD+"
	case "dtshd", "dtshdma":
		if strings.Contains(collapsed, "ma") {
			return "DTS-HD MA"
		}
		return "DTS-HD"
	case "truehd":
		return "TrueHD"
	case "atmos":
		return "Atmos"
	case "aac20":
		return "AAC2.0"
	case "dualaudio":
		if strings.Contains(value, "-") {
			return "Dual-Audio"
		}
		return "Dual Audio"
	default:
		return value
	}
}

func init() {
	for _, pat := range patterns {
		if pat.re.NumSubexp() != 2 {
			fmt.Printf("Pattern %q does not have enough capture groups. want 2, got %d\n", pat.name, pat.re.NumSubexp())
			os.Exit(1)
		}
	}
}
