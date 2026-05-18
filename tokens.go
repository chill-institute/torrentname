package torrentname

import (
	"regexp"
	"sort"
)

var (
	releaseStartPatterns = []*regexp.Regexp{
		seasonEpisodeRangePattern,
		yearTokenPattern,
		resolutionTokenPattern,
		qualityTokenPattern,
		codecTokenPattern,
		hdrTokenPattern,
		audioTokenPattern,
		bitDepthPattern,
	}
	sourceTokenPattern = regexp.MustCompile(`(?i)(?:^|[^A-Za-z0-9])((?:AMZN|NF|DSNP|HULU|CR|ATVP|PCOK|HMAX|HBO|MAX|iT))(?:$|[^A-Za-z0-9])`)
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
	rawMatches := sourceTokenPattern.FindAllStringSubmatchIndex(value, -1)
	matches := make([]tokenMatch, 0, len(rawMatches))
	for _, match := range rawMatches {
		if len(match) < 4 || match[2] < 0 || match[3] < 0 {
			continue
		}
		matches = append(matches, tokenMatch{start: match[2], end: match[3], raw: value[match[2]:match[3]]})
	}
	sortTokenMatches(matches)
	return matches
}

func sortTokenMatches(matches []tokenMatch) {
	sort.SliceStable(matches, func(i, j int) bool {
		return matches[i].start < matches[j].start
	})
}
