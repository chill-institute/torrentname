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
	trailingParenthesesPattern = regexp.MustCompile(`\([^)]*\)\s*$`)
)

func applyGroup(info *TorrentInfo, value string) {
	if info.Group != "" {
		if group := cleanGroup(info.Group); group != "" {
			info.Group = group
			return
		}
		info.Group = ""
	}
	if matches := dashGroupPattern.FindAllStringSubmatch(value, -1); len(matches) > 0 {
		if group := cleanGroup(matches[len(matches)-1][1]); group != "" {
			info.Group = group
			return
		}
	}

	trimmed := strings.TrimSpace(trailingParenthesesPattern.ReplaceAllString(value, ""))
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
