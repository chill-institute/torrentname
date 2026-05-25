package torrentname

import (
	"regexp"
	"strings"
)

var (
	audioChannelTokenPattern  = `[257][ .]?[01]|[268][ .-]*CH`
	seasonEpisodeRangePattern = regexp.MustCompile(`(?i)\bS([0-9]{1,2})E([0-9]{1,3})(?:[ .-]*E?([0-9]{1,3}))?\b`)
	seasonReleaseTokenPattern = regexp.MustCompile(`(?i)\bSeason[ .-]+[0-9]{1,2}\b`)
	compactRangePattern       = regexp.MustCompile(`(?i)\bS([0-9]{1,2})E([0-9]{1,3})[ .-]+([0-9]{1,3})(?:[ .-]+of[ .-]+[0-9]{1,3})?\b`)
	yearTokenPattern          = regexp.MustCompile(`\b(?:19[0-9]{2}|20[0-9]{2})\b`)
	resolutionTokenPattern    = regexp.MustCompile(`(?i)\b(?:` + tokenPatternAlternates(resolutionCatalog, broadResolutionAliasContextPattern, `[0-9]{3,4}p`) + `)\b`)
	qualityTokenPattern       = compileTokenPattern(qualityCatalog)
	codecTokenPattern         = compileTokenPattern(codecCatalog)
	hdrTokenPattern           = compileLooseEndTokenPattern(hdrCatalog)
	audioTokenPattern         = compileTokenPattern(audioCatalog)
	bitDepthPattern           = regexp.MustCompile(`(?i)\b(8|10|12|16|24)[ .-]?bits?\b`)
	partPattern               = regexp.MustCompile(`(?i)\bPart[ .-]+([0-9]+|One|Two|Three|Four|Five|Six|Seven|Eight|Nine|Ten|I|II|III|IV|V|VI|VII|VIII|IX|X)\b`)
	completePattern           = regexp.MustCompile(`(?i)\b(?:Complete(?:[ .-]+(?:Season|Series))?|Season[ .-]+[0-9]{1,2}(?:[ .-]*(?:to|-|\+|&)[ .-]*[0-9]{1,2})?[ .-]+Complete(?:[ .-]+Series)?|Seasons[ .-]+[0-9]{1,2}[ .-]+to[ .-]+[0-9]{1,2}[ .-]+Complete|S[0-9]{1,2}[ .-]*Complete)\b`)
	editionTokenPattern       = compileTokenPattern(editionCatalog)
	languageTokenPattern      = compileTokenPattern(languageCatalog)
	remasteredPattern         = regexp.MustCompile(`(?i)\b(?:Remastered|RM4K|4K[ .-]+Remaster(?:ed)?)\b`)
	imaxPattern               = regexp.MustCompile(`(?i)\bIMAX\b`)
)

