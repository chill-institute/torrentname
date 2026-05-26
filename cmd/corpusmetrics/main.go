package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/chill-institute/torrentname"
)

type fixture struct {
	Results []struct {
		Title string `json:"title"`
	} `json:"results"`
}

type fieldCounter struct {
	name string
	set  func(torrentname.TorrentInfo) bool
}

type fieldThreshold struct {
	name    string
	percent float64
}

type thresholdFlags []fieldThreshold

var fieldCounters = []fieldCounter{
	{name: "title", set: func(info torrentname.TorrentInfo) bool { return info.Title != "" }},
	{name: "season", set: func(info torrentname.TorrentInfo) bool { return info.Season != 0 }},
	{name: "episode", set: func(info torrentname.TorrentInfo) bool { return info.Episode != 0 }},
	{name: "year", set: func(info torrentname.TorrentInfo) bool { return info.Year != 0 }},
	{name: "resolution", set: func(info torrentname.TorrentInfo) bool { return info.Resolution != "" }},
	{name: "quality", set: func(info torrentname.TorrentInfo) bool { return info.Quality != "" }},
	{name: "codec", set: func(info torrentname.TorrentInfo) bool { return info.Codec != "" }},
	{name: "hdr", set: func(info torrentname.TorrentInfo) bool { return info.HDR != "" }},
	{name: "audio", set: func(info torrentname.TorrentInfo) bool { return info.Audio != "" }},
	{name: "source", set: func(info torrentname.TorrentInfo) bool { return info.Source != "" }},
	{name: "group", set: func(info torrentname.TorrentInfo) bool { return info.Group != "" }},
	{name: "edition", set: func(info torrentname.TorrentInfo) bool { return info.Edition != "" }},
	{name: "bit_depth", set: func(info torrentname.TorrentInfo) bool { return info.BitDepth != "" }},
	{name: "complete", set: func(info torrentname.TorrentInfo) bool { return info.Complete }},
}

func main() {
	dir := flag.String("dir", filepath.Join("testdata", "jackett"), "fixture directory")
	var minimums thresholdFlags
	flag.Var(&minimums, "min", "minimum field coverage as field=percent; repeatable")
	flag.Parse()

	if err := run(*dir, minimums...); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(dir string, minimums ...fieldThreshold) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		return fmt.Errorf("glob fixtures: %w", err)
	}
	if len(files) == 0 {
		return fmt.Errorf("no fixture files found in %s", dir)
	}
	sort.Strings(files)

	counts := make(map[string]int, len(fieldCounters))
	filesRead := 0
	titles := 0
	releaseInfo := 0
	parseErrors := 0

	for _, path := range files {
		entries, err := loadTitles(path)
		if err != nil {
			return err
		}
		filesRead++
		for _, title := range entries {
			title = strings.TrimSpace(title)
			if title == "" {
				continue
			}
			titles++
			info, err := torrentname.Parse(title)
			if err != nil {
				parseErrors++
				continue
			}
			if info.HasReleaseInfo() {
				releaseInfo++
			}
			for _, counter := range fieldCounters {
				if counter.set(*info) {
					counts[counter.name]++
				}
			}
		}
	}

	fmt.Printf("jackett corpus: files=%d titles=%d parsed=%d errors=%d\n", filesRead, titles, titles-parseErrors, parseErrors)
	fmt.Printf("release_info: %d/%d %s\n", releaseInfo, titles, percent(releaseInfo, titles))
	fmt.Println("fields:")
	for _, counter := range fieldCounters {
		fmt.Printf("  %-10s %4d/%-4d %s\n", counter.name+":", counts[counter.name], titles, percent(counts[counter.name], titles))
	}
	return checkThresholds(minimums, counts, releaseInfo, titles)
}

func loadTitles(path string) ([]string, error) {
	payload, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read fixture %s: %w", path, err)
	}
	var fixture fixture
	if err := json.Unmarshal(payload, &fixture); err != nil {
		return nil, fmt.Errorf("decode fixture %s: %w", path, err)
	}

	titles := make([]string, 0, len(fixture.Results))
	for _, result := range fixture.Results {
		titles = append(titles, result.Title)
	}
	return titles, nil
}

func percent(value int, total int) string {
	if total == 0 {
		return "(0.0%)"
	}
	return fmt.Sprintf("(%.1f%%)", percentValue(value, total))
}

func percentValue(value int, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(value) * 100 / float64(total)
}

func checkThresholds(minimums []fieldThreshold, counts map[string]int, releaseInfo int, total int) error {
	for _, minimum := range minimums {
		count := counts[minimum.name]
		if minimum.name == "release_info" {
			count = releaseInfo
		}
		actual := percentValue(count, total)
		if actual < minimum.percent {
			return fmt.Errorf("%s coverage %.1f%% is below minimum %.1f%%", minimum.name, actual, minimum.percent)
		}
	}
	return nil
}

func (flags *thresholdFlags) Set(raw string) error {
	name, value, ok := strings.Cut(raw, "=")
	if !ok {
		return fmt.Errorf("minimum %q must use field=percent", raw)
	}
	name = strings.TrimSpace(name)
	if !knownThresholdField(name) {
		return fmt.Errorf("unknown threshold field %q", name)
	}
	parsed, err := strconv.ParseFloat(strings.TrimSuffix(strings.TrimSpace(value), "%"), 64)
	if err != nil {
		return fmt.Errorf("parse threshold %q: %w", raw, err)
	}
	if parsed < 0 || parsed > 100 {
		return fmt.Errorf("threshold %q must be between 0 and 100", raw)
	}
	*flags = append(*flags, fieldThreshold{name: name, percent: parsed})
	return nil
}

func (flags *thresholdFlags) String() string {
	values := make([]string, 0, len(*flags))
	for _, threshold := range *flags {
		values = append(values, fmt.Sprintf("%s=%.1f", threshold.name, threshold.percent))
	}
	return strings.Join(values, ",")
}

func knownThresholdField(name string) bool {
	if name == "release_info" {
		return true
	}
	for _, counter := range fieldCounters {
		if counter.name == name {
			return true
		}
	}
	return false
}
