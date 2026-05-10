package torrentname

import (
	"html"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	seasonEpisodeRangePattern = regexp.MustCompile(`(?i)\bS([0-9]{1,2})E([0-9]{1,3})(?:[ .-]*E?([0-9]{1,3}))?\b`)
	compactRangePattern       = regexp.MustCompile(`(?i)\bS([0-9]{1,2})E([0-9]{1,3})[ .-]+([0-9]{1,3})(?:[ .-]+of[ .-]+[0-9]{1,3})?\b`)
	yearTokenPattern          = regexp.MustCompile(`\b(?:19[0-9]{2}|20[0-9]{2})\b`)
	resolutionTokenPattern    = regexp.MustCompile(`(?i)\b(?:[0-9]{3,4}p|4K)\b`)
	qualityTokenPattern       = regexp.MustCompile(`(?i)\b(?:UHD[ .-]+Blu[ .-]?Ray[ .-]+Remux|Blu[ .-]?Ray[ .-]+Remux|BDRemux|REMUX|WEB[ .-]?DL|WEB[ .-]?Rip|Blu[ .-]?Ray|HDRip|DVDRip|BRRip|BDRip|HDTV|PDTV|WEB)\b`)
	codecTokenPattern         = regexp.MustCompile(`(?i)\b(?:xvid|x[ .]?26[45]|h[ .]?26[45]|HEVC|AVC|AV1)\b`)
	hdrTokenPattern           = regexp.MustCompile(`(?i)\b(?:Dolby[ .-]+Vision|HDR10Plus|HDR10P|HDR10\+|HDR10|DoVi|DV|HLG|HDR)\b`)
	audioTokenPattern         = regexp.MustCompile(`(?i)\b(?:TrueHD(?:[ .-]+Atmos)?(?:[ .-]+[257][ .][01])?|DTS[ .-]?HD(?:[ .-]?MA)?(?:[ .-]+[257][ .][01])?|E-?AC-?3(?:[ .-]*[257][ .][01])?|DDP(?:[ .-]*Atmos)?[ .-]*[257][ .][01](?:[ .-]+Atmos)?|DD\+|AAC[ .-]*[257][ .][01]|AAC|AC3(?:[ .-]*[257][ .][01])?|FLAC(?:[ .-]*[257][ .][01])?|Atmos)\b`)
	sourceTokenPattern        = regexp.MustCompile(`(?i)(?:^|[^A-Za-z0-9])((?:AMZN|NF|DSNP|HULU|CR|ATVP|PCOK|HMAX|HBO|MAX|iT))(?:$|[^A-Za-z0-9])`)
	bitDepthPattern           = regexp.MustCompile(`(?i)\b(8|10|12)[ .-]?bit\b`)
	partPattern               = regexp.MustCompile(`(?i)\bPart[ .-]+([0-9]+|One|Two|Three|Four|Five|Six|Seven|Eight|Nine|Ten|I|II|III|IV|V|VI|VII|VIII|IX|X)\b`)
	completePattern           = regexp.MustCompile(`(?i)\b(?:Complete(?:[ .-]+(?:Season|Series))?|Season[ .-]+[0-9]{1,2}[ .-]+Complete|Seasons[ .-]+[0-9]{1,2}[ .-]+to[ .-]+[0-9]{1,2}[ .-]+Complete|S[0-9]{1,2}[ .-]*Complete)\b`)
	editionTokenPattern       = regexp.MustCompile(`(?i)\b(?:Director'?s[ .-]+Cut|Directors[ .-]+Cut|DC|Hybrid|B&W|BW|Black[ .-]+and[ .-]+White|DUBBED|DUAL|Dual[ .-]+Audio|Multi[ .-]*Subs?|MultiSub)\b`)
	dashGroupPattern          = regexp.MustCompile(`(?i)-\s*([A-Za-z0-9][A-Za-z0-9._=+\[\]{}-]{0,31})(?:\s*(?:\(|\[|$))`)
	bracketGroupPattern       = regexp.MustCompile(`(?i)\[([A-Za-z][A-Za-z0-9_ +.-]{1,31})\]\s*(?:\([^)]*\))?\s*$`)
	trailingGroupPattern      = regexp.MustCompile(`(?i)([A-Za-z][A-Za-z0-9_+-]{1,31})\s*$`)
	advancedGroupContext      = regexp.MustCompile(`(?i)\b(?:REMUX|TrueHD|Atmos|DTS[ .-]?HD|HEVC|AVC|HDR10|Dolby[ .-]+Vision|DDP[ .-]*[257][ .][01]|E-?AC-?3)\b`)
)

type tokenMatch struct {
	start int
	end   int
	raw   string
}

