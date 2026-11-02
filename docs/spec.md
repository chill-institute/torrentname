# Parser contract

`Parse` extracts only what a release name supports, leaves unknown fields at
zero values, and never consults external catalogs. Parsing is deterministic
and dependency-free at runtime.

- [TorrentInfo and Parse](../parser.go) define the output fields and parser.
- [catalog_tokens.go](../catalog_tokens.go) defines token aliases and canonical values.
- [normalize.go](../normalize.go) defines normalization.
- [parser_test.go](../parser_test.go) and [testdata](../testdata/) define exact
  input/output examples and boundary cases.

Keep heuristics general to release names. Corpus field-presence floors measure
how often fields are found, not whether the values are accurate.
