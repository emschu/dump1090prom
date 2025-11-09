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
	"log"
	"strings"
	"sync"

	"github.com/golang/geo/s2"
	"github.com/prometheus/client_golang/prometheus"
)

const EarthRadiusKm = 6371.01
const Dump1090Namespace = "dump1090prom"

const (
	METRIC_AIRCRAFT_ADSB_VERSION                  = "aircraft_adsb_version"
	METRIC_AIRCRAFT_ALT_BARO_FEET                 = "aircraft_altitude_baro_feet"
	METRIC_AIRCRAFT_ALT_GEOM_FEET                 = "aircraft_altitude_geom_feet"
	METRIC_AIRCRAFT_BARO_RATE                     = "aircraft_barometric_vertical_rate_feet_per_minute"
	METRIC_AIRCRAFT_COUNT                         = "aircraft_count"
	METRIC_AIRCRAFT_DISTANCE_FROM_POSITION_METERS = "aircraft_distance_from_position_meters"
	METRIC_AIRCRAFT_FLIGHT_INFO                   = "aircraft_flight_info"
	METRIC_AIRCRAFT_GEOM_RATE                     = "aircraft_geometric_vertical_rate_feet_per_minute"
	METRIC_AIRCRAFT_GROUNDSPEED_KNOTS             = "aircraft_ground_speed_knots"
	METRIC_AIRCRAFT_GVA                           = "aircraft_gva"
	METRIC_AIRCRAFT_HEADING_DEGREE                = "aircraft_magnetic_heading_degrees"
	METRIC_AIRCRAFT_INDICATED_AIR_SPEED_KNOTS     = "aircraft_indicated_airspeed_knots"
	METRIC_AIRCRAFT_LAT                           = "aircraft_latitude"
	METRIC_AIRCRAFT_LON                           = "aircraft_longitude"
	METRIC_AIRCRAFT_MACH_NUMBER                   = "aircraft_mach_number"
	METRIC_AIRCRAFT_MODE_A                        = "aircraft_mode_a"
	METRIC_AIRCRAFT_MODE_C                        = "aircraft_mode_c"
	METRIC_AIRCRAFT_NAC_P                         = "aircraft_nac_p"
	METRIC_AIRCRAFT_NAC_V                         = "aircraft_nac_v"
	METRIC_AIRCRAFT_NAV_ALTITUDE_MCP_FEET         = "aircraft_nav_altitude_mcp_feet"
	METRIC_AIRCRAFT_NAV_HEADING_DEGREE            = "aircraft_nav_heading_degrees"
	METRIC_AIRCRAFT_NAV_QNH_MiLLIBAR              = "aircraft_nav_qnh_millibar"
	METRIC_AIRCRAFT_NIC                           = "aircraft_nic"
	METRIC_AIRCRAFT_NIC_BARO                      = "aircraft_nic_baro"
	METRIC_AIRCRAFT_OAT                           = "aircraft_oat"
	METRIC_AIRCRAFT_RC                            = "aircraft_rc"
	METRIC_AIRCRAFT_ROLL_DEGREE                   = "aircraft_roll_degrees"
	METRIC_AIRCRAFT_RSSI_DBM                      = "aircraft_rssi_dbm"
	METRIC_AIRCRAFT_SYSTEM_DESIGN_ASSURANCE       = "aircraft_sda"
	METRIC_AIRCRAFT_SEEN_POS_SECOND               = "aircraft_seen_pos_seconds"
	METRIC_AIRCRAFT_SEEN_SECOND                   = "aircraft_seen_seconds"
	METRIC_AIRCRAFT_SIL                           = "aircraft_sil"
	METRIC_AIRCRAFT_SPI                           = "aircraft_spi"
	METRIC_AIRCRAFT_TAT                           = "aircraft_tat"
	METRIC_AIRCRAFT_TRACK_DEGREES                 = "aircraft_track_degrees"
	METRIC_AIRCRAFT_TRACK_RATE                    = "aircraft_track_rate_degrees_per_second"
	METRIC_AIRCRAFT_TRUE_AIR_SPEED_KNOTS          = "aircraft_true_airspeed_knots"
	METRIC_AIRCRAFT_TRUE_HEADING_DEGREE           = "aircraft_true_heading_degrees"
	METRIC_NOW_TIMESTAMP                          = "now_timestamp"
	METRIC_TOTAL_MESSAGES                         = "total_messages"
)

