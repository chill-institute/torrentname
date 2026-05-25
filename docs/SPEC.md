# Parser Spec

`torrentname` parses torrent-style release names into `TorrentInfo`. The parser is deterministic, best-effort, and dependency-free at runtime.

## Contract

- `Parse` accepts a single filename or release title string and returns a populated `TorrentInfo`.
- Parse failures are reserved for future hard errors; current parsing is best-effort and should not reject unusual release names.
- The parser normalizes common release tokens while preserving human-readable title text.
- Missing or unrecognized fields stay at their Go zero value.
- The parser does not fetch metadata, validate titles against external catalogs, or infer facts that are not present in the release name.

## Supported Fields

| Field | Supported inputs | Normalization |
| --- | --- | --- |
| `Title` | Text before release metadata such as season, year, resolution, quality, or source tokens | Separators collapse to spaces; wrapper website/group tags are removed when recognized |
| `Season` | `S01`, `S01E02`, `1x02`, `Season 1`, and complete-season markers | Integer without leading zero |
| `Episode` | `S01E02`, `1x02`, and anime-style ` - 10 ` forms | Integer without leading zero |
| `EpisodeEnd` | Episode ranges such as `S01E01-E03` and `S01E01 03 of 10` | Integer without leading zero |
| `Part` | `Part 1`, `Part One`, and roman numerals through `X` | Integer |
| `Year` | Years from `1900` through `2099` | Integer |
| `Resolution` | `SD`, `HD`, `FHD`, `480p`, `720p`, `1080p`, `2160p`, uppercase `P` variants, `4K`, and `UHD` | Common aliases normalize to numeric `p` forms; numeric tokens canonicalize to lowercase; explicit numeric tokens take precedence when both appear |
| `Quality` | WEB, WEB-DL, WEBRip, HDTV, PDTV, BluRay, REMUX, HDRip, DVDRip, BRRip, BDRip, DvDScr, CAM/TS/TC-style tokens | Common web, BluRay, remux, telecine, telesync, and cam variants collapse to canonical spellings such as `WEB-DL`, `REMUX`, `TC`, `TS`, and `CAM` |
| `Codec` | x264, x265, H.264, H.265, AVC, HEVC, AV1, XViD | Common aliases collapse to `x264`, `x265`, `H264`, `H265`, `AV1`, or `XViD` |
| `HDR` | HDR, HDR10, HDR10+, DV/DoVi/Dolby Vision, HLG | Ordered unique tokens such as `HDR10+ DV` |
| `Audio` | MP3, AAC/AAC-LC, AC3, EAC3, DDP, DDPA, DD, DTS, DTS X, DTS-HD MA/HRA, TrueHD, Atmos, FLAC, LiNE, PCM/LPCM, Opus, dual-audio markers, and channel tokens including `2CH`, `6CH`, `8CH`, `2.0`, `5.1`, and `7.1` | Common channel variants collapse to forms such as `DDP5.1`, `DDP Atmos 5.1`, `DD Atmos 5.1`, `EAC3 5.1`, `EAC3 Atmos 5.1`, `PCM 2.0`, `TrueHD Atmos 7.1`; real audio tokens take precedence over dual-audio edition markers |
| `Source` | Known source tags after release metadata, including ABEMA, AMZN, ATVP, BILI, BCORE, CR, DSNP, HULU, NF, PCOK, PMTP, ROKU, STAN, HMAX/HBO/MAX, iT/iP, and other current streaming-service tags | Canonical uppercase except stylized tags such as `iT`, `iP`, and `CRiT` |
| `Group` | Dash suffixes, bracket suffixes, and advanced release trailing group names | Wrapper characters and spaces are stripped; metadata-looking tokens are ignored |
| `Website` | Leading bracket tags such as `[Source]` | Trimmed bracket content |
| `Language` | Current explicit language pairs such as `rus.eng`, `ita.eng`, and `ENG.LAT`, strong single post-release markers such as `JAPANESE` and `KOREAN`, plus clusters such as `Eng.Rus.Multi-Subs`, `VOSTFR`, `VFF`, and `MultiLang` | Preserved explicit pairs or uppercase normalized clusters |
| `BitDepth` | `8-bit`, `10bit`, `12 bit`, `16Bit`, `24bit` | `N-bit` |
| `Edition` | Director's Cut, DC, hybrid, theatrical cut, special edition, open matte, B&W, dubbed, dual audio, `2Audios`, `M SUB`, multi subs | Canonical descriptive labels |
| `Size` | `MB` and `GB` size tokens | Preserved token |
| Flags | Extended, hardcoded (`HC`, `HC-SUB`, `KORSUB`), proper, repack, remastered (`Remastered`, `RM4K`), widescreen, unrated, 3D, IMAX, complete | Boolean |
| `Container` | MKV, AVI, MP4 tokens | Preserved token |
| `Sbs` | SBS and Half-SBS | Preserved token |

## Accuracy Loop

Accuracy is guarded by three layers:

- unit tests for curated parser examples
- golden real-world cases selected from `testdata/jackett`
- a corpus metrics command that reports field coverage across every committed Jackett fixture title

Use:

```bash
mise run test
mise run corpus:metrics
```

Refresh the Jackett corpus only when intentionally updating real-world samples:

```bash
JACKETT_API_KEY=... mise run fixtures:jackett
```

Fixture refreshes often change tracker timestamps and result ordering. Commit only fixture changes that improve the representative corpus.

## Performance Loop

Performance is tracked with parser microbenchmarks and the 1k-row batch benchmark.

Use this workflow around parser changes:

```bash
BENCH_OUT=tmp/bench/baseline.txt mise run bench:record
# edit parser
BENCH_OUT=tmp/bench/current.txt mise run bench:record
mise run bench:compare
```

Review `ns/op`, `B/op`, `allocs/op`, `rows/sec`, and `ms/1krows`. Treat allocation increases on hot parser paths as regressions unless the accuracy gain is worth the cost.
