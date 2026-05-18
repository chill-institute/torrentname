package torrentname

import (
	"regexp"
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
	bitDepthPattern           = regexp.MustCompile(`(?i)\b(8|10|12)[ .-]?bit\b`)
	partPattern               = regexp.MustCompile(`(?i)\bPart[ .-]+([0-9]+|One|Two|Three|Four|Five|Six|Seven|Eight|Nine|Ten|I|II|III|IV|V|VI|VII|VIII|IX|X)\b`)
	completePattern           = regexp.MustCompile(`(?i)\b(?:Complete(?:[ .-]+(?:Season|Series))?|Season[ .-]+[0-9]{1,2}[ .-]+Complete|Seasons[ .-]+[0-9]{1,2}[ .-]+to[ .-]+[0-9]{1,2}[ .-]+Complete|S[0-9]{1,2}[ .-]*Complete)\b`)
	editionTokenPattern       = regexp.MustCompile(`(?i)\b(?:Director'?s[ .-]+Cut|Directors[ .-]+Cut|DC|Hybrid|B&W|BW|Black[ .-]+and[ .-]+White|DUBBED|DUAL|Dual[ .-]+Audio|Multi[ .-]*Subs?|MultiSub)\b`)
	remasteredPattern         = regexp.MustCompile(`(?i)\bRemastered\b`)
	imaxPattern               = regexp.MustCompile(`(?i)\bIMAX\b`)
)

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