func augmentTorrentInfo(info *TorrentInfo, filename string) {
	normalized := normalizeReleaseString(filename)
	releaseStart := firstReleaseTokenPosition(normalized)
	info.Title = normalizeTitleText(info.Title)
	applyEpisodeRanges(info, normalized)
	applyPart(info, normalized)
	applyResolution(info, normalized)
	applyQuality(info, normalized)
	applyCodec(info, normalized)
	applyHDR(info, normalized)
	applyAudio(info, normalized)
	applySource(info, normalized, releaseStart)
	applyLanguage(info, normalized, releaseStart)
	applyBitDepth(info, normalized)
	applyEdition(info, normalized, releaseStart)
	applyFlags(info, normalized)
	applyGroup(info, normalized, releaseStart)
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

func applyResolution(info *TorrentInfo, value string) {
	matches := resolutionTokenPattern.FindAllStringIndex(value, -1)
	fallback := ""
	skippedCurrentResolution := false
	for _, match := range matches {
		raw := value[match[0]:match[1]]
		normalized := normalizeResolution(raw)
		if normalized == "" {
			continue
		}
		if isBlockedResolutionAlias(value, match[0], match[1]) {
			if normalized == info.Resolution {
				skippedCurrentResolution = true
			}
			continue
		}
		if isNumericResolutionToken(raw) {
			info.Resolution = normalized
			return
		}
		if fallback == "" {
			fallback = normalized
		}
	}
	if fallback != "" {
		info.Resolution = fallback
		return
	}
	if skippedCurrentResolution {
		info.Resolution = ""
	}
}

func isBlockedResolutionAlias(value string, start int, end int) bool {
	if compactKey(value[start:end]) != "HD" {
		return false
	}
	suffix := strings.TrimLeft(value[end:], ".-_ ")
	suffixKey := compactKey(suffix)
	return strings.HasPrefix(suffixKey, "TS") ||
		strings.HasPrefix(suffixKey, "TC") ||
		strings.HasPrefix(suffixKey, "CAM")
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
	tokens, skippedCurrentAudio := orderedAudioTokens(value, info.Audio)
	if len(tokens) == 0 {
		if skippedCurrentAudio {
			info.Audio = ""
		}
		return
	}
	if hasDualAudioToken(tokens) {
		withoutDual := nonDualAudioTokens(tokens)
		if len(withoutDual) == 0 {
			if isDualAudioToken(info.Audio) {
				return
			}
			info.Audio = "Dual Audio"
			return
		}
		if allBareChannelAudio(withoutDual) {
			info.Audio = firstDualAudioToken(tokens)
			if info.Audio == "" {
				info.Audio = "Dual Audio"
			}
			return
		}
		info.Audio = strings.Join(withoutDual, " ")
		return
	}
	info.Audio = strings.Join(tokens, " ")
}

func orderedAudioTokens(value string, currentAudio string) ([]string, bool) {
	matches := make([]tokenMatch, 0)
	for _, match := range audioTokenPattern.FindAllStringIndex(value, -1) {
		matches = append(matches, tokenMatch{start: match[0], end: match[1], raw: value[match[0]:match[1]]})
	}
	sortTokenMatches(matches)

	out := make([]string, 0, len(matches))
	seen := map[string]struct{}{}
	skippedCurrentAudio := false
	for _, match := range matches {
		normalized := normalizeAudioRich(match.raw)
		if normalized == "" {
			continue
		}
		if isLineTVAudioMatch(value, match) {
			if normalized == normalizeAudioRich(currentAudio) {
				skippedCurrentAudio = true
			}
			continue
		}
		if _, ok := seen[normalized]; ok {
			continue
		}
		seen[normalized] = struct{}{}
		out = append(out, normalized)
	}
	return out, skippedCurrentAudio
}

func isLineTVAudioMatch(value string, match tokenMatch) bool {
	if compactKey(match.raw) != "LINE" {
		return false
	}
	suffix := strings.TrimLeft(value[match.end:], ".-_ ")
	return strings.HasPrefix(compactKey(suffix), "TV")
}

func isBareChannelAudio(value string) bool {
	return value == "2.0" || value == "5.1" || value == "7.1"
}

func hasDualAudioToken(tokens []string) bool {
	for _, token := range tokens {
		if isDualAudioToken(token) {
			return true
		}
	}
	return false
}

func nonDualAudioTokens(tokens []string) []string {
	out := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if !isDualAudioToken(token) {
			out = append(out, token)
		}
	}
	return out
}

func isDualAudioToken(value string) bool {
	return value == "Dual Audio" || value == "Dual-Audio"
}

func firstDualAudioToken(tokens []string) string {
	for _, token := range tokens {
		if isDualAudioToken(token) {
			return token
		}
	}
	return ""
}

func allBareChannelAudio(tokens []string) bool {
	for _, token := range tokens {
		if !isBareChannelAudio(token) {
			return false
		}
	}
	return len(tokens) > 0
}

func applySource(info *TorrentInfo, value string, releaseStart int) {
	if releaseStart < 0 {
		return
	}
	for _, match := range sourceMatches(value) {
		if match.start < releaseStart {
			continue
		}
		source := normalizeSource(match.raw)
		if source != "" && !hasSourceTokenContext(value, match, source) {
			continue
		}
		if source != "" && isFinalGroupSourceToken(info, value, match, source) {
			continue
		}
		if source != "" {
			info.Source = source
			return
		}
	}
}

func hasSourceTokenContext(value string, match tokenMatch, source string) bool {
	if !requiresAdjacentReleaseContext(source) {
		return true
	}
	return hasQualityBefore(value, match.start) || hasQualityAfter(value, match.end)
}

func requiresAdjacentReleaseContext(source string) bool {
	if token, ok := sourceLookup[compactKey(source)]; ok {
		return token.requiresContext
	}
	return false
}

func hasQualityBefore(value string, end int) bool {
	prefix := strings.TrimRight(value[:end], ".-_ ")
	matches := qualityTokenPattern.FindAllStringIndex(prefix, -1)
	return len(matches) > 0 && matches[len(matches)-1][1] == len(prefix)
}

func hasQualityAfter(value string, start int) bool {
	suffix := strings.TrimLeft(value[start:], ".-_ ")
	if match := qualityTokenPattern.FindStringIndex(suffix); match != nil && match[0] == 0 {
		return true
	}
	resolutionMatch := resolutionTokenPattern.FindStringIndex(suffix)
	if resolutionMatch == nil || resolutionMatch[0] != 0 {
		return false
	}
	if resolutionMatchContainsQuality(suffix[:resolutionMatch[1]]) {
		return true
	}
	afterResolution := strings.TrimLeft(suffix[resolutionMatch[1]:], ".-_ ")
	for {
		hdrMatch := hdrTokenPattern.FindStringIndex(afterResolution)
		if hdrMatch == nil || hdrMatch[0] != 0 {
			break
		}
		afterResolution = strings.TrimLeft(afterResolution[hdrMatch[1]:], ".-_ ")
	}
	match := qualityTokenPattern.FindStringIndex(afterResolution)
	return match != nil && match[0] == 0
}

