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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func writeJSONFile(t *testing.T, dir, name string, v any) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, b, 0o600); err != nil {
		t.Fatalf("failed to write %s: %v", p, err)
	}
	return p
}

func TestFileMetricSourceValidateAndFetch(t *testing.T) {
	tmp := t.TempDir()

	// minimal valid payloads
	receiver := Dump1090ReceiverJson{Version: "v", Refresh: 1, History: 1}
	lat := 10.0
	lon := 20.0
	receiver.Lat = &lat
	receiver.Lon = &lon
	writeJSONFile(t, tmp, "receiver.json", receiver)

	aircraft := Dump1090AircraftJson{Now: 1, Messages: 2, Aircraft: []AircraftStateRecord{{Hex: "abc123", Messages: 1, Seen: 0, RSSI: -3, Version: 2}}}
	writeJSONFile(t, tmp, "aircraft.json", aircraft)

	src := &FileMetricSource{fileBasePath: tmp}
	if !src.validateInput() {
		t.Fatalf("expected validateInput to be true")
	}

	rcv, err := src.fetchReceiverMetrics()
	if err != nil {
		t.Fatalf("fetchReceiverMetrics error: %v", err)
	}
	if rcv == nil || rcv.Lat == nil || *rcv.Lat != lat || *rcv.Lon != lon {
		t.Fatalf("unexpected receiver: %#v", rcv)
	}

	ac, err := src.fetchMetrics()
	if err != nil {
		t.Fatalf("fetchMetrics error: %v", err)
	}
	if ac == nil || len(ac.Aircraft) != 1 || ac.Messages != 2 {
		t.Fatalf("unexpected aircraft payload: %#v", ac)
	}
}

func TestFileMetricSourceValidatePanicsOnNonDir(t *testing.T) {
	tmp := t.TempDir()
	file := filepath.Join(tmp, "file.txt")
	if err := os.WriteFile(file, []byte("x"), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	src := &FileMetricSource{fileBasePath: file}
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic for non-directory path")
		}
	}()
	_ = src.validateInput()
}

func TestURLMetricSourceValidateAndFetch(t *testing.T) {
	// Create httptest server that serves valid JSON for both endpoints
	aircraft := Dump1090AircraftJson{Now: 5, Messages: 7}
	receiver := Dump1090ReceiverJson{Version: "x", Refresh: 1, History: 1}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/aircraft.json":
			_ = json.NewEncoder(w).Encode(aircraft)
		case "/receiver.json":
			_ = json.NewEncoder(w).Encode(receiver)
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	src := &URLMetricSource{url: ts.URL}
	if !src.validateInput() {
		t.Fatalf("expected validateInput to be true")
	}

	ac, err := src.fetchMetrics()
	if err != nil || ac == nil || ac.Now != 5 || ac.Messages != 7 {
		t.Fatalf("unexpected aircraft result: %#v, err=%v", ac, err)
	}

	rcv, err := src.fetchReceiverMetrics()
	if err != nil || rcv == nil || rcv.Version != "x" {
		t.Fatalf("unexpected receiver result: %#v, err=%v", rcv, err)
	}
}

func TestURLMetricSourceValidatePanicsOnInvalidURL(t *testing.T) {
	// Invalid scheme should panic
	src := &URLMetricSource{url: "ftp://example.com"}
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic for invalid URL scheme")
		}
	}()
	_ = src.validateInput()
}

func TestFetchPayloadInvalidJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Write invalid JSON
		_, _ = w.Write([]byte("{"))
	}))
	defer ts.Close()

	// Even with invalid JSON, we call fetchPayload directly to see error
	_, err := fetchPayload[Dump1090ReceiverJson](ts.URL)
	if err == nil {
		t.Fatalf("expected JSON unmarshal error, got nil")
	}
}

func TestFetchPayloadNetworkError(t *testing.T) {
	// assume port 1 is unavailable and will cause a dial error
	_, err := fetchPayload[Dump1090ReceiverJson]("http://127.0.0.1:1/receiver.json")
	if err == nil {
		t.Fatalf("expected network error, got nil")
	}
}

func TestReadJsonFromFile_NotFound(t *testing.T) {
	if _, err := readAircraftJsonFromFile("/non/existent/aircraft.json"); err == nil {
		t.Fatalf("expected error for missing aircraft.json")
	}
	if _, err := readReceiverJsonFromFile("/non/existent/receiver.json"); err == nil {
		t.Fatalf("expected error for missing receiver.json")
	}
}

func TestReadFile_InvalidJSON_ShouldReturnError(t *testing.T) {
	tmp := t.TempDir()
	file := filepath.Join(tmp, "invalid.json")
	_ = os.WriteFile(file, []byte("{"), 0o644)
	_, err := readFile[Dump1090AircraftJson](file)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
