package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type benchmarkMetrics map[string]map[string]float64

type metricSample struct {
	total float64
	count int
}

var benchmarkSuffixPattern = regexp.MustCompile(`-\d+$`)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "usage: %s <baseline-bench.txt> <current-bench.txt>\n", os.Args[0])
		os.Exit(2)
	}
	if err := run(os.Args[1], os.Args[2]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(baselinePath string, currentPath string) error {
	baseline, err := parseBenchmarkFile(baselinePath)
	if err != nil {
		return err
	}
	current, err := parseBenchmarkFile(currentPath)
	if err != nil {
		return err
	}

	rows := comparisonRows(baseline, current)
	if len(rows) == 0 {
		return fmt.Errorf("no comparable benchmark metrics found")
	}

	fmt.Println("| benchmark | metric | baseline | current | delta |")
	fmt.Println("| --- | ---: | ---: | ---: | ---: |")
	for _, row := range rows {
		fmt.Printf("| %s | %s | %.2f | %.2f | %s |\n", row.name, row.metric, row.baseline, row.current, row.delta)
	}
	return nil
}

type comparisonRow struct {
	name     string
	metric   string
	baseline float64
	current  float64
	delta    string
}

func comparisonRows(baseline benchmarkMetrics, current benchmarkMetrics) []comparisonRow {
	names := make([]string, 0, len(baseline))
	for name := range baseline {
		if _, ok := current[name]; ok {
			names = append(names, name)
		}
	}
	sort.Strings(names)

	rows := make([]comparisonRow, 0)
	for _, name := range names {
		metrics := make([]string, 0, len(baseline[name]))
		for metric := range baseline[name] {
			if _, ok := current[name][metric]; ok {
				metrics = append(metrics, metric)
			}
		}
		sort.Strings(metrics)
		for _, metric := range metrics {
			base := baseline[name][metric]
			now := current[name][metric]
			rows = append(rows, comparisonRow{
				name:     name,
				metric:   metric,
				baseline: base,
				current:  now,
				delta:    formatDelta(base, now),
			})
		}
	}
	return rows
}

func parseBenchmarkFile(path string) (benchmarkMetrics, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open benchmark file %s: %w", path, err)
	}
	defer file.Close()

	samples := make(map[string]map[string]metricSample)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 4 || !strings.HasPrefix(fields[0], "Benchmark") {
			continue
		}
		name := benchmarkSuffixPattern.ReplaceAllString(fields[0], "")
		metrics := samples[name]
		if metrics == nil {
			metrics = make(map[string]metricSample)
			samples[name] = metrics
		}
		for index := 2; index+1 < len(fields); index += 2 {
			value, err := strconv.ParseFloat(fields[index], 64)
			if err != nil {
				continue
			}
			sample := metrics[fields[index+1]]
			sample.total += value
			sample.count++
			metrics[fields[index+1]] = sample
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan benchmark file %s: %w", path, err)
	}

	out := make(benchmarkMetrics, len(samples))
	for name, metrics := range samples {
		out[name] = make(map[string]float64, len(metrics))
		for metric, sample := range metrics {
			if sample.count == 0 {
				continue
			}
			out[name][metric] = sample.total / float64(sample.count)
		}
	}
	return out, nil
}

func formatDelta(baseline float64, current float64) string {
	if baseline == 0 {
		if current == 0 {
			return "+0.0%"
		}
		return "n/a"
	}
	return fmt.Sprintf("%+0.1f%%", (current-baseline)*100/baseline)
}
