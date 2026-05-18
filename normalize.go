package torrentname

import (
	"html"
	"strconv"
	"strings"
)

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
