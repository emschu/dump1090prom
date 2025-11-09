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
	"math"
	"testing"
)

type stubSource struct{ p string }

func (s *stubSource) fetchMetrics() (*Dump1090AircraftJson, error) {
	return &Dump1090AircraftJson{}, nil
}
func (s *stubSource) fetchReceiverMetrics() (*Dump1090ReceiverJson, error) {
	return &Dump1090ReceiverJson{}, nil
}
func (s *stubSource) validateInput() bool { return true }
func (s *stubSource) getPath() string     { return s.p }

func TestGetDirection(t *testing.T) {
	cases := []struct {
		in   float64
		want string
	}{
		{0, "N"},
		{10, "N"},
		{22.5, "NNE"},
		{45, "NE"},
		{90, "E"},
		{135, "SE"},
		{180, "S"},
		{225, "SW"},
		{270, "W"},
		{315, "NW"},
		{359, "N"},
	}
	for _, c := range cases {
		got := getDirection(c.in)
		if got != c.want {
			t.Fatalf("direction(%v)=%s want %s", c.in, got, c.want)
		}
	}
}

func TestCalculateDistanceToAircraft(t *testing.T) {
	// Set config to enable distance calc with a receiver position
	lat := 52.5200 // Berlin
	lon := 13.4050
	source := &stubSource{p: "stub"}
	CONFIG = *newDump1090MetricCollectorConfig(source)
	CONFIG.ReceiverPositionLat = &lat
	CONFIG.ReceiverPositionLon = &lon

	d := newDump1090Metric(source)

	// Aircraft in Paris approx
	alat := 48.8566
	alon := 2.3522
	rec := &AircraftStateRecord{Lat: &alat, Lon: &alon}

	meters := d.calculateDistanceToAircraft(rec)
	if meters <= 0 {
		t.Fatalf("expected positive distance, got %f", meters)
	}
	// Expect roughly ~ 878km (great-circle) between Berlin and Paris; in meters about 878k
	if math.Abs(meters-878000) > 200000 { // allow big tolerance across sphere calc
		t.Fatalf("unexpected distance: %f m", meters)
	}

	// If disabled, distance should be 0
	CONFIG.IsDistanceCalculationEnabled = false
	if got := d.calculateDistanceToAircraft(rec); got != 0 {
		t.Fatalf("expected 0 distance when disabled, got %f", got)
	}
}
