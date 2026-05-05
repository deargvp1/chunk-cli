package upgrade

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func writeFakeBrew(t *testing.T, upgradeExit int) string {
	t.Helper()
	dir := t.TempDir()

	if runtime.GOOS == "windows" {
		t.Skip("fake brew script not supported on Windows")
	}

	script := fmt.Sprintf(`#!/bin/sh
exit %d
`, upgradeExit)
	brewPath := filepath.Join(dir, "brew")
	if err := os.WriteFile(brewPath, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestRun(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		upgradeExit int
		wantErr     bool
		errContains string
	}{
		{
			name:        "brew not found",
			path:        "/nonexistent",
			wantErr:     true,
			errContains: "brew not found",
		},
		{
			name: "success",
		},
		{
			name:        "upgrade fails",
			upgradeExit: 1,
			wantErr:     true,
			errContains: "upgrade failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.path != "" {
				t.Setenv("PATH", tt.path)
			} else {
				dir := writeFakeBrew(t, tt.upgradeExit)
				t.Setenv("PATH", dir)
			}

			err := Run()
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Fatalf("expected error containing %q, got: %v", tt.errContains, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
