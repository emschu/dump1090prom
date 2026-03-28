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
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var VERSION = "1.0.4"

//go:embed data/wikipedia-airlines.csv
var wikipediaAirlinesCsvFileContent string

//go:embed data/countries.csv
var countriesCsvFileContent string

var CONFIG dump1090MetricCollectorConfig

func main() {
	url := flag.String("base-url", "", "Base URL (path) to aircraft.json and receiver.json")
	file := flag.String("base-path", "", "Path to the directory where aircraft.json and receiver.json are located")
	lat := flag.Float64("lat", 0.0, "Custom latitude of the receiver")
	lon := flag.Float64("lon", 0.0, "Custom longitude of the receiver")
	host := flag.String("host", "127.0.0.1", "Host to listen on (default: 127.0.0.1)")
	port := flag.Int("port", 8080, "Port to listen on (default: 8080)")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	distanceCalc := flag.Bool("distance-calc", true, "Enable distance calculation to aircraft")
	airlineLabel := flag.Bool("airline-label", true, "Enable airline labelling")
	exposeFiles := flag.Bool("expose-files", true, "Expose original aircraft.json and receiver.json files")
	rollingMapSize := flag.Int("rolling-map-size", 1000, "Default size of the rolling map for caching")
	collector := flag.String("collector", "dump1090prom", "Constant 'collector' label value of this instance")
	globalLabels := flag.String("labels", "", "Global labels to add to all metrics (e.g., 'location=home,env=prod')")
	flag.Parse()
	log.Printf("Dump1090Prom - A bridge between the dump1090 JSON data and prometheus – Version %s", VERSION)
	source := findMetricSource(url, file)
	source.validateInput()

	CONFIG = *newDump1090MetricCollectorConfig(source)
	if verbose != nil {
		log.Printf("Verbose: %t", *verbose)
		CONFIG.IsVerbose = *verbose
	}
	if distanceCalc != nil {
		CONFIG.IsDistanceCalculationEnabled = *distanceCalc
	}
	if airlineLabel != nil {
		CONFIG.IsAirlineLabellingEnabled = *airlineLabel
	}
	if exposeFiles != nil {
		CONFIG.IsExposingOriginalFilesEnabled = *exposeFiles
	}
	if rollingMapSize != nil {
		CONFIG.RollingMapDefaultSize = *rollingMapSize
	}
	if collector != nil {
		CONFIG.Collector = *collector
	}
	CONFIG.GlobalLabels["collector"] = CONFIG.Collector
	if globalLabels != nil && *globalLabels != "" {
		labels := strings.Split(*globalLabels, ",")
		for _, label := range labels {
			parts := strings.SplitN(label, "=", 2)
			if len(parts) == 2 {
				CONFIG.GlobalLabels[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}
	}
	initRollingMaps()
	receiverData, err := source.fetchReceiverMetrics()
	if err != nil {
		log.Printf("Cannot open receiver.json: %v", err)
	}
	if receiverData != nil {
		marshal, err := json.Marshal(receiverData)
		if err != nil {
			log.Printf("Error marshalling receiver data: %v", err)
		} else {
			log.Printf("Receiver data: %s", marshal)
		}
	}
	if receiverData != nil {
		if receiverData.Lat != nil && receiverData.Lon != nil {
			log.Printf("Set receiver position for distance calculation from %s: Lat: %f, Lon: %f. You can override this with '-lat' and '-lon' arguments", CONFIG.Source.getPath(), *receiverData.Lat, *receiverData.Lon)
			CONFIG.ReceiverPositionLat = receiverData.Lat
			CONFIG.ReceiverPositionLon = receiverData.Lon
		}
	}
	if lat != nil && lon != nil && *lat != 0.0 && *lon != 0.0 {
		log.Printf("Set custom receiver position for distance calculation to: Lat: %f, Lon: %f", *lat, *lon)
		CONFIG.ReceiverPositionLat = lat
		CONFIG.ReceiverPositionLon = lon
	}
	if port != nil {
		CONFIG.Port = *port
	}
	if host != nil {
		CONFIG.Host = *host
	}

	// Read airline data
	airlineData, err := readAirlinesData(wikipediaAirlinesCsvFileContent)
	if err != nil {
		log.Fatalf("Cannot open airline data: %v", err)
	}
	airlines = airlineData
	for _, airline := range airlines {
		icao := strings.TrimSpace(airline.ICAO)
		if icao == "" {
			log.Printf("Invalid airline ICAO: %v", airline)
			continue
		}
		icaoToFilter = append(icaoToFilter, icao)
	}
	icaoCountries, err = readCountryData(countriesCsvFileContent)
	if err != nil {
		log.Fatalf("Cannot open country data: %v", err)
	}

	metric := newDump1090Metric(CONFIG.Source)
	metric.LoadAll()
	prometheus.MustRegister(metric)

	// Output metrics to a file for testing
	registry := prometheus.DefaultRegisterer.(*prometheus.Registry)
	metricFamilies, err := registry.Gather()
	if err != nil {
		log.Printf("Error gathering metrics: %v\n", err)
	} else {
		if CONFIG.IsVerbose {
			for _, mf := range metricFamilies {
				var metric *dump1090MetricItem
				for _, v := range metrics {
					metricName := fmt.Sprintf("%s_%s", v.Namespace, v.Name)
					if metricName == *mf.Name {
						metric = &v
						break
					}
				}
				if metric == nil {
					continue
				}
				log.Printf("Metric: %s – Type: %s – Labels: %s – Description: %s\n",
					*mf.Name, mf.Type.String(), strings.Join(metric.Labels, ", "), *mf.Help)
			}
		}
	}

	log.Printf("Starting server on %s:%d\n", CONFIG.Host, CONFIG.Port)

	if CONFIG.IsExposingOriginalFilesEnabled {
		http.HandleFunc("/aircraft.json", func(w http.ResponseWriter, r *http.Request) {
			// Write the string content
			fetchMetrics, err := CONFIG.Source.fetchMetrics()
			writeFileAsJson(w, err, fetchMetrics)
		})
		http.HandleFunc("/receiver.json", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			// Write the string content
			receiverMetrics, err := CONFIG.Source.fetchReceiverMetrics()
			writeFileAsJson(w, err, receiverMetrics)
		})
	}

	http.Handle("/metrics", promhttp.Handler())
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", CONFIG.Host, CONFIG.Port), nil)
	if err != nil {
		log.Fatalf("failed to start server on %s:%d: %v", CONFIG.Host, CONFIG.Port, err)
	}
}

func writeFileAsJson(w http.ResponseWriter, err error, fetchMetrics any) {
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		log.Printf("Error fetching metrics: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	marshal, err := json.MarshalIndent(fetchMetrics, "", "  ")
	if err != nil {
		log.Printf("Error marshalling metrics: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(marshal)
	if err != nil {
		log.Printf("Error writing metrics: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// findMetricSource returns a MetricSource based on the provided arguments
func findMetricSource(url *string, fileBasePath *string) MetricSource {
	if url != nil && len(*url) > 0 {
		log.Printf("URL: %s", *url)
	}
	if fileBasePath != nil && len(*fileBasePath) > 0 {
		log.Printf("File: %s", *fileBasePath)
	}

	if (url == nil || *url == "") && (fileBasePath == nil || *fileBasePath == "") {
		log.Fatal("No input fileBasePath or URL provided")
	}
	if (url != nil && *url != "") && (fileBasePath != nil && *fileBasePath != "") {
		log.Fatal("Both input fileBasePath and URL provided, use only one")
	}
	var source MetricSource
	if url != nil && *url != "" {
		source = &URLMetricSource{*url}
	}
	if fileBasePath != nil && *fileBasePath != "" {
		source = &FileMetricSource{*fileBasePath}
	}
	return source
}
