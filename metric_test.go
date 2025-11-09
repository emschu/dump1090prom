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
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func checkMetricPresence(t *testing.T, output string, metricName string) {
	if !strings.Contains(output, Dump1090Namespace+"_"+metricName) {
		t.Errorf("Expected metric %s not found in output", metricName)
	}
	if !strings.Contains(output, "# HELP "+Dump1090Namespace+"_"+metricName) {
		t.Errorf("Help text for metric %s not found", metricName)
	}
	if !strings.Contains(output, "# TYPE "+Dump1090Namespace+"_"+metricName) {
		t.Errorf("Type information for metric %s not found", metricName)
	}
}

// Test that metric generation works
func TestFunction(t *testing.T) {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run metric_test.go <input_file>\n")
		return
	}

	emptyString := ""
	s := "example"
	source := findMetricSource(&emptyString, &s)
	// Set the current input file
	CONFIG = *newDump1090MetricCollectorConfig(source)

	// Register the metric
	metric := newDump1090Metric(source)
	prometheus.MustRegister(metric)

	// Create a test HTTP server
	handler := promhttp.Handler()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}))
	defer server.Close()

	// Make a request to the test server
	resp, err := http.Get(server.URL)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Read and print the response
	buf := make([]byte, 1024*1024) // 1MB buffer
	n, err := resp.Body.Read(buf)
	if err != nil && err.Error() != "EOF" {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	metricsOutput := string(buf[:n])
	if strings.Contains(metricsOutput, "NaN") {
		t.Errorf("NaN found in metrics output")
	}
	if strings.Contains(metricsOutput, "undefined") {
		t.Errorf("'undefined' found in metrics output")
	}

	metric.itemMap.Range(func(key, value interface{}) bool {
		item := value.(*dump1090MetricItem)
		checkMetricPresence(t, metricsOutput, item.Name)
		return true
	})
}
