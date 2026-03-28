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

import "testing"

func TestGetAirlineLabel_MatchAndNoMatch(t *testing.T) {
	// Prepare global state
	airlines = []airline{
		{Name: "My Test Air", IATA: "", ICAO: "TST", Callsign: "TEST", Country: "Nowhere"},
	}
	// Ensure rolling maps are initialized with sane size
	flightAirLineMap = NewRollingAirLineMap(10)
	flightAirLineLabelMap = NewRollingAirLineLabelMap(10)
	flightIcaoCountryMap = NewRollingAirLineCountryMap(10)

	// Matching prefix
	if lbl := getAirlineLabel("TST1234", ""); lbl == nil {
		t.Fatalf("expected label for matching ICAO prefix")
	} else if got := *lbl; got == "" {
		t.Fatalf("expected non-empty label")
	}

	// No match
	if lbl := getAirlineLabel("ABC1234", ""); lbl != nil {
		t.Fatalf("expected nil for non-matching prefix, got %v", *lbl)
	}

	// Empty flight code
	if lbl := getAirlineLabel("", ""); lbl != nil {
		t.Fatalf("expected nil for empty flight code")
	}

	label := getAirlineLabel("ABCDE", "")
	if label != nil {
		t.Fatalf("Expected no airline match on 5-letter code, got %v", *label)
	}
}

func TestReadAirlinesData_File(t *testing.T) {
	airlinesList, err := readAirlinesData(wikipediaAirlinesCsvFileContent)
	if err != nil {
		t.Fatalf("readAirlinesData error: %v", err)
	}
	if len(airlinesList) == 0 {
		t.Fatalf("expected some airlines parsed from CSV")
	}
	// ensure all entries have ICAO non-empty as function filters empties
	for _, a := range airlinesList {
		if a.ICAO == "" {
			t.Fatalf("unexpected empty ICAO in parsed airlines: %#v", a)
		}
	}
}
