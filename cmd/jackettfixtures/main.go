package main

import (
	"encoding/json"
	"errors"
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

type fixtureQuery struct {
	slug  string
	query string
	limit int
}

var curatedQueries = []fixtureQuery{
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
	{slug: "widow_bay_s01e04", query: "Widow Bay S01E04", limit: 11},
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
	Configured bool `json:"configured"`
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
	HasError    bool   `json:"has_error,omitempty"`
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
	client := &http.Client{Timeout: 30 * time.Second}
	if err := refreshFixtures(client, baseURL, apiKey, fixtureDir, curatedQueries, fetchFixture); err != nil {
		fail("refresh fixtures: %v", err)
	}
}

type fetchFixtureFunc func(*http.Client, string, string, string, int) (fixtureFile, error)

type stagedFixture struct {
	name    string
	payload []byte
	results int
}

func refreshFixtures(client *http.Client, baseURL, apiKey, dir string, queries []fixtureQuery, fetch fetchFixtureFunc) error {
	staged := make([]stagedFixture, 0, len(queries))
	for _, item := range queries {
		fixture, err := fetch(client, baseURL, apiKey, item.query, item.limit)
		if err != nil {
			return fmt.Errorf("fetch fixture %s: %w", item.slug, err)
		}
		payload, err := json.MarshalIndent(fixture, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal fixture %s: %w", item.slug, err)
		}
		staged = append(staged, stagedFixture{
			name:    item.slug + ".json",
			payload: append(payload, '\n'),
			results: len(fixture.Results),
		})
	}

	if err := replaceFixtureDir(dir, staged); err != nil {
		return err
	}
	for _, item := range staged {
		fmt.Fprintf(os.Stdout, "wrote %s (%d results)\n", filepath.Join(dir, item.name), item.results)
	}
	return nil
}

func replaceFixtureDir(dir string, fixtures []stagedFixture) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create fixture directory: %w", err)
	}
	dirInfo, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("inspect fixture directory: %w", err)
	}
	parent := filepath.Dir(dir)
	stage, err := os.MkdirTemp(parent, ".jackett-fixtures-stage-")
	if err != nil {
		return fmt.Errorf("create fixture staging directory: %w", err)
	}
	defer os.RemoveAll(stage)
	if err := os.Chmod(stage, dirInfo.Mode().Perm()); err != nil {
		return fmt.Errorf("set fixture staging permissions: %w", err)
	}

	if err := os.CopyFS(stage, os.DirFS(dir)); err != nil {
		return fmt.Errorf("stage fixture directory: %w", err)
	}
	if err := restoreCopiedModes(dir, stage); err != nil {
		return fmt.Errorf("preserve staged fixture permissions: %w", err)
	}
	if err := clearFixtureDir(stage); err != nil {
		return fmt.Errorf("clear staged fixtures: %w", err)
	}
	for _, fixture := range fixtures {
		if err := os.WriteFile(filepath.Join(stage, fixture.name), fixture.payload, 0o644); err != nil {
			return fmt.Errorf("stage fixture %s: %w", fixture.name, err)
		}
	}

	backup, err := os.MkdirTemp(parent, ".jackett-fixtures-backup-")
	if err != nil {
		return fmt.Errorf("reserve fixture backup: %w", err)
	}
	if err := os.Remove(backup); err != nil {
		return fmt.Errorf("prepare fixture backup: %w", err)
	}
	if err := os.Rename(dir, backup); err != nil {
		return fmt.Errorf("back up fixture directory: %w", err)
	}
	if err := os.Rename(stage, dir); err != nil {
		restoreErr := os.Rename(backup, dir)
		if restoreErr != nil {
			return fmt.Errorf("activate staged fixtures: %v; restore original fixtures: %w", err, restoreErr)
		}
		return fmt.Errorf("activate staged fixtures: %w", err)
	}
	if err := os.RemoveAll(backup); err != nil {
		return fmt.Errorf("remove fixture backup: %w", err)
	}
	return nil
}

func restoreCopiedModes(source, target string) error {
	return filepath.WalkDir(source, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		return os.Chmod(filepath.Join(target, relative), info.Mode())
	})
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

func fetchFixture(client *http.Client, baseURL, apiKey, queryValue string, limit int) (fixtureFile, error) {
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
		return fixtureFile{}, safeRequestError(err, requestURL, apiKey)
	}
	response, err := client.Do(request)
	if err != nil {
		return fixtureFile{}, safeRequestError(err, requestURL, apiKey)
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
			Configured: true,
		},
		Results:  sanitizeResults(raw.Results, limit),
		Indexers: sanitizeIndexers(raw.Indexers),
	}
	if len(fixture.Results) == 0 {
		return fixtureFile{}, fmt.Errorf("query %q returned no usable results", queryValue)
	}
	return fixture, nil
}

func safeRequestError(err error, requestURL *url.URL, apiKey string) error {
	message := err.Error()
	var requestErr *url.Error
	if errors.As(err, &requestErr) {
		message = requestErr.Err.Error()
	}
	sensitiveValues := []string{apiKey, url.QueryEscape(apiKey)}
	if requestURL.User != nil {
		encodedUserInfo := requestURL.User.String()
		rawUserInfo := requestURL.User.Username()
		if password, ok := requestURL.User.Password(); ok {
			rawUserInfo += ":" + password
		}
		sensitiveValues = append(sensitiveValues, encodedUserInfo+"@", rawUserInfo+"@")
	}
	sort.Slice(sensitiveValues, func(i, j int) bool {
		return len(sensitiveValues[i]) > len(sensitiveValues[j])
	})
	for _, sensitive := range sensitiveValues {
		if sensitive != "" {
			message = strings.ReplaceAll(message, sensitive, "[REDACTED]")
		}
	}
	return fmt.Errorf("request Jackett fixture from %s: %s", requestURL.Host, message)
}

func sanitizeResults(rawResults []rawResult, limit int) []fixtureResult {
	if limit <= 0 {
		limit = 25
	}
	results := make([]fixtureResult, 0, min(len(rawResults), limit))
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
		if len(results) == limit {
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
		fixture.HasError = item.Error != nil && strings.TrimSpace(*item.Error) != ""
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
