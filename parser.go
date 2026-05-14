package torrentname

import (
	"strings"
)

type TorrentInfo struct {
	Title      string
	Season     int    `json:"season,omitempty"`
	Episode    int    `json:"episode,omitempty"`
	EpisodeEnd int    `json:"episode_end,omitempty"`
	Part       int    `json:"part,omitempty"`
	Year       int    `json:"year,omitempty"`
	Resolution string `json:"resolution,omitempty"`
	Quality    string `json:"quality,omitempty"`
	Codec      string `json:"codec,omitempty"`
	HDR        string `json:"hdr,omitempty"`
	Audio      string `json:"audio,omitempty"`
	Source     string `json:"source,omitempty"`
	Group      string `json:"group,omitempty"`
	Region     string `json:"region,omitempty"`
	Extended   bool   `json:"extended,omitempty"`
	Hardcoded  bool   `json:"hardcoded,omitempty"`
	Proper     bool   `json:"proper,omitempty"`
	Repack     bool   `json:"repack,omitempty"`
	Remastered bool   `json:"remastered,omitempty"`
	Container  string `json:"container,omitempty"`
	Widescreen bool   `json:"widescreen,omitempty"`
	Website    string `json:"website,omitempty"`
	Language   string `json:"language,omitempty"`
	BitDepth   string `json:"bit_depth,omitempty"`
	Edition    string `json:"edition,omitempty"`
	Sbs        string `json:"sbs,omitempty"`
	Unrated    bool   `json:"unrated,omitempty"`
	Size       string `json:"size,omitempty"`
	ThreeD     bool   `json:"3d,omitempty"`
	IMAX       bool   `json:"imax,omitempty"`
	Complete   bool   `json:"complete,omitempty"`
	Excess     string `json:"excess,omitempty"`
}

func Parse(filename string) (*TorrentInfo, error) {
	tor := &TorrentInfo{}

	var startIndex, endIndex = 0, len(filename)
	cleanName := strings.ReplaceAll(filename, "_", " ")
	for _, pattern := range patterns {
		match := findPatternMatch(pattern, cleanName)
		if match == nil {
			continue
		}

		rawStart, rawEnd := match[2], match[3]
		valueStart, valueEnd := match[4], match[5]
		if valueStart < 0 || valueEnd < 0 {
			valueStart, valueEnd = rawStart, rawEnd
		}
		if rawStart == 0 {
			startIndex = rawEnd
		} else if rawStart < endIndex {
			endIndex = rawStart
		}
		pattern.apply(tor, cleanName[valueStart:valueEnd])
	}

	if startIndex < 0 {
		startIndex = 0
	}
	if endIndex < 0 || endIndex > len(filename) {
		endIndex = len(filename)
	}
	if startIndex > endIndex {
		startIndex = 0
		endIndex = len(filename)
	}

	raw := strings.Split(filename[startIndex:endIndex], "(")[0]
	cleanName = strings.TrimSpace(raw)
	cleanName = strings.TrimLeft(cleanName, "- ")
	cleanName = strings.Trim(cleanName, ".-_ ")
	if strings.ContainsRune(cleanName, '.') && !strings.ContainsRune(cleanName, ' ') {
		cleanName = strings.ReplaceAll(cleanName, ".", " ")
	}
	cleanName = strings.ReplaceAll(cleanName, "_", " ")
	tor.Title = strings.TrimSpace(cleanName)
	augmentTorrentInfo(tor, filename)

	return tor, nil
}

func (info TorrentInfo) HasReleaseInfo() bool {
	return info.Title != "" ||
		info.Season != 0 ||
		info.Episode != 0 ||
		info.EpisodeEnd != 0 ||
		info.Part != 0 ||
		info.Year != 0 ||
		info.Resolution != "" ||
		info.Quality != "" ||
		info.Codec != "" ||
		info.HDR != "" ||
		info.Audio != "" ||
		info.Source != "" ||
		info.Group != "" ||
		info.Region != "" ||
		info.Extended ||
		info.Hardcoded ||
		info.Proper ||
		info.Repack ||
		info.Remastered ||
		info.Container != "" ||
		info.Widescreen ||
		info.Language != "" ||
		info.BitDepth != "" ||
		info.Edition != "" ||
		info.Unrated ||
		info.Size != "" ||
		info.ThreeD ||
		info.IMAX ||
		info.Complete ||
		info.Excess != ""
}

func findPatternMatch(pattern pattern, cleanName string) []int {
	if !pattern.last {
		return pattern.re.FindStringSubmatchIndex(cleanName)
	}

	matches := pattern.re.FindAllStringSubmatchIndex(cleanName, -1)
	if len(matches) == 0 {
		return nil
	}
	return matches[len(matches)-1]
}
