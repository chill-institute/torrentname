package torrentname

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type pattern struct {
	name  string
	last  bool
	re    *regexp.Regexp
	apply func(*TorrentInfo, string)
}

var patterns = []pattern{
	{name: "season", re: regexp.MustCompile(`(?i)(s?([0-9]{1,2}))[ex]`), apply: setInt(func(t *TorrentInfo, value int) { t.Season = value })},
	{name: "episode", re: regexp.MustCompile(`(?i)([ex]([0-9]{2})(?:[^0-9]|$))`), apply: setInt(func(t *TorrentInfo, value int) { t.Episode = value })},
	{name: "episode", re: regexp.MustCompile(`(-\s+([0-9]{1,})(?:[^0-9]|$))`), apply: setInt(func(t *TorrentInfo, value int) { t.Episode = value })},
	{name: "year", last: true, re: regexp.MustCompile(`\b(((?:19[0-9]|20[0-9])[0-9]))\b`), apply: setInt(func(t *TorrentInfo, value int) { t.Year = value })},
	{name: "resolution", re: regexp.MustCompile(`\b(([0-9]{3,4}p))\b`), apply: func(t *TorrentInfo, value string) { t.Resolution = value }},
	{name: "quality", re: regexp.MustCompile(`(?i)\b(((?:PPV\.)?[HP]DTV|(?:HD)?CAM|B[DR]Rip|(?:HD-?)?TS|(?:PPV )?WEB-?DL(?: DVDRip)?|HDRip|DVDRip|DVDRIP|CamRip|W[EB]BRip|BluRay|DvDScr|telesync))\b`), apply: func(t *TorrentInfo, value string) { t.Quality = value }},
	{name: "codec", re: regexp.MustCompile(`(?i)\b((xvid|[hx]\.?26[45]))\b`), apply: func(t *TorrentInfo, value string) { t.Codec = value }},
	{name: "audio", re: regexp.MustCompile(`(?i)\b((MP3|DD5\.?1|Dual[\- ]Audio|LiNE|DTS|AAC[.-]LC|AAC(?:\.?2\.0)?|AC3(?:\.5\.1)?))\b`), apply: func(t *TorrentInfo, value string) { t.Audio = value }},
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
	{name: "widescreen", re: regexp.MustCompile(`(?i)\b((WS))\b`), apply: func(t *TorrentInfo, _ string) { t.Widescreen = true }},
	{name: "unrated", re: regexp.MustCompile(`(?i)\b((UNRATED))\b`), apply: func(t *TorrentInfo, _ string) { t.Unrated = true }},
	{name: "threeD", re: regexp.MustCompile(`(?i)\b((3D))\b`), apply: func(t *TorrentInfo, _ string) { t.ThreeD = true }},
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

func init() {
	for _, pat := range patterns {
		if pat.re.NumSubexp() != 2 {
			fmt.Printf("Pattern %q does not have enough capture groups. want 2, got %d\n", pat.name, pat.re.NumSubexp())
			os.Exit(1)
		}
	}
}
