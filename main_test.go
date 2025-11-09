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
	"path/filepath"
	"testing"
)

func TestFindMetricSource_FileBasePath(t *testing.T) {
	tmp := t.TempDir()
	empty := ""
	src := findMetricSource(&empty, &tmp)
	if src == nil {
		t.Fatalf("expected non-nil source")
	}
	if src.getPath() != tmp {
		t.Fatalf("expected path %s, got %s", tmp, src.getPath())
	}
	if _, ok := src.(*FileMetricSource); !ok {
		t.Fatalf("expected *FileMetricSource, got %T", src)
	}
}

func TestFindMetricSource_URL(t *testing.T) {
	url := "http://example.com"
	empty := ""
	src := findMetricSource(&url, &empty)
	if src == nil {
		t.Fatalf("expected non-nil source")
	}
	if src.getPath() != url {
		t.Fatalf("expected url %s, got %s", url, src.getPath())
	}
	if _, ok := src.(*URLMetricSource); !ok {
		t.Fatalf("expected *URLMetricSource, got %T", src)
	}
}

func TestIsDir(t *testing.T) {
	tmp := t.TempDir()
	if !isDir(tmp) {
		t.Fatalf("expected %s to be a directory", tmp)
	}

	// Create a file
	f := filepath.Join(tmp, "file.txt")
	if err := os.WriteFile(f, []byte("x"), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if isDir(f) {
		t.Fatalf("expected %s to be a file, not directory", f)
	}

	// Non-existing path
	if isDir(filepath.Join(tmp, "doesnotexist")) {
		t.Fatalf("expected non-existing path to return false")
	}
}
