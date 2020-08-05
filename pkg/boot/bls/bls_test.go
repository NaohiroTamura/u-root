// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bls

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/boot/boottest"
	"github.com/u-root/u-root/pkg/ulog/ulogtest"
)

var blsEntries = []struct {
	fsRoot string
	entry string
	err   string
}{
	{
		fsRoot: "./testdata/madeup",
		entry: "entry-1.conf",
	},
	{
		fsRoot: "./testdata/madeup",
		entry: "entry-2.conf",
		err:   "neither linux, efi, nor multiboot present in BootLoaderSpec config",
	},
	{
		fsRoot: "./testdata/centos_8",
		entry: "9151c934991a474b9f86e3e7be1d25e6-4.18.0-147.el8.x86_64.conf",
	},
	{
		fsRoot: "./testdata/centos_8",
		entry: "9151c934991a474b9f86e3e7be1d25e6-0-rescue.conf",
	},
}

func TestParseBLSEntries(t *testing.T) {

	for _, tt := range blsEntries {
		dir := filepath.Join(tt.fsRoot, "loader/entries")

		t.Run(tt.entry, func(t *testing.T) {
			image, err := parseBLSEntry(filepath.Join(dir, tt.entry), tt.fsRoot)
			if err != nil {
				if tt.err == "" {
					t.Fatalf("Got error %v", err)
				}
				if !strings.Contains(err.Error(), tt.err) {
					t.Fatalf("Got error %v, expected error to contain %s", err, tt.err)
				}
				return
			}
			if tt.err != "" {
				t.Fatalf("Expected error %s, got no error", tt.err)
			}
			t.Logf("Got image: %s", image.String())
		})
	}
}

func TestScanBLSEntries(t *testing.T) {
	// find all saved configs
	tests, err := filepath.Glob("testdata/*.json")
	if err != nil {
		t.Error("Failed to find test config files:", err)
	}

	for _, test := range tests {
		configPath := strings.TrimRight(test, ".json")
		t.Run(configPath, func(t *testing.T) {
			want, err := ioutil.ReadFile(test)
			if err != nil {
				t.Errorf("Failed to read test json '%v':%v", test, err)
			}

			imgs, err := ScanBLSEntries(ulogtest.Logger{t}, configPath)
			if err != nil {
				t.Fatalf("Failed to parse %s: %v", test, err)
			}

			if err := boottest.CompareImagesToJSON(imgs, want); err != nil {
				t.Errorf("ParseLocalConfig(): %v", err)
			}
		})
	}
}

// Enable this temporarily to generate new configs. Double-check them by hand.
func DISABLEDTestGenerateConfigs(t *testing.T) {
	tests, err := filepath.Glob("testdata/*.json")
	if err != nil {
		t.Error("Failed to find test config files:", err)
	}

	for _, test := range tests {
		configPath := strings.TrimRight(test, ".json")
		t.Run(configPath, func(t *testing.T) {
			imgs, err := ScanBLSEntries(ulogtest.Logger{t}, configPath)
			if err != nil {
				t.Fatalf("Failed to parse %s: %v", test, err)
			}

			if err := boottest.ToJSONFile(imgs, test); err != nil {
				t.Errorf("failed to generate file: %v", err)
			}
		})
	}
}
