package torrentname

import "testing"

func TestCatalogAliasNormalization(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		tokens    []aliasToken
		normalize func(string) string
	}{
		{name: "hdr", tokens: hdrCatalog, normalize: normalizeHDR},
		{name: "codec", tokens: codecCatalog, normalize: normalizeCodec},
		{name: "language", tokens: languageCatalog, normalize: normalizeLanguage},
		{name: "edition", tokens: editionCatalog, normalize: normalizeEdition},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			for _, token := range tc.tokens {
				for _, alias := range append([]string{token.canonical}, token.aliases...) {
					if got := tc.normalize(alias); got != token.canonical {
						t.Fatalf("%s normalize(%q) = %q, want %q", tc.name, alias, got, token.canonical)
					}
				}
			}
		})
	}
}

func TestSourceCatalogNormalizesAndMatches(t *testing.T) {
	t.Parallel()

	for _, token := range sourceCatalog {
		token := token
		t.Run(token.canonical, func(t *testing.T) {
			t.Parallel()

			for _, alias := range append([]string{token.canonical}, token.aliases...) {
				if got := normalizeSource(alias); got != token.canonical {
					t.Fatalf("normalizeSource(%q) = %q, want %q", alias, got, token.canonical)
				}

				title := "Sample Title 2026 1080p " + alias + " WEB-DL H264-GRP"
				info, err := Parse(title)
				if err != nil {
					t.Fatalf("Parse(%q) error: %v", title, err)
				}
				if info.Source != token.canonical {
					t.Fatalf("Parse(%q).Source = %q, want %q", title, info.Source, token.canonical)
				}
			}
		})
	}
}
