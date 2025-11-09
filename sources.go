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
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

// MetricSource is an interface for fetching and validating metric data from a source related to aircraft and receiver metrics.
type MetricSource interface {
	fetchMetrics() (*Dump1090AircraftJson, error)
	fetchReceiverMetrics() (*Dump1090ReceiverJson, error)
	validateInput() bool
	getPath() string
}

type FileMetricSource struct {
	fileBasePath string
}

func (f *FileMetricSource) getPath() string {
	return f.fileBasePath
}

func (f *FileMetricSource) fetchReceiverMetrics() (*Dump1090ReceiverJson, error) {
	metrics, err := readReceiverJsonFromFile(fmt.Sprintf("%s/receiver.json", f.fileBasePath))
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

func (f *FileMetricSource) validateInput() bool {
	if !isDir(f.fileBasePath) {
		panic(fmt.Sprintf("Invalid input file '%s'", f.fileBasePath))
	}
	return true
}

func (f *FileMetricSource) fetchMetrics() (*Dump1090AircraftJson, error) {
	metrics, err := readAircraftJsonFromFile(fmt.Sprintf("%s/aircraft.json", f.fileBasePath))
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

func readAircraftJsonFromFile(inFile string) (*Dump1090AircraftJson, error) {
	return readFile[Dump1090AircraftJson](inFile)
}

func readReceiverJsonFromFile(inFile string) (*Dump1090ReceiverJson, error) {
	return readFile[Dump1090ReceiverJson](inFile)
}

func readFile[T payloads](inFile string) (*T, error) {
	// TODO make secure
	file, err := os.ReadFile(inFile)
	if err != nil {
		return nil, err
	}
	var msg *T
	jsonErr := json.Unmarshal(file, &msg)
	if jsonErr != nil {
		log.Fatalf("Error parsing receiver.json: %v", jsonErr)
	}
	return msg, nil
}

type URLMetricSource struct {
	url string
}

func (u *URLMetricSource) getPath() string {
	return u.url
}

func (u *URLMetricSource) fetchReceiverMetrics() (*Dump1090ReceiverJson, error) {
	return fetchPayload[Dump1090ReceiverJson](fmt.Sprintf("%s/receiver.json", u.url))
}

func (u *URLMetricSource) fetchMetrics() (*Dump1090AircraftJson, error) {
	return fetchPayload[Dump1090AircraftJson](fmt.Sprintf("%s/aircraft.json", u.url))
}

func (u *URLMetricSource) validateInput() bool {
	parsedURL, err := url.Parse(u.url)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		panic(fmt.Sprintf("Invalid URL: %s", u.url))
	}
	return true
}

var httpClient = &http.Client{
	Transport: &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).DialContext,
		MaxIdleConns:        5,
		MaxIdleConnsPerHost: 5,
		Proxy:               nil,
	},
	Timeout: time.Second * 15,
}

type payloads = interface {
	Dump1090AircraftJson | Dump1090ReceiverJson | Dump1090StatsJson
}

func fetchPayload[T payloads](targetURL string) (*T, error) {
	resp, err := httpClient.Get(targetURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch payload from URL '%s': %v", targetURL, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to close response body '%s': %v", targetURL, err)
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response '%s': %v", targetURL, err)
	}
	var payload T
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse JSON '%s': %v", targetURL, err)
	}
	return &payload, nil
}
