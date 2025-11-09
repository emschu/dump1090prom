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
	"testing"
)

func TestRollingMapSetGetEvict(t *testing.T) {
	rm := NewRollingFlightMap(2)

	// Initially empty
	if _, ok := rm.Get("a"); ok {
		t.Fatalf("expected key 'a' to be absent initially")
	}

	// Insert up to capacity
	rm.Set("a", "A1")
	rm.Set("b", "B1")

	if v, ok := rm.Get("a"); !ok || v != "A1" {
		t.Fatalf("expected to get A1 for key 'a', got %v, ok=%v", v, ok)
	}
	if v, ok := rm.Get("b"); !ok || v != "B1" {
		t.Fatalf("expected to get B1 for key 'b', got %v, ok=%v", v, ok)
	}

	// Insert beyond capacity -> evict oldest ('a')
	rm.Set("c", "C1")

	if _, ok := rm.Get("a"); ok {
		t.Fatalf("expected key 'a' to be evicted")
	}
	if v, ok := rm.Get("b"); !ok || v != "B1" {
		t.Fatalf("expected to get B1 for key 'b' after eviction, got %v, ok=%v", v, ok)
	}
	if v, ok := rm.Get("c"); !ok || v != "C1" {
		t.Fatalf("expected to get C1 for key 'c', got %v, ok=%v", v, ok)
	}
}

func TestRollingMapUpdateExisting(t *testing.T) {
	rm := NewRollingAirLineLabelMap(2)

	rm.Set("k", "v1")
	rm.Set("k", "v2") // update existing should not duplicate keys and should keep value

	if v, ok := rm.Get("k"); !ok || v != "v2" {
		t.Fatalf("expected updated value 'v2' for key 'k', got %v ok=%v", v, ok)
	}

	// Add another key and then another to force eviction once; capacity=2, keys order is [k,x], adding y evicts oldest 'k'
	rm.Set("x", "x1")
	rm.Set("y", "y1")

	if _, ok := rm.Get("k"); ok {
		t.Fatalf("expected key 'k' to be evicted as the oldest entry")
	}

	// Ensure the two most recent keys remain
	if v, ok := rm.Get("x"); !ok || v != "x1" {
		t.Fatalf("expected key 'x' to remain with value 'x1', got %v ok=%v", v, ok)
	}
	if v, ok := rm.Get("y"); !ok || v != "y1" {
		t.Fatalf("expected key 'y' to remain with value 'y1', got %v ok=%v", v, ok)
	}
}

func TestRollingMapDelete(t *testing.T) {
	rm := NewRollingAirLineMap(3)

	rm.Set("a", airline{ICAO: "AAA"})
	rm.Set("b", airline{ICAO: "BBB"})

	// Delete existing
	rm.Delete("a")
	if _, ok := rm.Get("a"); ok {
		t.Fatalf("expected key 'a' to be deleted")
	}

	// Delete non-existing should be no-op
	rm.Delete("missing")

	if v, ok := rm.Get("b"); !ok || v.ICAO != "BBB" {
		t.Fatalf("expected key 'b' to remain, got %v ok=%v", v, ok)
	}
}
