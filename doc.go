// Package torrentname parses torrent-style release names into structured
// metadata without consulting external catalogs or services.
//
// The parser is deterministic and best-effort. Unknown fields remain at their
// Go zero value, so callers can safely use Parse on noisy search results and
// decide which fields are strong enough for their own matching workflow.
package torrentname
