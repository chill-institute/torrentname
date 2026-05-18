package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const fixtureDir = "testdata/jackett"

var curatedQueries = []struct {
	slug  string
	query string
}{
	{slug: "1080p_bluray_x264", query: "1080p BluRay x264"},
	{slug: "2160p_hevc_hdr", query: "2160p HEVC HDR"},
	{slug: "anime_1080p", query: "Anime 1080p"},
	{slug: "complete_season_1080p", query: "Complete Season 1080p"},
	{slug: "dual_audio_1080p", query: "Dual Audio 1080p"},
	{slug: "dv_hdr_hevc", query: "DV HDR HEVC"},
	{slug: "hdrip_xvid_ac3", query: "HDRip XviD AC3"},
	{slug: "part_1_1080p", query: "Part 1 1080p"},
	{slug: "remux_truehd_atmos", query: "REMUX TrueHD Atmos"},
	{slug: "s01e01_720p_hdtv", query: "S01E01 720p HDTV"},
	{slug: "web_dl_ddp5_1", query: "WEB-DL DDP5.1"},
	{slug: "x265_aac_720p", query: "x265 AAC 720p"},
}

type fixtureFile struct {
	Query     string           `json:"query"`
	FetchedAt string           `json:"fetched_at"`
	Source    fixtureSource    `json:"source"`
	Results   []fixtureResult  `json:"results"`
	Indexers  []fixtureIndexer `json:"indexers"`
}

type fixtureSource struct {
	BaseURL    string `json:"base_url"`
	Configured bool   `json:"configured"`
}

type fixtureResult struct {
	Title       string `json:"title"`
	Tracker     string `json:"tracker"`
	TrackerID   string `json:"tracker_id"`
	Category    string `json:"category"`
	Size        int64  `json:"size"`
	Seeders     int    `json:"seeders"`
	Peers       int    `json:"peers"`
	PublishDate string `json:"publish_date,omitempty"`
	IMDb        string `json:"imdb,omitempty"`
	HasLink     bool   `json:"has_link,omitempty"`
	HasDetails  bool   `json:"has_details,omitempty"`
	HasMagnet   bool   `json:"has_magnet,omitempty"`
}

type fixtureIndexer struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Status      int    `json:"status"`
	Results     int    `json:"results"`
	ElapsedTime int    `json:"elapsed_time_ms"`
	Error       string `json:"error,omitempty"`
}

type rawResponse struct {
	Results  []rawResult  `json:"Results"`
	Indexers []rawIndexer `json:"Indexers"`
}

type rawResult struct {
	Title        string          `json:"Title"`
	Tracker      string          `json:"Tracker"`
	TrackerID    string          `json:"TrackerId"`
	CategoryDesc string          `json:"CategoryDesc"`
	Size         json.RawMessage `json:"Size"`
	Seeders      int             `json:"Seeders"`
	Peers        int             `json:"Peers"`
	PublishDate  string          `json:"PublishDate"`
	IMDb         json.RawMessage `json:"Imdb"`
	Link         *string         `json:"Link"`
	Details      *string         `json:"Details"`
	MagnetURI    *string         `json:"MagnetUri"`
}

type rawIndexer struct {
	ID          string  `json:"ID"`
	Name        string  `json:"Name"`
	Status      int     `json:"Status"`
	Results     int     `json:"Results"`
	Error       *string `json:"Error"`
	ElapsedTime int     `json:"ElapsedTime"`
}

func main() {
	baseURL := strings.TrimSpace(os.Getenv("JACKETT_BASE_URL"))
	if baseURL == "" {
		baseURL = "http://localhost:9117"
	}
	apiKey := strings.TrimSpace(os.Getenv("JACKETT_API_KEY"))
	if apiKey == "" {
		fail("JACKETT_API_KEY is required")
	}
	if err := os.MkdirAll(fixtureDir, 0o755); err != nil {
		fail("create fixture directory: %v", err)
	}
	if err := clearFixtureDir(fixtureDir); err != nil {
		fail("clear fixture directory: %v", err)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	for _, item := range curatedQueries {
		fixture, err := fetchFixture(client, baseURL, apiKey, item.query)
		if err != nil {
			fail("fetch fixture %s: %v", item.slug, err)
		}

		target := filepath.Join(fixtureDir, item.slug+".json")
		payload, err := json.MarshalIndent(fixture, "", "  ")
		if err != nil {
			fail("marshal fixture %s: %v", item.slug, err)
		}
		payload = append(payload, '\n')
		if err := os.WriteFile(target, payload, 0o644); err != nil {
			fail("write fixture %s: %v", item.slug, err)
		}
		fmt.Fprintf(os.Stdout, "wrote %s (%d results)\n", target, len(fixture.Results))
	}
}

func clearFixtureDir(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		if err := os.Remove(filepath.Join(dir, entry.Name())); err != nil {
			return err
		}
	}
	return nil
}

