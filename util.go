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
	"os"
	"sync"
)

type RollingFlightMap = GenericRollingStringKeyValueStore[string]
type RollingAirlineMap = GenericRollingStringKeyValueStore[airline]
type RollingAirlineLabelMap = GenericRollingStringKeyValueStore[string]

type GenericRollingStringKeyValueStore[T any] struct {
	capacity int
	keys     []string
	mu       sync.Mutex
	store    sync.Map // string -> T
}

func NewRollingFlightMap(n int) *RollingFlightMap {
	return &RollingFlightMap{
		capacity: n,
		keys:     make([]string, 0, n),
	}
}

func NewRollingAirLineLabelMap(i int) *RollingAirlineLabelMap {
	return &RollingAirlineLabelMap{
		capacity: i,
		keys:     make([]string, 0, i),
	}
}

func NewRollingAirLineMap(i int) *RollingAirlineMap {
	return &RollingAirlineMap{
		capacity: i,
		keys:     make([]string, 0, i),
	}
}

func (rm *GenericRollingStringKeyValueStore[T]) Set(key string, value T) {
	// Fast path: update existing value
	if _, exists := rm.store.Load(key); exists {
		rm.store.Store(key, value)
		return
	}
	// New key path with eviction and keys management
	rm.mu.Lock()
	defer rm.mu.Unlock()
	// Double-check under lock to avoid duplicate keys slice entries
	if _, exists := rm.store.Load(key); exists {
		rm.store.Store(key, value)
		return
	}
	if len(rm.keys) >= rm.capacity {
		if len(rm.keys) > 0 {
			oldest := rm.keys[0]
			rm.keys = rm.keys[1:]
			rm.store.Delete(oldest)
		}
	}
	rm.keys = append(rm.keys, key)
	rm.store.Store(key, value)
}

func (rm *GenericRollingStringKeyValueStore[T]) Get(key string) (T, bool) {
	var zero T
	if val, ok := rm.store.Load(key); ok {
		return val.(T), true
	}
	return zero, false
}

func (rm *GenericRollingStringKeyValueStore[T]) Delete(key string) {
	// Check existence first
	if _, exists := rm.store.Load(key); !exists {
		return
	}
	// Remove from map and keys slice
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.store.Delete(key)
	for i, k := range rm.keys {
		if k == key {
			rm.keys = append(rm.keys[:i], rm.keys[i+1:]...)
			break
		}
	}
}

var windDirections = []string{
	"N",   // 0°
	"NNE", // 22.5°
	"NE",  // 45°
	"ENE", // 67.5°
	"E",   // 90°
	"ESE", // 112.5°
	"SE",  // 135°
	"SSE", // 157.5°
	"S",   // 180°
	"SSW", // 202.5°
	"SW",  // 225°
	"WSW", // 247.5°
	"W",   // 270°
	"WNW", // 292.5°
	"NW",  // 315°
	"NNW", // 337.5°
}

func getDirection(track float64) string {
	var degree = int((track / float64(360.0/len(windDirections))) + .5)
	return windDirections[(degree % len(windDirections))]
}

func isDir(file string) bool {
	fileStat, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}
	return fileStat.IsDir()
}
