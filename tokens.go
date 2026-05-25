package torrentname

import (
	"regexp"
	"sort"
	"strings"
)

var (
	releaseStartPattern = regexp.MustCompile(`(?i)(?:` +
		`\bS[0-9]{1,2}E[0-9]{1,3}(?:[ .-]*E?[0-9]{1,3})?\b` +
		`|\bSeason[ .-]+[0-9]{1,2}\b` +
		`|\b(?:19[0-9]{2}|20[0-9]{2})\b` +
		`|\b(?:` + tokenPatternAlternates(resolutionCatalog, broadResolutionAliasContextPattern, `[0-9]{3,4}p`) + `)\b` +
		`|\b(?:` + strings.Join(tokenPatterns(qualityCatalog), "|") + `)\b` +
		`|\b(?:` + strings.Join(tokenPatterns(codecCatalog), "|") + `)\b` +
		`|\b(?:` + strings.Join(tokenPatterns(hdrCatalog), "|") + `)` +
		`|\b(?:` + strings.Join(tokenPatterns(audioCatalog), "|") + `)\b` +
		`|\b(?:8|10|12|16|24)[ .-]?bits?\b` +
		`)`)
	sourceTokenPattern = compileSourcePattern(sourceCatalog)
)

type tokenMatch struct {
	start int
	end   int
	raw   string
}

func firstReleaseTokenPosition(value string) int {
	if match := releaseStartPattern.FindStringIndex(value); match != nil {
		return match[0]
	}
	return -1
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
