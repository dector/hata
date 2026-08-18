package version

import "testing"

func TestGetDefaultVersion(t *testing.T) {
	t.Setenv("HATA_VERSION", "")
	oldVersion := Version
	Version = ""
	t.Cleanup(func() { Version = oldVersion })

	if got := Get(); got != "dev" {
		t.Fatalf("Get() = %q, want %q", got, "dev")
	}
}

func TestGetBuildTimeVersion(t *testing.T) {
	t.Setenv("HATA_VERSION", "")
	oldVersion := Version
	Version = "v1.2.3-4-gabcdef0-dirty"
	t.Cleanup(func() { Version = oldVersion })

	if got := Get(); got != "v1.2.3-4-gabcdef0-dirty" {
		t.Fatalf("Get() = %q", got)
	}
}

func TestGetRuntimeVersionOverride(t *testing.T) {
	t.Setenv("HATA_VERSION", "runtime-version")
	oldVersion := Version
	Version = "build-version"
	t.Cleanup(func() { Version = oldVersion })

	if got := Get(); got != "runtime-version" {
		t.Fatalf("Get() = %q, want %q", got, "runtime-version")
	}
}