func resolutionMatchContainsQuality(value string) bool {
	match := qualityTokenPattern.FindStringIndex(value)
	return match != nil && match[0] > 0
}

func isFinalGroupSourceToken(info *TorrentInfo, value string, match tokenMatch, source string) bool {
	if match.start == 0 {
		return false
	}
	if isFinalBracketSourceGroup(value, match, source) {
		return true
	}
	separator := value[match.start-1]
	if separator != '-' && separator != '.' {
		return false
	}
	if separator == '.' && !hasDotSuffixGroupContext(value, match.start) {
		return false
	}
	suffix := trimKnownContainerSuffix(value[match.start:])
	hasGroupSuffix := strings.EqualFold(cleanGroup(suffix), source)
	if separator == '-' && hasGroupSuffix {
		return true
	}
	if separator == '.' && hasGroupSuffix {
		return true
	}
	if !isAmbiguousFinalSourceGroup(source) && (info.Group == "" || !strings.EqualFold(info.Group, source)) {
		return false
	}
	return hasGroupSuffix || info.Group != "" && strings.EqualFold(info.Group, source)
}

func isFinalBracketSourceGroup(value string, match tokenMatch, source string) bool {
	if match.start == 0 || value[match.start-1] != '[' {
		return false
	}
	rest := strings.TrimSpace(value[match.end:])
	if !strings.HasPrefix(rest, "]") {
		return false
	}
	suffix := strings.TrimSpace(strings.TrimPrefix(rest, "]"))
	suffix = strings.TrimSpace(trimKnownContainerSuffix(suffix))
	if suffix != "" {
		return false
	}
	return strings.EqualFold(cleanGroup(source), source)
}

func hasDotSuffixGroupContext(value string, end int) bool {
	previous := previousReleaseToken(value, end)
	if previous == "" {
		return false
	}
	return codecTokenPattern.MatchString(previous) ||
		hdrTokenPattern.MatchString(previous) ||
		audioTokenPattern.MatchString(previous) ||
		bitDepthPattern.MatchString(previous)
}

func previousReleaseToken(value string, end int) string {
	prefix := strings.TrimRight(value[:end], ".-_ ")
	start := len(prefix)
	for start > 0 {
		switch prefix[start-1] {
		case '.', '-', '_', ' ':
			return prefix[start:]
		default:
			start--
		}
	}
	return prefix
}

func isAmbiguousFinalSourceGroup(source string) bool {
	if token, ok := sourceLookup[compactKey(source)]; ok {
		return token.ambiguousGroup
	}
	return false
}

func applyLanguage(info *TorrentInfo, value string, releaseStart int) {
	if info.Language != "" {
		return
	}
	if releaseStart < 0 {
		return
	}

	releaseSuffix := value[releaseStart:]
	tokens := orderedNormalizedTokens(releaseSuffix, languageTokenPattern, normalizeLanguage)
	hasExplicitMulti := hasMultiLanguageMarker(releaseSuffix)
	for _, language := range tokens {
		isSubtitleMulti := language == "MULTI" && hasMultiSubtitleMarker(releaseSuffix)
		if language == "MULTI" && !isSubtitleMulti || language == "VOSTFR" || language == "VFF" || language == "VFQ" {
			hasExplicitMulti = true
			break
		}
	}
	if len(tokens) >= 2 || hasExplicitMulti {
		info.Language = strings.Join(tokens, " ")
		return
	}
	if len(tokens) == 1 && isStrongSingleLanguageMarker(tokens[0]) {
		info.Language = tokens[0]
	}
}

func isStrongSingleLanguageMarker(value string) bool {
	return value == "JPN" || value == "KOR"
}

func hasMultiSubtitleMarker(value string) bool {
	return strings.Contains(compactKey(value), "MULTISUB")
}

func hasMultiLanguageMarker(value string) bool {
	return strings.Contains(compactKey(value), "MULTILANG")
}

func applyBitDepth(info *TorrentInfo, value string) {
	if match := bitDepthPattern.FindStringSubmatch(value); len(match) == 2 {
		info.BitDepth = match[1] + "-bit"
	}
}

func applyEdition(info *TorrentInfo, value string, releaseStart int) {
	tokens := make([]tokenMatch, 0)
	for _, match := range editionTokenPattern.FindAllStringIndex(value, -1) {
		if releaseStart >= 0 && match[0] < releaseStart {
			continue
		}
		tokens = append(tokens, tokenMatch{start: match[0], end: match[1], raw: value[match[0]:match[1]]})
	}
	sortTokenMatches(tokens)

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
	info.Remastered = info.Remastered || remasteredPattern.MatchString(value)
	info.IMAX = info.IMAX || imaxPattern.MatchString(value)
}
