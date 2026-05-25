package torrentname

import (
	"regexp"
	"strings"
)

var compactKeyReplacer = strings.NewReplacer(" ", "", ".", "", "-", "", "_", "")

type aliasToken struct {
	canonical string
	aliases   []string
	patterns  []string
}

type sourceToken struct {
	canonical       string
	aliases         []string
	patterns        []string
	requiresContext bool
	ambiguousGroup  bool
}

func buildAliasLookup(tokens []aliasToken) map[string]string {
	lookup := make(map[string]string)
	for _, token := range tokens {
		lookup[compactKey(token.canonical)] = token.canonical
		for _, alias := range token.aliases {
			lookup[compactKey(alias)] = token.canonical
		}
	}
	return lookup
}

func mergeAliasTokens(catalogs ...[]aliasToken) []aliasToken {
	var merged []aliasToken
	for _, catalog := range catalogs {
		merged = append(merged, catalog...)
	}
	return merged
}

func selectAliasTokens(tokens []aliasToken, canonicals ...string) []aliasToken {
	wanted := make(map[string]struct{}, len(canonicals))
	for _, canonical := range canonicals {
		wanted[canonical] = struct{}{}
	}

	selected := make([]aliasToken, 0, len(canonicals))
	for _, token := range tokens {
		if _, ok := wanted[token.canonical]; ok {
			selected = append(selected, token)
		}
	}
	return selected
}

func buildSourceLookup(tokens []sourceToken) map[string]sourceToken {
	lookup := make(map[string]sourceToken)
	for _, token := range tokens {
		lookup[compactKey(token.canonical)] = token
		for _, alias := range token.aliases {
			lookup[compactKey(alias)] = token
		}
	}
	return lookup
}

func compactKey(value string) string {
	return strings.ToUpper(compactKeyReplacer.Replace(strings.TrimSpace(value)))
}

func compileTokenPattern(tokens []aliasToken) *regexp.Regexp {
	return regexp.MustCompile(`(?i)\b(?:` + tokenPatternAlternates(tokens) + `)\b`)
}

func compileLooseEndTokenPattern(tokens []aliasToken) *regexp.Regexp {
	return regexp.MustCompile(`(?i)\b(?:` + tokenPatternAlternates(tokens) + `)`)
}

func compileCapturedTokenPattern(tokens []aliasToken, extraPatterns ...string) *regexp.Regexp {
	return regexp.MustCompile(`(?i)\b((` + tokenPatternAlternates(tokens, extraPatterns...) + `))\b`)
}

func compileSourcePattern(tokens []sourceToken) *regexp.Regexp {
	patterns := make([]string, 0, len(tokens))
	for _, token := range tokens {
		patterns = append(patterns, token.patterns...)
	}
	return regexp.MustCompile(`(?i)(?:^|[^A-Za-z0-9])((?:` + strings.Join(patterns, "|") + `))(?:$|[^A-Za-z0-9])`)
}

func tokenPatterns(tokens []aliasToken) []string {
	patterns := make([]string, 0, len(tokens))
	for _, token := range tokens {
		patterns = append(patterns, token.patterns...)
	}
	return patterns
}

func sourceTokenPatterns(tokens []sourceToken) []string {
	patterns := make([]string, 0, len(tokens))
	for _, token := range tokens {
		patterns = append(patterns, token.patterns...)
	}
	return patterns
}

func tokenPatternAlternates(tokens []aliasToken, extraPatterns ...string) string {
	patterns := tokenPatterns(tokens)
	patterns = append(patterns, extraPatterns...)
	return strings.Join(patterns, "|")
}

func sourceTokenPatternAlternates(tokens []sourceToken) string {
	return strings.Join(sourceTokenPatterns(tokens), "|")
}