var metrics = []dump1090MetricItem{
	{
		Name:      METRIC_NOW_TIMESTAMP,
		Namespace: Dump1090Namespace,
		Help:      "Current timestamp in seconds since the epoch",
		Type:      METRIC_ITEM_TYPE_GAUGE,
	}, {
		Name:      METRIC_TOTAL_MESSAGES,
		Namespace: Dump1090Namespace,
		Help:      "Total number of messages received",
		Type:      METRIC_ITEM_TYPE_COUNTER,
	}, {
		Name:      METRIC_AIRCRAFT_COUNT,
		Namespace: Dump1090Namespace,
		Help:      "Number of different aircraft currently seen",
		Type:      METRIC_ITEM_TYPE_GAUGE,
	}, {
		Name:      METRIC_AIRCRAFT_ALT_BARO_FEET,
		Namespace: Dump1090Namespace,
		Help:      "Barometric altitude in feet",
		Labels:    []string{"hex", "flight"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	}, {
		Name:      METRIC_AIRCRAFT_ALT_GEOM_FEET,
		Namespace: Dump1090Namespace,
		Help:      "Geometric altitude in feet",
		Labels:    []string{"hex", "flight"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	}, {
		Name:      METRIC_AIRCRAFT_BARO_RATE,
		Namespace: Dump1090Namespace,
		Help:      "Barometric vertical rate in feet per minute",
		Labels:    []string{"hex", "flight"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	}, {
		Name:      METRIC_AIRCRAFT_GEOM_RATE,
		Namespace: Dump1090Namespace,
		Help:      "Geometric vertical rate in feet per minute",
		Labels:    []string{"hex", "flight"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	}, {
		Name:      METRIC_AIRCRAFT_SEEN_SECOND,
		Namespace: Dump1090Namespace,
		Help:      "Seconds since this aircraft was last seen",
		Labels:    []string{"hex", "flight"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	}, {
		Name:      METRIC_AIRCRAFT_RSSI_DBM,
		Namespace: Dump1090Namespace,
		Help:      "Received Signal Strength Indicator in dBm",
		Labels:    []string{"hex", "flight"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	}, {
		Name:      METRIC_AIRCRAFT_GROUNDSPEED_KNOTS,
		Namespace: Dump1090Namespace,
		Help:      "Ground speed in knots",
		Labels:    []string{"hex", "flight"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	}, {
		Name:      METRIC_AIRCRAFT_INDICATED_AIR_SPEED_KNOTS,
		Namespace: Dump1090Namespace,
		Help:      "Indicated airspeed in knots",
		Labels:    []string{"hex", "flight"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	}, {
		Name:      METRIC_AIRCRAFT_FLIGHT_INFO,
		Namespace: Dump1090Namespace,
		Help:      "Metadata about the flight and aircraft",
		Labels:    []string{"hex", "flight", "lat", "lon", "category", "squawk", "direction", "distance", "airline"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	}, {
		Name:      METRIC_AIRCRAFT_TRUE_AIR_SPEED_KNOTS,
		Namespace: Dump1090Namespace,
		Help:      "True airspeed in knots",
		Labels:    []string{"hex", "flight"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	}, {
		Name:      METRIC_AIRCRAFT_MACH_NUMBER,
		Namespace: Dump1090Namespace,
		Help:      "Mach number",
		Labels:    []string{"hex", "flight"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	}, {
		Name:      METRIC_AIRCRAFT_TRACK_DEGREES,
		Namespace: Dump1090Namespace,
		Help:      "Track angle in degrees (0-359)",
		Labels:    []string{"hex", "flight", "direction"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	}, {
		Name:      METRIC_AIRCRAFT_TRACK_RATE,
		Namespace: Dump1090Namespace,
		Help:      "Rate of change of track angle in degrees per second",
		Labels:    []string{"hex", "flight", "direction"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	}, {
		Name:      METRIC_AIRCRAFT_ROLL_DEGREE,
		Namespace: Dump1090Namespace,
		Help:      "Roll angle in degrees",
		Labels:    []string{"hex", "flight"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	}, {
		Name:      METRIC_AIRCRAFT_HEADING_DEGREE,
		Namespace: Dump1090Namespace,
		Help:      "Magnetic heading in degrees",
		Labels:    []string{"hex", "flight", "direction"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	}, {
		Name:      METRIC_AIRCRAFT_NAV_QNH_MiLLIBAR,
		Namespace: Dump1090Namespace,
		Help:      "QNH setting in millibars",
		Labels:    []string{"hex", "flight"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	}, {
		Name:      METRIC_AIRCRAFT_NAV_ALTITUDE_MCP_FEET,
		Namespace: Dump1090Namespace,
		Help:      "MCP/FCU selected altitude in feet",
		Labels:    []string{"hex", "flight"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	}, {
		Name:      METRIC_AIRCRAFT_NAV_HEADING_DEGREE,
		Namespace: Dump1090Namespace,
		Help:      "Selected heading in degrees",
		Labels:    []string{"hex", "flight", "direction"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	}, {
		Name:      METRIC_AIRCRAFT_TRUE_HEADING_DEGREE,
		Namespace: Dump1090Namespace,
		Help:      "Selected true heading in degrees",
		Labels:    []string{"hex", "flight", "direction"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	}, {
		Name:      METRIC_AIRCRAFT_LAT,
		Namespace: Dump1090Namespace,
		Help:      "Latitude of aircraft",
		Labels:    []string{"hex", "flight"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	}, {
		Name:      METRIC_AIRCRAFT_LON,
		Namespace: Dump1090Namespace,
		Help:      "Longitude of aircraft",
		Labels:    []string{"hex", "flight"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	}, {
		Name:      METRIC_AIRCRAFT_SEEN_POS_SECOND,
		Namespace: Dump1090Namespace,
		Help:      "Seconds since position was last updated",
		Labels:    []string{"hex", "flight"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	}, {
		Name:      METRIC_AIRCRAFT_ADSB_VERSION,
		Namespace: Dump1090Namespace,
		Help:      "Version of the ADS-B protocol in use",
		Labels:    []string{"hex", "flight"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	}, {
		Name:      METRIC_AIRCRAFT_NIC,
		Namespace: Dump1090Namespace,
		Help:      "Navigation Integrity Category",
		Labels:    []string{"hex", "flight"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	}, {
		Name:      METRIC_AIRCRAFT_RC,
		Namespace: Dump1090Namespace,
		Help:      "Radius of containment",
		Labels:    []string{"hex", "flight"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	}, {
		Name:      METRIC_AIRCRAFT_NIC_BARO,
		Namespace: Dump1090Namespace,
		Help:      "Navigation Integrity Category for barometric altitude",
		Labels:    []string{"hex", "flight"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	}, {
		Name:      METRIC_AIRCRAFT_NAC_P,
		Namespace: Dump1090Namespace,
		Help:      "Navigation Accuracy Category for Position",
		Labels:    []string{"hex", "flight"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	}, {
		Name:      METRIC_AIRCRAFT_NAC_V,
		Namespace: Dump1090Namespace,
		Help:      "Navigation Accuracy Category for Velocity",
		Labels:    []string{"hex", "flight"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	}, {
		Name:      METRIC_AIRCRAFT_SIL,
		Namespace: Dump1090Namespace,
		Help:      "Surveillance Integrity Level",
		Labels:    []string{"hex", "flight"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	}, {
		Name:      METRIC_AIRCRAFT_GVA,
		Namespace: Dump1090Namespace,
		Help:      "Geometric Vertical Accuracy",
		Labels:    []string{"hex", "flight"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	}, {
		Name:      METRIC_AIRCRAFT_SYSTEM_DESIGN_ASSURANCE,
		Namespace: Dump1090Namespace,
		Help:      "System Design Assurance",
		Labels:    []string{"hex", "flight"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	}, {
		Name:      METRIC_AIRCRAFT_MODE_A,
		Namespace: Dump1090Namespace,
		Help:      "Mode A (ident) capability (1 if present)",
		Labels:    []string{"hex", "flight"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	}, {
		Name:      METRIC_AIRCRAFT_MODE_C,
		Namespace: Dump1090Namespace,
		Help:      "Mode C (altitude) capability (1 if present)",
		Labels:    []string{"hex", "flight"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	}, {
		Name:      METRIC_AIRCRAFT_DISTANCE_FROM_POSITION_METERS,
		Namespace: Dump1090Namespace,
		Help:      "Distance in meters from the recording position",
		Labels:    []string{"hex", "flight", "direction"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	}, {
		Name:      METRIC_AIRCRAFT_SPI,
		Namespace: Dump1090Namespace,
		Help:      "Special Position Identification (IDENT)",
		Labels:    []string{"hex", "flight"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	}, {
		Name:      METRIC_AIRCRAFT_OAT,
		Namespace: Dump1090Namespace,
		Help:      "Outside Air Temperature in Celsius",
		Labels:    []string{"hex", "flight"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	}, {
		Name:      METRIC_AIRCRAFT_TAT,
		Namespace: Dump1090Namespace,
		Help:      "Total Air Temperature in Celsius",
		Labels:    []string{"hex", "flight"},
		Type:      METRIC_ITEM_TYPE_GAUGE_VEC,
	},
}

// supported metric types
const (
	METRIC_ITEM_TYPE_GAUGE = iota
	METRIC_ITEM_TYPE_GAUGE_VEC
	METRIC_ITEM_TYPE_COUNTER
)

type metricItemType int

type dump1090MetricItem struct {
	Name      string
	Namespace string
	Help      string
	Labels    []string
	Type      metricItemType

	Gauge    prometheus.Gauge
	GaugeVec *prometheus.GaugeVec
	Counter  prometheus.Counter
}

type dump1090Metric struct {
	aircraftJsonFilePath string
	itemMap              sync.Map
	// TODO: missing: nav_altitude_fms, wind_speed, mrar_source, temperature, pressure, turbulence, humidity
}

// used to keep the flight info cached once found, and limits entries
var flightMap = NewRollingFlightMap(CONFIG.RollingMapDefaultSize)
var flightAirLineMap = NewRollingAirLineMap(CONFIG.RollingMapDefaultSize)
var flightAirLineLabelMap = NewRollingAirLineLabelMap(CONFIG.RollingMapDefaultSize)

func newDump1090Metric(source MetricSource) *dump1090Metric {
	var metric = &dump1090Metric{
		aircraftJsonFilePath: source.getPath(),
	}
	// Ensure all metric definitions are loaded so collectors can find them
	metric.LoadAll()
	return metric
}

func (d *dump1090Metric) Describe(ch chan<- *prometheus.Desc) {
	d.DescribeAll(ch)
}

func (d *dump1090Metric) Collect(ch chan<- prometheus.Metric) {
	data, err := CONFIG.Source.fetchMetrics()
	if err != nil {
		log.Printf("Error fetching metrics: %s", err)
		return
	}
	d.mustFindMetric(METRIC_NOW_TIMESTAMP).Gauge.Set(data.Now)
	d.mustFindMetric(METRIC_TOTAL_MESSAGES).Counter.Add(float64(data.Messages))
	d.mustFindMetric(METRIC_AIRCRAFT_COUNT).Gauge.Set(float64(len(data.Aircraft)))

	if CONFIG.IsVerbose {
		log.Printf("Collection request – Total aircraft: %d\n", len(data.Aircraft))
	}

	// ResetAll all per-aircraft metrics to remove stale data
	d.ResetAll()

	// Update per-aircraft metrics
	for _, ac := range data.Aircraft {
		hex := strings.TrimSpace(ac.Hex)
		flight := strings.TrimSpace(ac.Flight)

		var distance float64
		if ac.Lat != nil && ac.Lon != nil {
			distance = d.calculateDistanceToAircraft(&ac)
		}

		if hex != "" && flight != "" {
			flightMap.Set(hex, flight)
		}
		var exist bool
		flight, exist = flightMap.Get(hex)
		if !exist {
			flight = ""
		}

		var aircraftLat string
		if ac.Lat == nil {
			aircraftLat = ""
		} else {
			aircraftLat = fmt.Sprintf("%.12f", *ac.Lat)
		}
		var aircraftLon string
		if ac.Lon == nil {
			aircraftLon = ""
		} else {
			aircraftLon = fmt.Sprintf("%.12f", *ac.Lon)
		}
		var direction string
		if ac.Track != nil {
			direction = getDirection(*ac.Track)
		}
		var airline string
		if CONFIG.IsAirlineLabellingEnabled {
			possibleAirline := getAirlineLabel(flight)
			if possibleAirline != nil {
				airline = strings.TrimSpace(*possibleAirline)
			}
		}

		// Basic metrics
		d.mustFindMetric(METRIC_AIRCRAFT_SEEN_SECOND).GaugeVec.WithLabelValues(hex, flight).Set(ac.Seen)
		d.mustFindMetric(METRIC_AIRCRAFT_RSSI_DBM).GaugeVec.WithLabelValues(hex, flight).Set(ac.RSSI)

		// Position and movement metrics
		if ac.AltBaro != nil {
			d.mustFindMetric(METRIC_AIRCRAFT_ALT_BARO_FEET).GaugeVec.WithLabelValues(hex, flight).Set(*ac.AltBaro)
		}
		if ac.AltGeom != nil {
			d.mustFindMetric(METRIC_AIRCRAFT_ALT_GEOM_FEET).GaugeVec.WithLabelValues(hex, flight).Set(*ac.AltGeom)
		}
		if ac.GS != nil {
			d.mustFindMetric(METRIC_AIRCRAFT_GROUNDSPEED_KNOTS).GaugeVec.WithLabelValues(hex, flight).Set(*ac.GS)
		}
		if ac.IAS != nil {
			d.mustFindMetric(METRIC_AIRCRAFT_INDICATED_AIR_SPEED_KNOTS).GaugeVec.WithLabelValues(hex, flight).Set(*ac.IAS)
		}
		if ac.TAS != nil {
			d.mustFindMetric(METRIC_AIRCRAFT_TRUE_AIR_SPEED_KNOTS).GaugeVec.WithLabelValues(hex, flight).Set(*ac.TAS)
		}
		if ac.Mach != nil {
			d.mustFindMetric(METRIC_AIRCRAFT_MACH_NUMBER).GaugeVec.WithLabelValues(hex, flight).Set(*ac.Mach)
		}
		if ac.Track != nil {
			d.mustFindMetric(METRIC_AIRCRAFT_TRACK_DEGREES).GaugeVec.WithLabelValues(hex, flight, direction).Set(*ac.Track)
		}
		if ac.TrackRate != nil {
			d.mustFindMetric(METRIC_AIRCRAFT_TRACK_RATE).GaugeVec.WithLabelValues(hex, flight, direction).Set(*ac.TrackRate)
		}
		if ac.Roll != nil {
			d.mustFindMetric(METRIC_AIRCRAFT_ROLL_DEGREE).GaugeVec.WithLabelValues(hex, flight).Set(*ac.Roll)
		}
		if ac.MagHeading != nil {
			d.mustFindMetric(METRIC_AIRCRAFT_HEADING_DEGREE).GaugeVec.WithLabelValues(hex, flight, direction).Set(*ac.MagHeading)
		}
		if ac.BaroRate != nil {
			d.mustFindMetric(METRIC_AIRCRAFT_BARO_RATE).GaugeVec.WithLabelValues(hex, flight).Set(*ac.BaroRate)
		}
		if ac.GeomRate != nil {
			d.mustFindMetric(METRIC_AIRCRAFT_GEOM_RATE).GaugeVec.WithLabelValues(hex, flight).Set(*ac.GeomRate)
		}
		if ac.NavQNH != nil {
			d.mustFindMetric(METRIC_AIRCRAFT_NAV_QNH_MiLLIBAR).GaugeVec.WithLabelValues(hex, flight).Set(*ac.NavQNH)
		}
		if ac.NavAltitudeMCP != nil {
			d.mustFindMetric(METRIC_AIRCRAFT_NAV_ALTITUDE_MCP_FEET).GaugeVec.WithLabelValues(hex, flight).Set(*ac.NavAltitudeMCP)
		}
		if ac.NavHeading != nil {
			d.mustFindMetric(METRIC_AIRCRAFT_NAV_HEADING_DEGREE).GaugeVec.WithLabelValues(hex, flight, direction).Set(*ac.NavHeading)
		}
		if ac.TrueHeading != nil {
			d.mustFindMetric(METRIC_AIRCRAFT_TRUE_HEADING_DEGREE).GaugeVec.WithLabelValues(hex, flight, direction).Set(*ac.TrueHeading)
		}
		if ac.Lat != nil {
			d.mustFindMetric(METRIC_AIRCRAFT_LAT).GaugeVec.WithLabelValues(hex, flight).Set(*ac.Lat)
		}
		if ac.Lon != nil {
			d.mustFindMetric(METRIC_AIRCRAFT_LON).GaugeVec.WithLabelValues(hex, flight).Set(*ac.Lon)
		}
		if ac.SeenPos != nil {
			d.mustFindMetric(METRIC_AIRCRAFT_SEEN_POS_SECOND).GaugeVec.WithLabelValues(hex, flight).Set(*ac.SeenPos)
		}

		d.mustFindMetric(METRIC_AIRCRAFT_FLIGHT_INFO).GaugeVec.WithLabelValues(hex, flight, aircraftLat, aircraftLon, ac.Category, ac.Squawk, direction, fmt.Sprintf("%.6f", distance), airline).Set(1)

		d.mustFindMetric(METRIC_AIRCRAFT_ADSB_VERSION).GaugeVec.WithLabelValues(hex, flight).Set(float64(ac.Version))

		if ac.NIC != nil {
			d.mustFindMetric(METRIC_AIRCRAFT_NIC).GaugeVec.WithLabelValues(hex, flight).Set(float64(*ac.NIC))
		}
		if ac.RC != nil {
			d.mustFindMetric(METRIC_AIRCRAFT_RC).GaugeVec.WithLabelValues(hex, flight).Set(float64(*ac.RC))
		}
		if ac.NICBaro != nil {
			d.mustFindMetric(METRIC_AIRCRAFT_NIC_BARO).GaugeVec.WithLabelValues(hex, flight).Set(float64(*ac.NICBaro))
		}
		if ac.NACP != nil {
			d.mustFindMetric(METRIC_AIRCRAFT_NAC_P).GaugeVec.WithLabelValues(hex, flight).Set(float64(*ac.NACP))
		}
		if ac.NACV != nil {
			d.mustFindMetric(METRIC_AIRCRAFT_NAC_V).GaugeVec.WithLabelValues(hex, flight).Set(float64(*ac.NACV))
		}
		if ac.SIL != nil {
			d.mustFindMetric(METRIC_AIRCRAFT_SIL).GaugeVec.WithLabelValues(hex, flight).Set(float64(*ac.SIL))
		}
		if ac.GVA != nil {
			d.mustFindMetric(METRIC_AIRCRAFT_GVA).GaugeVec.WithLabelValues(hex, flight).Set(float64(*ac.GVA))
		}
		if ac.SDA != nil {
			d.mustFindMetric(METRIC_AIRCRAFT_SYSTEM_DESIGN_ASSURANCE).GaugeVec.WithLabelValues(hex, flight).Set(float64(*ac.SDA))
		}

		// Mode A/C flags
		if ac.ModeA {
			d.mustFindMetric(METRIC_AIRCRAFT_MODE_A).GaugeVec.WithLabelValues(hex, flight).Set(1)
		} else {
			d.mustFindMetric(METRIC_AIRCRAFT_MODE_A).GaugeVec.WithLabelValues(hex, flight).Set(0)
		}
		if ac.ModeC {
			d.mustFindMetric(METRIC_AIRCRAFT_MODE_C).GaugeVec.WithLabelValues(hex, flight).Set(1)
		} else {
			d.mustFindMetric(METRIC_AIRCRAFT_MODE_C).GaugeVec.WithLabelValues(hex, flight).Set(0)
		}

		if ac.Lat != nil && ac.Lon != nil {
			d.mustFindMetric(METRIC_AIRCRAFT_DISTANCE_FROM_POSITION_METERS).GaugeVec.WithLabelValues(hex, flight, direction).Set(distance)
		}

		if ac.SPI != nil {
			d.mustFindMetric(METRIC_AIRCRAFT_SPI).GaugeVec.WithLabelValues(hex, flight).Set(float64(*ac.SPI))
		}
		if ac.OAT != nil {
			d.mustFindMetric(METRIC_AIRCRAFT_OAT).GaugeVec.WithLabelValues(hex, flight).Set(float64(*ac.OAT))
		}
		if ac.TAT != nil {
			d.mustFindMetric(METRIC_AIRCRAFT_TAT).GaugeVec.WithLabelValues(hex, flight).Set(float64(*ac.TAT))
		}
	}

	d.CollectAll(ch)
}

func (d *dump1090Metric) calculateDistanceToAircraft(ac *AircraftStateRecord) float64 {
	if !CONFIG.IsDistanceCalculationEnabled {
		return 0.0
	}
	if (CONFIG.ReceiverPositionLat == nil || CONFIG.ReceiverPositionLon == nil) ||
		(*CONFIG.ReceiverPositionLat == 0.0 && *CONFIG.ReceiverPositionLon == 0.0) {
		return 0.0
	}
	var planeLat = ac.Lat
	var planeLon = ac.Lon
	var positionLat = *CONFIG.ReceiverPositionLat
	var positionLon = *CONFIG.ReceiverPositionLon
	p1 := s2.PointFromLatLng(s2.LatLngFromDegrees(*planeLat, *planeLon))
	p2 := s2.PointFromLatLng(s2.LatLngFromDegrees(positionLat, positionLon))
	distance := p1.Distance(p2).Radians() * EarthRadiusKm * 1000 // in meters
	return distance
}

func (d *dump1090Metric) mustFindMetric(metricName string) *dump1090MetricItem {
	item, exists := d.itemMap.Load(metricName)
	if !exists {
		// this is fatal!
		log.Fatalf("Metric '%s' not found", metricName)
	}
	return item.(*dump1090MetricItem)
}

func (d *dump1090Metric) LoadAll() {
	d.itemMap.Clear()
	for _, i := range metrics {
		var vec *prometheus.GaugeVec
		var gauge prometheus.Gauge
		var counter prometheus.Counter

		switch i.Type {
		case METRIC_ITEM_TYPE_GAUGE:
			gauge = prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: i.Namespace,
				Name:      i.Name,
				Help:      i.Help,
			})
		case METRIC_ITEM_TYPE_GAUGE_VEC:
			vec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: i.Namespace,
				Name:      i.Name,
				Help:      i.Help,
			}, i.Labels)
		case METRIC_ITEM_TYPE_COUNTER:
			counter = prometheus.NewCounter(prometheus.CounterOpts{
				Namespace: i.Namespace,
				Name:      i.Name,
				Help:      i.Help,
			})
		default:
			log.Fatalf("Unknown metric type: %d", i.Type)
		}

		item := dump1090MetricItem{
			Name:      i.Name,
			Namespace: i.Namespace,
			Help:      i.Help,
			Labels:    i.Labels,
			GaugeVec:  vec,
			Gauge:     gauge,
			Counter:   counter,
		}
		d.itemMap.Store(i.Name, &item)
	}
}

func (d *dump1090Metric) ResetAll() {
	d.itemMap.Range(func(key, value interface{}) bool {
		item := value.(*dump1090MetricItem)
		if item.GaugeVec != nil {
			item.GaugeVec.Reset()
		}
		return true
	})
}

func (d *dump1090Metric) CollectAll(ch chan<- prometheus.Metric) {
	d.itemMap.Range(func(key, value interface{}) bool {
		item := value.(*dump1090MetricItem)
		if item.GaugeVec != nil {
			item.GaugeVec.Collect(ch)
		}
		if item.Gauge != nil {
			ch <- item.Gauge
		}
		if item.Counter != nil {
			ch <- item.Counter
		}
		return true
	})
}

func (d *dump1090Metric) DescribeAll(ch chan<- *prometheus.Desc) {
	d.itemMap.Range(func(key, value interface{}) bool {
		item := value.(*dump1090MetricItem)
		if item.GaugeVec != nil {
			ch <- item.GaugeVec.WithLabelValues(item.Labels...).Desc()
		}
		if item.Gauge != nil {
			ch <- item.Gauge.Desc()
		}
		if item.Counter != nil {
			ch <- item.Counter.Desc()
		}
		return true
	})
}