func augmentTorrentInfo(info *TorrentInfo, filename string) {
	normalized := normalizeReleaseString(filename)
	info.Title = normalizeTitleText(info.Title)
	applyEpisodeRanges(info, normalized)
	applyPart(info, normalized)
	applyQuality(info, normalized)
	applyCodec(info, normalized)
	applyHDR(info, normalized)
	applyAudio(info, normalized)
	applySource(info, normalized)
	applyBitDepth(info, normalized)
	applyEdition(info, normalized)
	applyFlags(info, normalized)
	applyGroup(info, normalized)
}

func normalizeReleaseString(value string) string {
	value = html.UnescapeString(strings.TrimSpace(value))
	value = strings.ReplaceAll(value, "_", " ")
	value = strings.ReplaceAll(value, " 039 ", "'")
	value = strings.ReplaceAll(value, "039", "'")
	value = strings.NewReplacer("[", " [", "]", "] ", "(", " (", ")", ") ").Replace(value)
	return collapseSpaces(value)
}

func normalizeTitleText(value string) string {
	value = html.UnescapeString(strings.TrimSpace(value))
	value = strings.ReplaceAll(value, " 039 ", "'")
	value = strings.ReplaceAll(value, "039", "'")
	value = strings.Trim(value, ".-_ ")
	return collapseSpaces(value)
}

