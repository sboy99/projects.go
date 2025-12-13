package version

import (
	"strings"
	"testing"
)

func TestGetVersion(t *testing.T) {
	version := GetVersion()
	// Version can be empty if not set via ldflags, which is acceptable
	_ = version
}

func TestGetBuildInfo(t *testing.T) {
	info := GetBuildInfo()
	if info == "" {
		t.Error("GetBuildInfo() returned empty string")
	}

	// Should contain version, commit, and build time labels
	if !strings.Contains(info, "Version:") {
		t.Error("GetBuildInfo() should contain 'Version:'")
	}
	if !strings.Contains(info, "Commit:") {
		t.Error("GetBuildInfo() should contain 'Commit:'")
	}
	if !strings.Contains(info, "Build Time:") {
		t.Error("GetBuildInfo() should contain 'Build Time:'")
	}
}

func TestString(t *testing.T) {
	str := String()
	if str == "" {
		t.Error("String() returned empty string")
	}

	// Should be same as GetBuildInfo
	info := GetBuildInfo()
	if str != info {
		t.Errorf("String() = %v, want %v", str, info)
	}
}
