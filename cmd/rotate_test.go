package cmd

import (
	"io/ioutil"
	"bytes"
	"os"
	"time"
	s "strings"
	"fmt"
	"testing"

	"github.com/mowzhja/portunus/librarian"
)

// Helper function, check whether the correct cipher has been used in overwriting the keys.
func usedCorrectCipher(keyFile string, cipher string) bool {
	bytes, err := ioutil.ReadFile(keyFile + ".pub")
	if err != nil {
		fmt.Println(err)
		return false
	}
	text := string(bytes)

	// ssh-<cipher> ...
	if s.Contains(text, "ssh-"+cipher) {
		return true
	}

	return false
}

// Tests key renewal for RSA.
func Test_rotateCmd_RSA(t *testing.T) {
	confPath := os.Getenv("HOME") + "/.portunus_data.gob"

	if _, err := os.Stat(confPath); err == nil {
		t.Fatal("config file should not exist at this point")
	}

	cmd := rootCmd

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"rotate", "-t", "30m", "-c", "rsa"})

	err := cmd.Execute()
	if err != nil {
		t.Fatal(err)
	}

	rotatedConf, err := librarian.ReadConfig()
	if err != nil {
		t.Fatal(err)
	}

	for keyFile, times := range rotatedConf {
		// test times
		if times[0].Add(30 * time.Minute) != times[1] {
			t.Fatalf("wrong expiration time: expected %q, got %q", 
			times[0].Add(30 * time.Minute),
			times[1])
		}

		// test cipher
		if !usedCorrectCipher(keyFile, "rsa") {
			t.Fatal("rotate used the wrong cipher (was supposed to be rsa)")
		}
	}

	// the config file is created directly by the rotateCmd
	err = cleanupTestConfig()
	if err != nil {
		t.Fatal(err)
	}
}

// Tests key renewal for Ed25519
func Test_rotateCmd_Ed25519(t *testing.T) {
	confPath := os.Getenv("HOME") + "/.portunus_data.gob"

	if _, err := os.Stat(confPath); err == nil {
		t.Fatal("config file should not exist at this point")
	}

	cmd := rootCmd

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"rotate", "-t", "30m", "-c", "ed25519"})

	err := cmd.Execute()
	if err != nil {
		t.Fatal(err)
	}

	rotatedConf, err := librarian.ReadConfig()
	if err != nil {
		t.Fatal(err)
	}

	for keyFile, times := range rotatedConf {
		// test times
		if times[0].Add(30 * time.Minute) != times[1] {
			t.Fatalf("wrong expiration time: expected %q, got %q", 
			times[0].Add(30 * time.Minute),
			times[1])
		}

		// test cipher
		if !usedCorrectCipher(keyFile, "ed25519") {
			t.Fatal("rotate used the wrong cipher (was supposed to be rsa)")
		}
	}

	// the config file is created directly by the rotateCmd
	err = cleanupTestConfig()
	if err != nil {
		t.Fatal(err)
	}
}
