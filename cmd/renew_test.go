package cmd

import (
	"bytes"
	"testing"
	"time"

	"github.com/de-lachende-cavalier/portunus/librarian"
)

// Tests key renewal.
func Test_renewCmd(t *testing.T) {
	oldConf, err := writeTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	cmd := rootCmd

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"renew", "-t", "1h"})

	err = cmd.Execute()
	if err != nil {
		t.Fatal(err)
	}

	newConf, err := librarian.ReadConfig()
	if err != nil {
		t.Fatal(err)
	}

	// test creation times
	for file, times := range newConf {
		if oldConf[file][0] != times[0] {
			t.Fatalf("incorrect creation time for %s: expected %q, got %q",
				file,
				oldConf[file][0],
				times[0])
		}
	}

	// test expiration times
	for file, times := range newConf {
		if oldConf[file][1].Add(3600*time.Second) != times[1] {
			t.Fatalf("expiration time set incorrectly by renew for %s: expected %q, got %q",
				file,
				oldConf[file][1].Add(3600*time.Second),
				times[1])
		}
	}

	// cleanup
	err = cleanupTestConfig()
	if err != nil {
		t.Fatal(err)
	}
}
