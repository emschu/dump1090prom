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

// dump1090MetricCollectorConfig config struct
type dump1090MetricCollectorConfig struct {
	Source                         MetricSource
	IsDistanceCalculationEnabled   bool
	IsAirlineLabellingEnabled      bool
	IsExposingOriginalFilesEnabled bool
	Host                           string
	Port                           int
	ReceiverPositionLat            *float64
	ReceiverPositionLon            *float64
	IsVerbose                      bool
	RollingMapDefaultSize          int
}

func newDump1090MetricCollectorConfig(source MetricSource) *dump1090MetricCollectorConfig {
	return &dump1090MetricCollectorConfig{
		Source:                         source,
		IsDistanceCalculationEnabled:   true,
		IsAirlineLabellingEnabled:      true,
		IsExposingOriginalFilesEnabled: true,
		Host:                           "127.0.0.1",
		Port:                           8080,
		ReceiverPositionLat:            nil,
		ReceiverPositionLon:            nil,
		IsVerbose:                      false,
		RollingMapDefaultSize:          1000,
	}
}
