package torrentname

import (
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
	{name: "season", re: regexp.MustCompile(`(?i)\b((?:s?([0-9]{1,2}))[ex][0-9])`), apply: setInt(func(t *TorrentInfo, value int) { t.Season = value })},
	{name: "season", re: regexp.MustCompile(`(?i)\b((?:s([0-9]{1,2}))(?:\b|[^a-z0-9]))`), apply: setInt(func(t *TorrentInfo, value int) { t.Season = value })},
	{name: "season", re: regexp.MustCompile(`(?i)\b((?:Season)[ .-]+([0-9]{1,2}))\b`), apply: setInt(func(t *TorrentInfo, value int) { t.Season = value })},
	{name: "episode", re: regexp.MustCompile(`(?i)([ex]([0-9]{2})(?:[^0-9]|$))`), apply: setInt(func(t *TorrentInfo, value int) { t.Episode = value })},
	{name: "episode", re: regexp.MustCompile(`(-\s+([0-9]{1,})(?:[^0-9]|$))`), apply: setInt(func(t *TorrentInfo, value int) { t.Episode = value })},
	{name: "year", last: true, re: regexp.MustCompile(`\b(((?:19[0-9]|20[0-9])[0-9]))\b`), apply: setInt(func(t *TorrentInfo, value int) { t.Year = value })},
	{name: "resolution", re: compileCapturedTokenPattern(resolutionCatalog, broadResolutionAliasContextPattern, `[0-9]{3,4}p`), apply: func(t *TorrentInfo, value string) { t.Resolution = normalizeResolution(value) }},
	{name: "quality", re: compileCapturedTokenPattern(qualityCatalog), apply: func(t *TorrentInfo, value string) { t.Quality = normalizeQuality(value) }},
	{name: "codec", re: compileCapturedTokenPattern(codecCatalog), apply: func(t *TorrentInfo, value string) { t.Codec = normalizeCodec(value) }},
	{name: "audio", re: compileCapturedTokenPattern(audioCatalog), apply: func(t *TorrentInfo, value string) { t.Audio = normalizeAudio(value) }},
	{name: "region", re: regexp.MustCompile(`(?i)\b(R([0-9]))\b`), apply: func(t *TorrentInfo, value string) { t.Region = value }},
	{name: "size", re: regexp.MustCompile(`(?i)\b((\d+(?:\.\d+)?(?:GB|MB)))\b`), apply: func(t *TorrentInfo, value string) { t.Size = value }},
	{name: "website", re: regexp.MustCompile(`^(\[ ?([^\]]+?) ?\])`), apply: func(t *TorrentInfo, value string) { t.Website = value }},
	{name: "language", re: regexp.MustCompile(`(?i)\b((rus\.eng|ita\.eng))\b`), apply: func(t *TorrentInfo, value string) { t.Language = value }},
	{name: "sbs", re: regexp.MustCompile(`(?i)\b(((?:Half-)?SBS))\b`), apply: func(t *TorrentInfo, value string) { t.Sbs = value }},
	{name: "container", re: compileCapturedTokenPattern(containerCatalog), apply: func(t *TorrentInfo, value string) { t.Container = value }},
	{name: "group", re: regexp.MustCompile(`(?i)(- ?([A-Za-z0-9\[{][A-Za-z0-9\[\]{}_=+-]*))$`), apply: func(t *TorrentInfo, value string) { t.Group = value }},
	{name: "extended", re: compileCapturedTokenPattern(extendedCatalog), apply: func(t *TorrentInfo, _ string) { t.Extended = true }},
	{name: "hardcoded", re: compileCapturedTokenPattern(hardcodedCatalog), apply: func(t *TorrentInfo, _ string) { t.Hardcoded = true }},
	{name: "proper", re: compileCapturedTokenPattern(properCatalog), apply: func(t *TorrentInfo, _ string) { t.Proper = true }},
	{name: "repack", re: compileCapturedTokenPattern(repackCatalog), apply: func(t *TorrentInfo, _ string) { t.Repack = true }},
	{name: "remastered", re: compileCapturedTokenPattern(remasteredCatalog), apply: func(t *TorrentInfo, _ string) { t.Remastered = true }},
	{name: "widescreen", re: compileCapturedTokenPattern(widescreenCatalog), apply: func(t *TorrentInfo, _ string) { t.Widescreen = true }},
	{name: "unrated", re: compileCapturedTokenPattern(unratedCatalog), apply: func(t *TorrentInfo, _ string) { t.Unrated = true }},
	{name: "threeD", re: compileCapturedTokenPattern(threeDCatalog), apply: func(t *TorrentInfo, _ string) { t.ThreeD = true }},
	{name: "imax", re: compileCapturedTokenPattern(imaxCatalog), apply: func(t *TorrentInfo, _ string) { t.IMAX = true }},
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
	if normalized, ok := qualityLookup[compactKey(value)]; ok {
		if !canonicalizesQuality(normalized) {
			return value
		}
		return normalized
	}
	return value
}

func canonicalizesQuality(value string) bool {
	switch value {
	case "WEB-DL", "WEBRip", "WEB", "BluRay", "REMUX", "TC", "TS", "CAM":
		return true
	default:
		return false
	}
}

func normalizeCodec(value string) string {
	if normalized, ok := codecLookup[compactKey(value)]; ok {
		return normalized
	}
	return value
}

func normalizeAudio(value string) string {
	collapsed := strings.ToLower(audioCompactReplacer.Replace(value))
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
	case "ddplus":
		return "DDPlus"
	case "dtshd", "dtshdma":
		if strings.Contains(collapsed, "ma") {
			return "DTS-HD MA"
		}
		return "DTS-HD"
	case "dtsx":
		return "DTS X"
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
