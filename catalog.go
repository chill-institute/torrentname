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
	return regexp.MustCompile(`(?i)\b(?:` + strings.Join(tokenPatterns(tokens), "|") + `)\b`)
}

func compileLooseEndTokenPattern(tokens []aliasToken) *regexp.Regexp {
	return regexp.MustCompile(`(?i)\b(?:` + strings.Join(tokenPatterns(tokens), "|") + `)`)
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
