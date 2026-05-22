package torrentname

import (
	"regexp"
	"strings"
)

var (
	dashGroupPattern           = regexp.MustCompile(`(?i)-\s*([A-Za-z0-9][A-Za-z0-9._=+\[\]{}-]{0,31})(?:\s*(?:\(|\[|$))`)
	bracketGroupPattern        = regexp.MustCompile(`(?i)\[([A-Za-z][A-Za-z0-9_ +.-]{1,31})\]\s*(?:\([^)]*\))?\s*$`)
	trailingGroupPattern       = regexp.MustCompile(`(?i)([A-Za-z][A-Za-z0-9_+-]{1,31})\s*$`)
	advancedGroupContext       = regexp.MustCompile(`(?i)\b(?:REMUX|TrueHD|Atmos|DTS[ .-]?HD|HEVC|AVC|HDR10|Dolby[ .-]+Vision|DDP[ .-]*[257][ .][01]|E-?AC-?3)\b`)
	legacyDotGroupContext      = regexp.MustCompile(`(?i)\b(?:WEB[ .-]?DL|WEB[ .-]?Rip|WEB|HDTV|HDRip|Blu[ .-]?Ray|BDRip|BRRip|x[ ._-]?26[45]|h[ ._-]?26[45]|XViD|AAC|AC3)\b`)
	trailingParenthesesPattern = regexp.MustCompile(`\([^)]*\)\s*$`)
	episodeVersionGroupPattern = regexp.MustCompile(`(?i)^[0-9]{1,4}v[0-9]+$`)
	seasonRangeGroupPattern    = regexp.MustCompile(`(?i)^[0-9]{1,2}[._-]?Complete\b`)
)

func applyGroup(info *TorrentInfo, value string, releaseStart int) {
	if info.Group != "" {
		if group := cleanGroup(info.Group); group != "" {
			if assignGroup(info, group) {
				return
			}
		}
		info.Group = ""
	}
	if matches := dashGroupPattern.FindAllStringSubmatch(value, -1); len(matches) > 0 {
		if group := cleanGroup(matches[len(matches)-1][1]); group != "" {
			if assignGroup(info, group) {
				return
			}
		}
	}

	trimmed := trimKnownContainerSuffix(strings.TrimSpace(trailingParenthesesPattern.ReplaceAllString(value, "")))
	if match := bracketGroupPattern.FindStringSubmatch(trimmed); len(match) == 2 {
		if group := cleanGroup(match[1]); group != "" && !looksLikeHash(group) && !isKnownBracketTag(group) {
			if assignGroup(info, group) {
				return
			}
		}
	}

	if releaseStart < 0 {
		return
	}
	if match := trailingGroupPattern.FindStringSubmatchIndex(trimmed); len(match) == 4 && match[2] >= releaseStart {
		hasAdvancedContext := advancedGroupContext.MatchString(trimmed[releaseStart:])
		hasLegacyDotContext := match[2] > 0 && trimmed[match[2]-1] == '.' && legacyDotGroupContext.MatchString(trimmed[releaseStart:])
		if !hasAdvancedContext && !hasLegacyDotContext {
			return
		}
		group := cleanGroup(trimmed[match[2]:match[3]])
		if group != "" && !looksLikeMetadataToken(group) {
			assignGroup(info, group)
		}
	}
}

func assignGroup(info *TorrentInfo, group string) bool {
	if info.Source != "" && strings.EqualFold(group, info.Source) {
		return false
	}
	info.Group = group
	return true
}

func cleanGroup(value string) string {
	value = strings.TrimSpace(value)
	value = strings.Trim(value, "[]{}() ")
	value = trimKnownContainerSuffix(value)
	value = strings.ReplaceAll(value, " ", "_")
	if parts := strings.Split(value, "-"); len(parts) > 1 && (looksLikeMetadataToken(parts[0]) || containsMetadataToken(parts[0])) {
		value = parts[len(parts)-1]
	}
	if containsMetadataToken(value) {
		return ""
	}
	if value == "" || atoiOrZero(value) > 0 || looksLikeMetadataToken(value) || looksLikeHash(value) {
		return ""
	}
	if episodeVersionGroupPattern.MatchString(value) || seasonRangeGroupPattern.MatchString(value) {
		return ""
	}
	return value
}

func trimKnownContainerSuffix(value string) string {
	value = strings.TrimSpace(value)
	lowered := strings.ToLower(value)
	for _, suffix := range []string{".mkv", ".mp4", ".avi"} {
		if strings.HasSuffix(lowered, suffix) {
			return value[:len(value)-len(suffix)]
		}
	}
	return value
}

func containsMetadataToken(value string) bool {
	fields := strings.FieldsFunc(value, func(r rune) bool {
		return r == '.' || r == '_' || r == ' ' || r == '-' || r == '[' || r == ']' || r == '{' || r == '}'
	})
	for _, field := range fields {
		if looksLikeMetadataToken(field) {
			return true
		}
	}
	return false
}

func looksLikeMetadataToken(value string) bool {
	lowered := strings.ToLower(strings.TrimSpace(value))
	switch lowered {
	case "mkv", "mp4", "avi", "web", "dl", "rip", "webrip", "webdl", "hdtv", "pdtv", "hdrip", "dvdrip", "bdrip", "brrip", "bluray", "remux", "hevc", "avc", "av1", "xvid", "x264", "x265", "h264", "h265", "aac", "ac3", "eac3", "dd+", "ddplus", "ddp", "dts", "dts-hd", "flac", "opus", "pcm", "atmos", "truehd", "multi", "multisub", "multisubs", "dual", "audio", "subs", "sub", "proper", "repack", "remastered", "imax", "hdr", "hdr10", "hdr10plus", "dv", "dovi", "hlg":
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
