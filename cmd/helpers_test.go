package cmd

import (
	"os"
	"strings"
	"testing"
	"time"
)

// Test parsing of the input times using days.
func Test_parseTime_Days(t *testing.T) {
	if secs, err := parseTime("1d"); secs != 86400 && err == nil {
		t.Fatalf("expected 86400, got %d", secs)
		t.Fatal(err)
	}

	if secs, err := parseTime("365d"); secs != 31536000 && err == nil {
		t.Fatalf("expected 31536000, got %d", secs)
		t.Fatal(err)
	}
}

// Test parsing of the input times using hours.
func Test_parseTime_Hours(t *testing.T) {
	if secs, err := parseTime("24h"); secs != 86400 && err == nil {
		t.Fatalf("expected 86400, got %d", secs)
		t.Fatal(err)
	}

	if secs, err := parseTime("1h"); secs != 3600 && err == nil {
		t.Fatalf("expected 3600, got %d", secs)
		t.Fatal(err)
	}
}

// Test parsing of the input times using minutes.
func Test_parseTime_Minutes(t *testing.T) {
	if secs, err := parseTime("60m"); secs != 3600 && err == nil {
		t.Fatalf("expected 3600, got %d", secs)
		t.Fatal(err)
	}

	if secs, err := parseTime("1m"); secs != 60 && err == nil {
		t.Fatalf("expected 60, got %d", secs)
		t.Fatal(err)
	}
}

// Test parsing of the input times using seconds.
func Test_parseTime_Seconds(t *testing.T) {
	if secs, err := parseTime("3600s"); secs != 3600 && err == nil {
		t.Fatalf("expected 3600, got %d", secs)
		t.Fatal(err)
	}
}

// Test parsing of the input containing an invalid specifier/no specifer.
func Test_parseTime_InvalidSpecifier(t *testing.T) {
	if _, err := parseTime("80l"); err == nil {
		t.Fatal("expected error, got nil (invalid specifier)")
	}

	if _, err := parseTime("80"); err == nil {
		t.Fatal("expected error, got nil (no specifier)")
	}
}

// Tests the generation of the complete config from the partial one.
func Test_getCompleteConfig(t *testing.T) {
	testPartialConfig := make(map[string]time.Time)

	testPartialConfig["hello"] = time.Now().Round(0)
	testPartialConfig["there"] = time.Now().Add(3 * time.Second).Round(0)
	testPartialConfig["friend"] = time.Now().Add(33 * time.Minute).Round(0)

	testCompleteConfig := getCompleteConfig(testPartialConfig, 3600)

	// test creation times
	for f, times := range testCompleteConfig {
		if testCompleteConfig[f][0] != testPartialConfig[f] {
			t.Fatalf("creation time should be the same (%q != %q)", times[0], testPartialConfig[f])
		}

		// test expiration times
		if testCompleteConfig[f][1] != testPartialConfig[f].Add(3600*time.Second) {
			t.Fatalf("expiration time incorrectly incremented: should be %q, is %q", testPartialConfig[f].Add(3600*time.Second), times[1])
		}
	}
}

// Tests the building of paths from file names.
func Test_buildPaths(t *testing.T) {
	home, _ := os.UserHomeDir()
	ssh := "/.ssh/"
	names := []string{"testing", "yo", "nayo", home + ssh + "shouldnotchange"}

	paths := buildPaths(names)
	for i, path := range paths {
		if path != home+ssh+names[i] && !strings.Contains(names[i], ssh) {
			t.Fatalf("path built incorrectly: expected %s, got %s", home+ssh+names[i], path)
		}
	}
}
