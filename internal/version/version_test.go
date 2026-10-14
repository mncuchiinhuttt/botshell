package version

import "testing"

func TestVersionConstants(t *testing.T) {
	if BaseVersion != "v1.0.0" {
		t.Errorf("expected v1.0.0, got %s", BaseVersion)
	}
	if CurrentVersion() != "v1.0.0" {
		t.Errorf("expected v1.0.0, got %s", CurrentVersion())
	}
	if Author != "@mncuchiinhuttt" {
		t.Errorf("expected @mncuchiinhuttt, got %s", Author)
	}
}