func fetchFixture(client *http.Client, baseURL, apiKey, queryValue string) (fixtureFile, error) {
	requestURL, err := url.Parse(strings.TrimRight(baseURL, "/") + "/api/v2.0/indexers/all/results")
	if err != nil {
		return fixtureFile{}, err
	}
	params := requestURL.Query()
	params.Set("apikey", apiKey)
	params.Set("configured", "true")
	params.Set("Query", queryValue)
	requestURL.RawQuery = params.Encode()

	request, err := http.NewRequest(http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return fixtureFile{}, err
	}
	response, err := client.Do(request)
	if err != nil {
		return fixtureFile{}, err
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fixtureFile{}, fmt.Errorf("unexpected status %d", response.StatusCode)
	}

	var raw rawResponse
	if err := json.NewDecoder(response.Body).Decode(&raw); err != nil {
		return fixtureFile{}, err
	}

	fixture := fixtureFile{
		Query:     queryValue,
		FetchedAt: time.Now().UTC().Format(time.RFC3339),
		Source: fixtureSource{
			BaseURL:    strings.TrimRight(baseURL, "/"),
			Configured: true,
		},
		Results:  sanitizeResults(raw.Results),
		Indexers: sanitizeIndexers(raw.Indexers),
	}
	if len(fixture.Results) == 0 {
		return fixtureFile{}, fmt.Errorf("query %q returned no usable results", queryValue)
	}
	return fixture, nil
}

func sanitizeResults(rawResults []rawResult) []fixtureResult {
	results := make([]fixtureResult, 0, min(len(rawResults), 25))
	for _, item := range rawResults {
		title := strings.TrimSpace(item.Title)
		if title == "" {
			continue
		}
		trackerID := strings.TrimSpace(item.TrackerID)
		if trackerID == "" {
			trackerID = strings.ToLower(strings.TrimSpace(item.Tracker))
		}
		if trackerID == "" {
			continue
		}
		results = append(results, fixtureResult{
			Title:       title,
			Tracker:     strings.TrimSpace(item.Tracker),
			TrackerID:   trackerID,
			Category:    strings.TrimSpace(item.CategoryDesc),
			Size:        parseSize(item.Size),
			Seeders:     item.Seeders,
			Peers:       item.Peers,
			PublishDate: normalizeTimestamp(item.PublishDate),
			IMDb:        parseIMDb(item.IMDb),
			HasLink:     item.Link != nil && strings.TrimSpace(*item.Link) != "",
			HasDetails:  item.Details != nil && strings.TrimSpace(*item.Details) != "",
			HasMagnet:   item.MagnetURI != nil && strings.TrimSpace(*item.MagnetURI) != "",
		})
		if len(results) == 25 {
			break
		}
	}
	return results
}

func sanitizeIndexers(rawIndexers []rawIndexer) []fixtureIndexer {
	indexers := make([]fixtureIndexer, 0, len(rawIndexers))
	for _, item := range rawIndexers {
		fixture := fixtureIndexer{
			ID:          strings.TrimSpace(item.ID),
			Name:        strings.TrimSpace(item.Name),
			Status:      item.Status,
			Results:     item.Results,
			ElapsedTime: item.ElapsedTime,
		}
		if item.Error != nil {
			fixture.Error = strings.TrimSpace(*item.Error)
		}
		indexers = append(indexers, fixture)
	}
	sort.Slice(indexers, func(i, j int) bool {
		return indexers[i].ID < indexers[j].ID
	})
	return indexers
}

func parseSize(raw json.RawMessage) int64 {
	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "" || trimmed == "null" {
		return 0
	}
	trimmed = strings.Trim(trimmed, `"`)
	value, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		return 0
	}
	return int64(value)
}

func parseIMDb(raw json.RawMessage) string {
	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "" || trimmed == "null" {
		return ""
	}
	trimmed = strings.Trim(trimmed, `"`)
	trimmed = strings.TrimSpace(trimmed)
	if trimmed == "" {
		return ""
	}
	if strings.HasPrefix(strings.ToLower(trimmed), "tt") {
		trimmed = trimmed[2:]
	}
	id, err := strconv.ParseInt(trimmed, 10, 64)
	if err != nil || id <= 0 {
		return ""
	}
	return fmt.Sprintf("tt%07d", id)
}

func normalizeTimestamp(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ""
	}
	parsed, err := time.Parse(time.RFC3339Nano, trimmed)
	if err != nil {
		return ""
	}
	return parsed.UTC().Format(time.RFC3339)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