func collapseSpaces(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

func applyEpisodeRanges(info *TorrentInfo, value string) {
	if match := seasonEpisodeRangePattern.FindStringSubmatch(value); len(match) == 4 {
		info.Season = atoiOrZero(match[1])
		info.Episode = atoiOrZero(match[2])
		if match[3] != "" {
			info.EpisodeEnd = atoiOrZero(match[3])
		}
		return
	}
	if match := compactRangePattern.FindStringSubmatch(value); len(match) == 4 {
		info.Season = atoiOrZero(match[1])
		info.Episode = atoiOrZero(match[2])
		info.EpisodeEnd = atoiOrZero(match[3])
	}
}

func applyPart(info *TorrentInfo, value string) {
	if match := partPattern.FindStringSubmatch(value); len(match) == 2 {
		info.Part = parsePart(match[1])
	}
}

func applyQuality(info *TorrentInfo, value string) {
	matches := qualityTokenPattern.FindAllString(value, -1)
	for _, match := range matches {
		normalized := normalizeQuality(match)
		if normalized == "REMUX" || info.Quality == "" || info.Quality == "BluRay" && strings.EqualFold(normalized, "REMUX") {
			info.Quality = normalized
		}
	}
}

func applyCodec(info *TorrentInfo, value string) {
	matches := codecTokenPattern.FindAllString(value, -1)
	if len(matches) == 0 {
		return
	}

	hasAV1 := false
	hasH265Tag := false
	hasX265 := false
	hasH264Tag := false
	hasX264 := false
	hasXvid := false
	for _, match := range matches {
		normalized := normalizeCodec(match)
		switch normalized {
		case "AV1":
			hasAV1 = true
		case "H265":
			hasH265Tag = true
		case "x265":
			hasX265 = true
		case "H264":
			hasH264Tag = true
		case "x264":
			hasX264 = true
		case "XViD":
			hasXvid = true
		}
	}

	switch {
	case hasAV1:
		info.Codec = "AV1"
	case hasH265Tag:
		info.Codec = "H265"
	case hasX265:
		info.Codec = "x265"
	case hasH264Tag:
		info.Codec = "H264"
	case hasX264:
		info.Codec = "x264"
	case hasXvid:
		info.Codec = "XViD"
	}
}

func applyHDR(info *TorrentInfo, value string) {
	tokens := orderedNormalizedTokens(value, hdrTokenPattern, normalizeHDR)
	if len(tokens) > 0 {
		info.HDR = strings.Join(tokens, " ")
	}
}

func applyAudio(info *TorrentInfo, value string) {
	tokens := orderedNormalizedTokens(value, audioTokenPattern, normalizeAudioRich)
	if len(tokens) == 0 {
		return
	}
	info.Audio = strings.Join(tokens, " ")
}

func applySource(info *TorrentInfo, value string) {
	releaseStart := firstReleaseTokenPosition(value)
	if releaseStart < 0 {
		return
	}
	for _, match := range sourceMatches(value) {
		if match.start < releaseStart {
			continue
		}
		source := normalizeSource(match.raw)
		if source != "" {
			info.Source = source
			return
		}
	}
}

func applyBitDepth(info *TorrentInfo, value string) {
	if match := bitDepthPattern.FindStringSubmatch(value); len(match) == 2 {
		info.BitDepth = match[1] + "-bit"
	}
}

func applyEdition(info *TorrentInfo, value string) {
	releaseStart := firstReleaseTokenPosition(value)
	tokens := make([]tokenMatch, 0)
	for _, match := range editionTokenPattern.FindAllStringIndex(value, -1) {
		if releaseStart >= 0 && match[0] < releaseStart {
			continue
		}
		tokens = append(tokens, tokenMatch{start: match[0], end: match[1], raw: value[match[0]:match[1]]})
	}
	sort.SliceStable(tokens, func(i, j int) bool { return tokens[i].start < tokens[j].start })

	normalized := make([]string, 0, len(tokens))
	seen := map[string]struct{}{}
	for _, token := range tokens {
		edition := normalizeEdition(token.raw)
		if edition == "" {
			continue
		}
		if _, ok := seen[edition]; ok {
			continue
		}
		seen[edition] = struct{}{}
		normalized = append(normalized, edition)
	}
	if len(normalized) > 0 {
		info.Edition = strings.Join(normalized, " ")
	}
}

func applyFlags(info *TorrentInfo, value string) {
	info.Complete = info.Complete || completePattern.MatchString(value)
	info.Remastered = info.Remastered || regexp.MustCompile(`(?i)\bRemastered\b`).MatchString(value)
	info.IMAX = info.IMAX || regexp.MustCompile(`(?i)\bIMAX\b`).MatchString(value)
}

func applyGroup(info *TorrentInfo, value string) {
	if group := cleanGroup(info.Group); group != "" {
		info.Group = group
		return
	}
	if matches := dashGroupPattern.FindAllStringSubmatch(value, -1); len(matches) > 0 {
		if group := cleanGroup(matches[len(matches)-1][1]); group != "" {
			info.Group = group
			return
		}
	}

	trimmed := strings.TrimSpace(regexp.MustCompile(`\([^)]*\)\s*$`).ReplaceAllString(value, ""))
	if match := bracketGroupPattern.FindStringSubmatch(trimmed); len(match) == 2 {
		if group := cleanGroup(match[1]); group != "" && !looksLikeHash(group) && !isKnownBracketTag(group) {
			info.Group = group
			return
		}
	}

	releaseStart := firstReleaseTokenPosition(trimmed)
	if releaseStart < 0 || !advancedGroupContext.MatchString(trimmed[releaseStart:]) {
		return
	}
	if match := trailingGroupPattern.FindStringSubmatchIndex(trimmed); len(match) == 4 && match[2] >= releaseStart {
		group := cleanGroup(trimmed[match[2]:match[3]])
		if group != "" && !looksLikeMetadataToken(group) {
			info.Group = group
		}
	}
}

func firstReleaseTokenPosition(value string) int {
	patterns := []*regexp.Regexp{
		seasonEpisodeRangePattern,
		yearTokenPattern,
		resolutionTokenPattern,
		qualityTokenPattern,
		codecTokenPattern,
		hdrTokenPattern,
		audioTokenPattern,
		bitDepthPattern,
	}
	first := -1
	for _, pattern := range patterns {
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
	sort.SliceStable(matches, func(i, j int) bool { return matches[i].start < matches[j].start })

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
	sort.SliceStable(matches, func(i, j int) bool { return matches[i].start < matches[j].start })
	return matches
}

func normalizeHDR(value string) string {
	collapsed := strings.ToLower(strings.NewReplacer(" ", "", "-", "").Replace(value))
	switch collapsed {
	case "hdr":
		return "HDR"
	case "hdr10":
		return "HDR10"
	case "hdr10+", "hdr10plus", "hdr10p":
		return "HDR10+"
	case "dv", "dovi", "dolbyvision":
		return "DV"
	case "hlg":
		return "HLG"
	default:
		return ""
	}
}

func normalizeAudioRich(value string) string {
	collapsed := strings.ToLower(strings.NewReplacer(" ", "", ".", "", "-", "").Replace(value))
	switch {
	case strings.HasPrefix(collapsed, "truehdatmos"):
		return "TrueHD Atmos" + normalizeOptionalChannel(collapsed, "truehdatmos")
	case strings.HasPrefix(collapsed, "truehd"):
		return "TrueHD" + normalizeOptionalChannel(collapsed, "truehd")
	case strings.HasPrefix(collapsed, "dtshdma"):
		return "DTS-HD MA" + normalizeOptionalChannel(collapsed, "dtshdma")
	case strings.HasPrefix(collapsed, "dtshd"):
		return "DTS-HD" + normalizeOptionalChannel(collapsed, "dtshd")
	case strings.HasPrefix(collapsed, "ddpatmos"):
		return "DDP Atmos" + normalizeOptionalChannel(collapsed, "ddpatmos")
	case strings.HasPrefix(collapsed, "ddp"):
		channel := normalizeChannelFromCollapsed(collapsed)
		hasAtmos := strings.Contains(collapsed, "atmos")
		if strings.HasSuffix(channel, "5.1") || strings.HasSuffix(channel, "7.1") || strings.HasSuffix(channel, "2.0") {
			if hasAtmos {
				return "DDP" + strings.TrimPrefix(channel, " ") + " Atmos"
			}
			return "DDP" + strings.TrimPrefix(channel, " ")
		}
		if hasAtmos {
			return "DDP Atmos"
		}
		return "DDP"
	case strings.HasPrefix(collapsed, "eac3"):
		return "EAC3" + normalizeOptionalChannel(collapsed, "eac3")
	case collapsed == "dd+":
		return "DD+"
	case strings.HasPrefix(collapsed, "aac"):
		return "AAC" + strings.TrimPrefix(normalizeOptionalChannel(collapsed, "aac"), " ")
	case strings.HasPrefix(collapsed, "ac3"):
		return "AC3" + normalizeOptionalChannel(collapsed, "ac3")
	case strings.HasPrefix(collapsed, "flac"):
		return "FLAC" + normalizeOptionalChannel(collapsed, "flac")
	case collapsed == "atmos":
		return "Atmos"
	default:
		return normalizeAudio(value)
	}
}

func normalizeChannelFromCollapsed(collapsed string) string {
	switch {
	case strings.Contains(collapsed, "20"):
		return " 2.0"
	case strings.Contains(collapsed, "51"):
		return " 5.1"
	case strings.Contains(collapsed, "71"):
		return " 7.1"
	default:
		return ""
	}
}

func normalizeOptionalChannel(collapsed string, prefix string) string {
	suffix := strings.TrimPrefix(collapsed, prefix)
	switch suffix {
	case "20":
		return " 2.0"
	case "51":
		return " 5.1"
	case "71":
		return " 7.1"
	default:
		return ""
	}
}

func normalizeSource(value string) string {
	upper := strings.ToUpper(strings.TrimSpace(value))
	switch upper {
	case "AMZN", "NF", "DSNP", "HULU", "CR", "ATVP", "PCOK", "HMAX":
		return upper
	case "HBO":
		return "HBO"
	case "MAX":
		return "MAX"
	case "IT":
		return "iT"
	default:
		return ""
	}
}

func normalizeEdition(value string) string {
	collapsed := strings.ToLower(strings.NewReplacer(" ", "", "-", "", "'", "").Replace(value))
	switch collapsed {
	case "directorscut", "dc":
		return "Director's Cut"
	case "hybrid":
		return "Hybrid"
	case "bw", "b&w", "blackandwhite":
		return "Black and White"
	case "dubbed":
		return "Dubbed"
	case "dual", "dualaudio":
		return "Dual Audio"
	case "multisub", "multisubs":
		return "Multi Subs"
	default:
		return ""
	}
}

func cleanGroup(value string) string {
	value = strings.TrimSpace(value)
	value = strings.Trim(value, "[]{}() ")
	value = strings.ReplaceAll(value, " ", "_")
	if value == "" || atoiOrZero(value) > 0 || looksLikeMetadataToken(value) || looksLikeHash(value) {
		return ""
	}
	return value
}

func looksLikeMetadataToken(value string) bool {
	lowered := strings.ToLower(strings.TrimSpace(value))
	switch lowered {
	case "mkv", "mp4", "avi", "web", "webrip", "webdl", "bluray", "remux", "hevc", "avc", "aac", "atmos", "truehd", "multi", "multisub", "multisubs", "dual", "audio", "subs", "sub":
		return true
	default:
		return false
	}
}

func looksLikeHash(value string) bool {
	if len(value) < 6 || len(value) > 16 {
		return false
	}
	for _, r := range value {
		if !('0' <= r && r <= '9' || 'a' <= r && r <= 'f' || 'A' <= r && r <= 'F') {
			return false
		}
	}
	return true
}

func isKnownBracketTag(value string) bool {
	normalized := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(value), " ", ""))
	switch normalized {
	case "multisub", "multisubs", "dualaudio", "weekly":
		return true
	default:
		return false
	}
}

func atoiOrZero(raw string) int {
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0
	}
	return value
}

func parsePart(raw string) int {
	if value := atoiOrZero(raw); value > 0 {
		return value
	}
	switch strings.ToLower(raw) {
	case "one", "i":
		return 1
	case "two", "ii":
		return 2
	case "three", "iii":
		return 3
	case "four", "iv":
		return 4
	case "five", "v":
		return 5
	case "six", "vi":
		return 6
	case "seven", "vii":
		return 7
	case "eight", "viii":
		return 8
	case "nine", "ix":
		return 9
	case "ten", "x":
		return 10
	default:
		return 0
	}
}
