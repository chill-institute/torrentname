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
		if strings.Contains(collapsed, "atmos") {
			return "TrueHD Atmos" + normalizeChannelFromCollapsed(collapsed)
		}
		return "TrueHD" + normalizeOptionalChannel(collapsed, "truehd")
	case strings.HasPrefix(collapsed, "dtshdma"):
		return "DTS-HD MA" + normalizeOptionalChannel(collapsed, "dtshdma")
	case strings.HasPrefix(collapsed, "dtshdhra"):
		return "DTS-HD HRA" + normalizeOptionalChannel(collapsed, "dtshdhra")
	case strings.HasPrefix(collapsed, "dtshd"):
		return "DTS-HD" + normalizeOptionalChannel(collapsed, "dtshd")
	case strings.HasPrefix(collapsed, "dtsx"):
		return "DTS X" + normalizeOptionalChannel(collapsed, "dtsx")
	case strings.HasPrefix(collapsed, "dts"):
		return "DTS" + normalizeOptionalChannel(collapsed, "dts")
	case strings.HasPrefix(collapsed, "dd+atmos"):
		return "DD+ Atmos" + normalizeOptionalChannel(collapsed, "dd+atmos")
	case strings.HasPrefix(collapsed, "ddplusatmos"):
		return "DD+ Atmos" + normalizeOptionalChannel(collapsed, "ddplusatmos")
	case strings.HasPrefix(collapsed, "dd+"):
		return normalizeDDPlus(collapsed)
	case strings.HasPrefix(collapsed, "ddplus"):
		return normalizeDDPlus(collapsed)
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
		if strings.Contains(collapsed, "atmos") {
			return "EAC3 Atmos" + normalizeChannelFromCollapsed(collapsed)
		}
		return "EAC3" + normalizeOptionalChannel(collapsed, "eac3")
	case strings.HasPrefix(collapsed, "aac"):
		return "AAC" + strings.TrimPrefix(normalizeOptionalChannel(collapsed, "aac"), " ")
	case strings.HasPrefix(collapsed, "ac3"):
		return "AC3" + normalizeOptionalChannel(collapsed, "ac3")
	case strings.HasPrefix(collapsed, "flac"):
		return "FLAC" + normalizeOptionalChannel(collapsed, "flac")
	case strings.HasPrefix(collapsed, "lpcm"):
		return "LPCM" + normalizeOptionalChannel(collapsed, "lpcm")
	case strings.HasPrefix(collapsed, "pcm"):
		return "PCM" + normalizeOptionalChannel(collapsed, "pcm")
	case strings.HasPrefix(collapsed, "opus"):
		return "Opus" + normalizeOptionalChannel(collapsed, "opus")
	case collapsed == "2ch", collapsed == "6ch", collapsed == "8ch":
		return strings.TrimSpace(normalizeChannelFromCollapsed(collapsed))
	case collapsed == "atmos":
		return "Atmos"
	default:
		return normalizeAudio(value)
	}
}

func compactUpperToken(value string) string {
	return strings.ToUpper(strings.NewReplacer(" ", "", ".", "", "-", "", "_", "").Replace(strings.TrimSpace(value)))
}

func normalizeDDPlus(collapsed string) string {
	channel := normalizeChannelFromCollapsed(collapsed)
	if strings.Contains(collapsed, "atmos") {
		return "DD+ Atmos" + channel
	}
	return "DD+" + channel
}

func normalizeChannelFromCollapsed(collapsed string) string {
	switch {
	case strings.Contains(collapsed, "20"):
		return formatChannel("20")
	case strings.Contains(collapsed, "51"):
		return formatChannel("51")
	case strings.Contains(collapsed, "71"):
		return formatChannel("71")
	case strings.Contains(collapsed, "2ch"):
		return formatChannel("2ch")
	case strings.Contains(collapsed, "6ch"):
		return formatChannel("6ch")
	case strings.Contains(collapsed, "8ch"):
		return formatChannel("8ch")
	default:
		return ""
	}
}

func normalizeOptionalChannel(collapsed string, prefix string) string {
	return formatChannel(strings.TrimPrefix(collapsed, prefix))
}

func formatChannel(value string) string {
	switch value {
	case "20", "2ch":
		return " 2.0"
	case "51", "6ch":
		return " 5.1"
	case "71", "8ch":
		return " 7.1"
	default:
		return ""
	}
}

func normalizeSource(value string) string {
	upper := compactUpperToken(value)
	switch upper {
	case "AMZN", "AUBC", "CBC", "CPNG", "CR", "DSNP", "FOD", "HAMI", "HBO", "HMAX", "HULU", "HTSR", "IQIY", "ITVX", "KCW", "KKTV", "LINETV", "MY5", "MYTVSUPER", "NF", "NOW", "OVID", "PCOK", "PLAY", "PMTP", "ROKU", "STAN", "TVING", "TVER", "UNEXT", "VIKI", "VIU", "VRV", "WAVVE", "WETV", "YOUKU":
		return upper
	case "ATV", "ATV+", "ATVP":
		return "ATVP"
	case "BCORE":
		return "BCORE"
	case "CRAVE", "CRAV":
		return "CRAVE"
	case "CRIT":
		return "CRiT"
	case "FRIDAY":
		return "friDay"
	case "HBOM", "HBOMAX":
		return "HMAX"
	case "MAX":
		return "MAX"
	case "IT", "ITUNES":
		return "iT"
	case "IP":
		return "iP"
	case "STAR", "STAR+":
		return "STAR+"
	case "STRP":
		return "STRP"
	default:
		return ""
	}
}

func normalizeLanguage(value string) string {
	upper := compactUpperToken(value)
	switch upper {
	case "MULTI", "MULTILANG":
		return "MULTI"
	case "VOSTFR", "VFF", "VFQ":
		return upper
	case "ENG", "ENGLISH":
		return "ENG"
	case "ITA", "ITALIAN":
		return "ITA"
	case "FRE", "FRENCH":
		return "FRE"
	case "GER", "GERMAN":
		return "GER"
	case "SPA", "SPANISH":
		return "SPA"
	case "RUS", "RUSSIAN":
		return "RUS"
	case "JPN", "JPS", "JAP", "JAPANESE":
		return "JPN"
	case "UKR", "UKRAINIAN":
		return "UKR"
	case "HIN", "HINDI":
		return "HIN"
	case "KOR", "KOREAN":
		return "KOR"
	case "CHI", "CHINESE":
		return "CHI"
	default:
		return ""
	}
}

func normalizeEdition(value string) string {
	collapsed := strings.ToLower(strings.NewReplacer(" ", "", ".", "", "-", "", "'", "").Replace(value))
	switch collapsed {
	case "directorscut", "dc":
		return "Director's Cut"
	case "hybrid":
		return "Hybrid"
	case "theatrical", "theatricalcut":
		return "Theatrical Cut"
	case "specialedition":
		return "Special Edition"
	case "openmatte":
		return "Open Matte"
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
