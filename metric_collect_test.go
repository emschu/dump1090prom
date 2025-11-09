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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type stubRichSource struct{ p string }

func (s *stubRichSource) getPath() string     { return s.p }
func (s *stubRichSource) validateInput() bool { return true }
func (s *stubRichSource) fetchReceiverMetrics() (*Dump1090ReceiverJson, error) {
	return &Dump1090ReceiverJson{}, nil
}
func (s *stubRichSource) fetchMetrics() (*Dump1090AircraftJson, error) {
	lat := 48.0
	lon := 11.0
	altb := 30000.0
	altg := 30500.0
	gs := 420.0
	ias := 250.0
	tas := 300.0
	mach := 0.78
	track := 270.0
	trackRate := 1.2
	roll := 2.0
	magHeading := 268.0
	baroRate := -500.0
	geomRate := -450.0
	navQNH := 1013.0
	navAlt := 32000.0
	navHeading := 265.0
	trueHeading := 266.0
	seenPos := 2.0
	nacp := 8
	nacv := 2
	nic := 7
	nicBaro := 1
	sil := 2
	gva := 1
	sda := 2
	rc := 75

	return &Dump1090AircraftJson{
		Now:      12345,
		Messages: 3,
		Aircraft: []AircraftStateRecord{{
			Hex:            "ABCDEF",
			Flight:         "TST123",
			AltBaro:        &altb,
			AltGeom:        &altg,
			GS:             &gs,
			IAS:            &ias,
			TAS:            &tas,
			Mach:           &mach,
			Track:          &track,
			TrackRate:      &trackRate,
			Roll:           &roll,
			MagHeading:     &magHeading,
			BaroRate:       &baroRate,
			GeomRate:       &geomRate,
			NavQNH:         &navQNH,
			NavAltitudeMCP: &navAlt,
			NavHeading:     &navHeading,
			TrueHeading:    &trueHeading,
			Lat:            &lat,
			Lon:            &lon,
			SeenPos:        &seenPos,
			Messages:       2,
			Seen:           1,
			RSSI:           -2,
			Category:       "A1",
			Squawk:         "7000",
			Emergency:      "none",
			SILType:        "perhour",
			Version:        2,
			NACP:           &nacp,
			NACV:           &nacv,
			NIC:            &nic,
			NICBaro:        &nicBaro,
			SIL:            &sil,
			GVA:            &gva,
			SDA:            &sda,
			RC:             &rc,
			ModeA:          true,
			ModeC:          true,
		}},
	}, nil
}

func TestCollectCoversMostBranches(t *testing.T) {
	// Receiver position for distance
	plat := 48.1
	plon := 11.6
	src := &stubRichSource{p: "stub"}
	CONFIG = *newDump1090MetricCollectorConfig(src)
	CONFIG.ReceiverPositionLat = &plat
	CONFIG.ReceiverPositionLon = &plon
	CONFIG.IsVerbose = true

	reg := prometheus.NewRegistry()
	m := newDump1090Metric(src)
	reg.MustRegister(m)

	// Expose via HTTP and query once to trigger Collect
	handler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	srv := httptest.NewServer(handler)
	defer srv.Close()

	resp, err := http.Get(srv.URL)
	if err != nil {
		t.Fatalf("http get: %v", err)
	}
	_ = resp.Body.Close()
}
