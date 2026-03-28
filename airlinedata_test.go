package main

import (
	"log"
	"maps"
	"os"
	"slices"
	"testing"
)

func TestReadCountryData(t *testing.T) {
	content, err := os.ReadFile("data/countries.csv")
	if err != nil {
		t.Fatalf("failed to read data/countries.csv: %v", err)
	}

	countries, err := readCountryData(string(content))
	if err != nil {
		t.Fatalf("readCountryData failed: %v", err)
	}

	if len(countries) == 0 {
		t.Error("expected at least one country, got zero")
	}

	// Verify the first entry (Afghanistan)
	if len(countries) > 0 {
		first := countries[0]
		if first.Country != "Afghanistan" {
			t.Errorf("expected Afghanistan, got %s", first.Country)
		}
		if !slices.Contains(first.CountryHexPrefixes, "700") {
			t.Errorf("expected 700 to be in %v", first.CountryHexPrefixes)
		}
	}
}

func TestGetAirlineCountry(t *testing.T) {
	// Setup icaoCountries if not already populated (it should be via main's init or read in tests)
	if len(icaoCountries) == 0 {
		var err error
		icaoCountries, err = readCountryData(countriesCsvFileContent)
		if err != nil {
			t.Fatalf("failed to read country data: %v", err)
		}
	}
	initRollingMaps()

	tests := []struct {
		icao     string
		expected []string // names of matching countries
	}{
		{"0a0123", []string{"Algeria"}},
		{"a00001", []string{"United States"}},
		{"afffff", []string{"United States"}},
		{"FFFFFF", nil},
		{"", nil},
	}

	for _, tt := range tests {
		got := getIcaoCountries(tt.icao)
		if tt.expected == nil {
			if got != nil {
				t.Errorf("getAirlineCountry(%s) = %v; want nil", tt.icao, got)
			}
		} else {
			if got == nil {
				t.Errorf("getAirlineCountry(%s) = nil; want %v", tt.icao, tt.expected)
				continue
			}
			var gotNames []string
			for _, c := range got {
				gotNames = append(gotNames, c.Country)
			}
			for _, wantName := range tt.expected {
				if !slices.Contains(gotNames, wantName) {
					t.Errorf("getAirlineCountry(%s) = %v; want to contain %s", tt.icao, gotNames, wantName)
				}
			}
		}
	}
}

func TestCountriesDoNotOverlap(t *testing.T) {
	countries, err := readCountryData(countriesCsvFileContent)
	if err != nil {
		t.Fatalf("readCountryData failed: %v", err)
	}

	for i, c1 := range countries {
		for j, c2 := range countries {
			if i == j {
				continue
			}
			// Check if ranges overlap
			if c1.HexRangeStart <= c2.HexRangeEnd && c2.HexRangeStart <= c1.HexRangeEnd {
				log.Fatalf("Countries %s and %s overlap", c1.Country, c2.Country)
			}
		}
	}
}

func TestAirLineDataCountries(t *testing.T) {
	countries, err := readCountryData(countriesCsvFileContent)
	if err != nil {
		t.Fatalf("readCountryData failed: %v", err)
	}

	airlines, err := readAirlinesData(wikipediaAirlinesCsvFileContent)
	if err != nil {
		t.Fatalf("readAirlinesData failed: %v", err)
	}

	validCountries := make(map[string]bool)
	for _, country := range countries {
		var countryKey = country.Country
		validCountries[countryKey] = true
	}
	for v, _ := range airlineCountryMapping {
		validCountries[v] = true
	}

	// data is missing atm
	knownLabelsToSkip := []string{
		"Aruba",
		"Bermuda",
		"British Virgin Islands",
		"Cayman Islands",
		"Faroe Islands",
		"French Guiana",
		"French Polynesia",
		"Hong Kong SAR of China",
		"Jersey",
		"Macao",
		"Montenegro",
		"Montserrat",
		"Netherlands Antilles",
		"North Macedonia",
		"Saint Kitts and Nevis",
		"Turks and Caicos Islands",
		"Norway, Sweden and San Diego",
		"Sweden, Denmark and Norway",
	}

	missingCountries := make(map[string]bool)
	for _, airline := range airlines {
		if len(airline.Country) > 0 && !validCountries[airline.Country] && !slices.Contains(knownLabelsToSkip, airline.Country) {
			missingCountries[airline.Country] = true
		}
	}

	if len(missingCountries) > 0 {
		t.Errorf("The following countries from wikipedia-airlines.csv are missing in countries.csv:")
		if len(missingCountries) == 0 {
			return
		}

		sorted := slices.Sorted(maps.Keys(missingCountries))
		for _, v := range sorted {
			t.Errorf("  - %v", v)
		}
	}
}

func TestAirlineMissmatch(t *testing.T) {
	initRollingMaps()
	var err error
	airlines, err = readAirlinesData(wikipediaAirlinesCsvFileContent)
	if err != nil {
		t.Fatalf("readAirlinesData failed: %v", err)
	}
	icaoCountries, err = readCountryData(countriesCsvFileContent)
	if err != nil {
		t.Fatalf("readCountryData failed: %v", err)
	}

	falsePositiveFlights := []string{
		"DEA99",
		"DMD99",
		"DKG99",
		"DCO99",
		"DEL99",
		"DFA99",
		"DGO99",
	}

	for _, falsePositive := range falsePositiveFlights {
		result := getAirlineLabel(falsePositive, "3c123")
		if result != nil {
			t.Fatalf("expected nil for non-matching airline '%s', got %v", falsePositive, *result)
		}
	}
}
