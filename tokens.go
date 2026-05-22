package torrentname

import (
	"regexp"
	"sort"
)

var (
	releaseStartPatterns = []*regexp.Regexp{
		seasonEpisodeRangePattern,
		seasonReleaseTokenPattern,
		yearTokenPattern,
		resolutionTokenPattern,
		qualityTokenPattern,
		codecTokenPattern,
		hdrTokenPattern,
		audioTokenPattern,
		bitDepthPattern,
	}
	sourceTokenPattern = regexp.MustCompile(`(?i)(?:^|[^A-Za-z0-9])((?:ABEMA|AMZN|ATV\+?|ATVP|AUBC|BILI|B[ ._-]?CORE|CBC|CPNG|CR|CRAVE|CRAV|CRiT|DSNP|FOD|friDay|Hami|HBO[ ._-]?MAX|HBOM|HMAX|HBO|HTSR|HULU|iQIY|ITVX|iTUNES|iT|KCW|KKTV|LINE[ ._-]?TV|MAX|MY5|MyTVSuper|NF|NOW|OVID|PCOK|PLAY|PMTP|ROKU|STAN|STAR\+?|STRP|TVING|TVER|U[ ._-]?NEXT|Viki|VIU|VRV|WAVVE|WETV|YOUKU|iP))(?:$|[^A-Za-z0-9])`)
)

type tokenMatch struct {
	start int
	end   int
	raw   string
}

func firstReleaseTokenPosition(value string) int {
	first := -1
	for _, pattern := range releaseStartPatterns {
		if match := pattern.FindStringIndex(value); match != nil && (first < 0 || match[0] < first) {
			first = match[0]
		}
	}
	return first
}

func orderedNormalizedTokens(value string, pattern *regexp.Regexp, normalize func(string) string) []string {
	matches := make([]tokenMatch, 0)
	for _, match := range pattern.FindAllStringIndex(value, -1) {
		matches = append(matches, tokenMatch{start: match[0], end: match[1], raw: value[match[0]:match[1]]})
	}
	sortTokenMatches(matches)

	out := make([]string, 0, len(matches))
	seen := map[string]struct{}{}
	for _, match := range matches {
		normalized := normalize(match.raw)
		if normalized == "" {
			continue
		}
		if _, ok := seen[normalized]; ok {
			continue
		}
		seen[normalized] = struct{}{}
		out = append(out, normalized)
	}
	return out
}

func sourceMatches(value string) []tokenMatch {
	matches := make([]tokenMatch, 0)
	searchStart := 0
	for searchStart < len(value) {
		match := sourceTokenPattern.FindStringSubmatchIndex(value[searchStart:])
		if match == nil {
			break
		}
		if len(match) < 4 || match[2] < 0 || match[3] < 0 {
			searchStart += match[1]
			continue
		}
		start := searchStart + match[2]
		end := searchStart + match[3]
		matches = append(matches, tokenMatch{start: start, end: end, raw: value[start:end]})
		if end <= searchStart {
			searchStart += match[1]
			continue
		}
		searchStart = end
	}
	sortTokenMatches(matches)
	return matches
}

func sortTokenMatches(matches []tokenMatch) {
	sort.SliceStable(matches, func(i, j int) bool {
		return matches[i].start < matches[j].start
	})
}
