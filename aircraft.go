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

// countryData represents a country entry from the countries.csv file
type countryData struct {
	Country            string
	CountryHexPrefixes []string
	CountryICAOCode    string
	HexRangeStart      string
	HexRangeEnd        string
	PossibleRegistry   string
}

// readAirlinesData reads the airlines data from the .csv file
// and returns a slice of airline structs
func readAirlinesData(fileContent string) ([]airline, error) {
	reader := csv.NewReader(bufio.NewReader(strings.NewReader(fileContent)))
	reader.Comma = ';'
	reader.LazyQuotes = true

	// skip first line
	_, err := reader.Read()
	if err == io.EOF {
		log.Fatal("airlines.csv content is empty")
	}

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

// readCountryData reads the country data from the countries.csv file
// and returns a slice of countryData structs
func readCountryData(fileContent string) ([]countryData, error) {
	reader := csv.NewReader(bufio.NewReader(strings.NewReader(fileContent)))
	reader.Comma = ';'
	reader.LazyQuotes = true

	var countries []countryData
	// skip header
	_, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading header: %w", err)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV: %w", err)
		}

		if len(record) < 6 {
			log.Printf("Invalid line in countries.csv detected: %v", record)
			continue
		}

		var prefixes []string
		for _, p := range strings.Split(strings.TrimSpace(record[1]), ",") {
			p = strings.ToUpper(strings.TrimSpace(p))
			if p != "" {
				prefixes = append(prefixes, p)
			}
		}

		country := countryData{
			Country:            strings.TrimSpace(record[0]),
			CountryHexPrefixes: prefixes,
			CountryICAOCode:    strings.TrimSpace(record[2]),
			HexRangeStart:      strings.TrimSpace(record[3]),
			HexRangeEnd:        strings.TrimSpace(record[4]),
			PossibleRegistry:   strings.TrimSpace(record[5]),
		}
		countries = append(countries, country)
	}

	return countries, nil
}

var icaoToFilter []string

// global variable to store airlines known
var airlines []airline
var icaoCountries []countryData

// the following mapping is used to enable full mapping of airline countries to ICAO country names
var airlineCountryMapping map[string]string = map[string]string{
	"Brunei":                                "Brunei Darussalam",
	"Burma":                                 "Myanmar",
	"Côte d'Ivoire":                         "Côte d Ivoire",
	"Democratic People's Republic of Korea": "North Korea",
	"Democratic Republic of Congo":          "Democratic Republic of the Congo",
	"Greenland":                             "Denmark",
	"Hong Kong":                             "Hong Kong SAR of China",
	"Iran":                                  "Iran, Islamic Republic of",
	"Iraqi Kurdistan":                       "Iraq",
	"Ivory Coast":                           "Côte d Ivoire",
	"Laos":                                  "Lao People's Democratic Republic",
	"Libya":                                 "Libyan Arab Jamahiriya",
	"Macedonia":                             "North Macedonia",
	"Moldova":                               "Republic of Moldova",
	"Netherlands":                           "Netherlands, Kingdom of the",
	"North Korea":                           "Democratic People's Republic of Korea",
	"Palau":                                 "United States",
	"Republic of the Congo":                 "Congo, Republic of the",
	"Russia":                                "Russian Federation",
	"Serbia":                                "Yugoslavia",
	"Somali Republic":                       "Somalia",
	"South Korea":                           "Republic of Korea",
	"Syria":                                 "Syrian Arab Republic",
	"São Tomé and Príncipe":                 "Sao Tome and Principe",
	"THAILAND":                              "Thailand",
	"Taiwan":                                "China",
	"Tanzania":                              "United Republic of Tanzania",
	"The Gambia":                            "Gambia",
	"Turkiye":                               "Turkey",
	"UAE":                                   "United Arab Emirates",
	"Vietnam":                               "Viet Nam",
}

func getAirlineLabel(flight string, hex string) *string {
	if flight == "" {
		return nil
	}
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

	var countryOfAirlineAndIcaoMatches = false
	countries := getIcaoCountries(hex)
	for _, match := range matchingAirlines {
		airlineCountry := match.Country
		// consider that wikipedia airlines have not consistent exactly icao conforming country names
		if mapped, ok := airlineCountryMapping[airlineCountry]; ok {
			airlineCountry = mapped
		}
		for _, country := range countries {
			if airlineCountry == country.Country {
				countryOfAirlineAndIcaoMatches = true
				break
			}
		}
	}

	if len(matchingAirlines) == 0 {
		// no match
		return nil
	}
	if CONFIG.IsVerbose && len(matchingAirlines) > 1 {
		log.Printf("multiple matches for %s: %#v", flight, matchingAirlines)
	}

	var ignoreAirlineMatch = false
	airline := matchingAirlines[0]

	if !countryOfAirlineAndIcaoMatches {
		countries := getIcaoCountries(hex)
		for _, country := range countries {
			if strings.HasPrefix(strings.ToUpper(flight), strings.ToUpper(country.CountryICAOCode)) &&
				strings.HasPrefix(strings.ToUpper(airline.ICAO), strings.ToUpper(country.CountryICAOCode)) {
				if CONFIG.IsVerbose {
					log.Printf("skipped airline detection. no country match (%s) for %s with hex %s: %#v", country.Country, flight, hex, matchingAirlines)
				}
				ignoreAirlineMatch = true
				break
			} else {
				countryLabels := getIcaoCountryLabel(hex)
				if CONFIG.IsVerbose {
					log.Printf("no country match (%s) for %s with hex %s: %#v", countryLabels, flight, hex, matchingAirlines)
				}
			}
		}
	}

	if ignoreAirlineMatch {
		flightAirLineLabelMap.Set(flight, "")
		return nil
	}

	flightAirLineMap.Set(flight, airline)
	label := fmt.Sprintf("%s - %s/%s %s", airline.ICAO, airline.Name, airline.Callsign, airline.Country)
	flightAirLineLabelMap.Set(flight, label)
	return &label
}

func getIcaoCountryLabel(hex string) string {
	if hex == "" {
		return ""
	}
	countries := getIcaoCountries(hex)
	if len(countries) == 0 || countries == nil {
		return ""
	}
	var countryLabels []string
	for _, country := range countries {
		countryLabels = append(countryLabels, country.Country)
	}
	return strings.Join(countryLabels, ", ")
}

func getIcaoCountries(hex string) []countryData {
	if hex == "" {
		return nil
	}
	if countries, exists := flightIcaoCountryMap.Get(hex); exists {
		return countries
	}
	hex = strings.ToUpper(hex)
	var matches []countryData
	for _, country := range icaoCountries {
		for _, prefix := range country.CountryHexPrefixes {
			if strings.HasPrefix(strings.ToUpper(hex), strings.ToUpper(prefix)) {
				matches = append(matches, country)
				break
			}
		}
	}
	flightIcaoCountryMap.Set(hex, matches)
	return matches
}
