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
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strings"
)

// airline represents an airline from the openflights-airlines.dat file
// Only includes the required fields: name, alias, iata, icao, callsign, country
type airline struct {
	Name     string
	IATA     string
	ICAO     string
	Callsign string
	Country  string
	Comments string
}

// readAirlinesData reads the airlines data from the .dat file
// and returns a slice of airline structs
func readAirlinesData(fileContent string) ([]airline, error) {
	reader := csv.NewReader(bufio.NewReader(strings.NewReader(fileContent)))
	reader.Comma = ';'
	reader.LazyQuotes = true

	var airlines []airline
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV: %w", err)
		}

		if len(record) != 6 {
			log.Printf("Invalid line in airlines.csv detected: %v", record)
			continue
		}

		iata := strings.TrimSpace(record[0])
		icao := strings.TrimSpace(record[1])
		name := strings.TrimSpace(record[2])
		callSign := strings.TrimSpace(record[3])
		countryOrRegion := strings.TrimSpace(record[4])
		comments := strings.TrimSpace(record[5])

		// at this moment we only care about airlines with ICAO codes
		if icao == "" {
			continue
		}

		airline := airline{
			Name:     name,
			IATA:     iata,
			ICAO:     icao,
			Callsign: callSign,
			Country:  countryOrRegion,
			Comments: comments,
		}
		airlines = append(airlines, airline)
	}

	return airlines, nil
}

var icaoToFilter []string

// global variable to store airlines known
var airlines []airline

func getAirlineLabel(flight string) *string {
	if flight == "" {
		return nil
	}
	// TODO make nice
	airlineLabel, exists := flightAirLineLabelMap.Get(flight)
	if exists {
		return &airlineLabel
	}
	var matchingAirlines []airline
	var prefixes []string
	for _, singleAirline := range airlines {
		prefix := strings.TrimSpace(singleAirline.ICAO)
		if len(prefix) == 0 {
			continue
		}

		if strings.HasPrefix(flight, prefix) {
			prefixes = append(prefixes, prefix)
			matchingAirlines = append(matchingAirlines, singleAirline)
		}
	}

	if len(prefixes) == 0 || len(matchingAirlines) == 0 {
		// no match
		return nil
	}
	if CONFIG.IsVerbose && len(matchingAirlines) > 1 {
		log.Printf("multiple matches for %s: %#v", flight, matchingAirlines)
	}
	for _, airline := range matchingAirlines {
		flightAirLineMap.Set(flight, airline)
		label := fmt.Sprintf("%s - %s/%s %s", airline.ICAO, airline.Name, airline.Callsign, airline.Country)
		flightAirLineLabelMap.Set(flight, label)
		return &label
	}
	return nil
}
