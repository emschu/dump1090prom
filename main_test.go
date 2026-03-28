// dump1090prom
// Copyright (C) 2025 emschu[aet]mailbox.org
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public
// License along with this program.
// If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func TestFindMetricSource_FileBasePath(t *testing.T) {
	tmp := t.TempDir()
	empty := ""
	src := findMetricSource(&empty, &tmp)
	if src == nil {
		t.Fatalf("expected non-nil source")
	}
	if src.getPath() != tmp {
		t.Fatalf("expected path %s, got %s", tmp, src.getPath())
	}
	if _, ok := src.(*FileMetricSource); !ok {
		t.Fatalf("expected *FileMetricSource, got %T", src)
	}
}

func TestFindMetricSource_URL(t *testing.T) {
	url := "http://example.com"
	empty := ""
	src := findMetricSource(&url, &empty)
	if src == nil {
		t.Fatalf("expected non-nil source")
	}
	if src.getPath() != url {
		t.Fatalf("expected url %s, got %s", url, src.getPath())
	}
	if _, ok := src.(*URLMetricSource); !ok {
		t.Fatalf("expected *URLMetricSource, got %T", src)
	}
}

func TestIsDir(t *testing.T) {
	tmp := t.TempDir()
	if !isDir(tmp) {
		t.Fatalf("expected %s to be a directory", tmp)
	}

	// Create a file
	f := filepath.Join(tmp, "file.txt")
	if err := os.WriteFile(f, []byte("x"), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if isDir(f) {
		t.Fatalf("expected %s to be a file, not directory", f)
	}

	// Non-existing path
	if isDir(filepath.Join(tmp, "doesnotexist")) {
		t.Fatalf("expected non-existing path to return false")
	}
}

func TestGlobalLabels(t *testing.T) {
	// Setup config with global labels
	empty := ""
	path := "example"
	source := findMetricSource(&empty, &path)
	CONFIG = *newDump1090MetricCollectorConfig(source)
	CONFIG.GlobalLabels = map[string]string{
		"location": "home",
		"env":      "test",
	}

	// Create and register the collector
	// Note: since prometheus.DefaultRegisterer is global, we should use a local registry for testing
	registry := prometheus.NewRegistry()
	metric := newDump1090Metric(source)
	registry.MustRegister(metric)

	// Create a test HTTP server
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	ts := httptest.NewServer(handler)
	defer ts.Close()

	// Make a request to the test server
	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("failed to get metrics: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	output := string(body)

	// Check if global labels are present in at least one metric
	// e.g. dump1090prom_aircraft_count{env="test",location="home"}
	if !strings.Contains(output, `location="home"`) {
		t.Errorf("expected global label location=\"home\" not found in output")
	}
	if !strings.Contains(output, `env="test"`) {
		t.Errorf("expected global label env=\"test\" not found in output")
	}

	// Verify it's on a specific metric too
	if !strings.Contains(output, `dump1090prom_aircraft_count{env="test",location="home"}`) {
		t.Errorf("metric dump1090prom_aircraft_count does not have expected global labels")
	}
}
