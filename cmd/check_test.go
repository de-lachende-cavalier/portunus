package cmd

import (
	"bytes"
	"io"
	"os"
	s "strings"
	"testing"
	"time"
)

// Tests the check command when all keys have expired.
func Test_checkCmd_AllExpired(t *testing.T) {
	_, err := writeTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	// to make sure all keys expire
	time.Sleep(6 * time.Second)

	// code snippet from: https://stackoverflow.com/questions/10473800/in-go-how-do-i-capture-stdout-of-a-function-into-a-string
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	runCheckCmd(checkCmd, []string{})

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC

	if !s.Contains(out, "Either renew or rotate them!") &&
		!(s.Contains(out, "hello") && s.Contains(out, "friend") && s.Contains(out, "leave")) {
		t.Fatal("the check command did not detect expiration correctly")
	}

	err = cleanupTestConfig()
	if err != nil {
		t.Fatal(err)
	}
}

// Tests the check command when none of the keys have expired.
func Test_checkCmd_NoneExpired(t *testing.T) {
	_, err := writeTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	// code snippet from: https://stackoverflow.com/questions/10473800/in-go-how-do-i-capture-stdout-of-a-function-into-a-string
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	runCheckCmd(checkCmd, []string{})

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC

	if !(s.Contains(out, "The keys are still fresh")) {
		t.Fatal("the check command did not detect expiration correctly")
	}

	err = cleanupTestConfig()
	if err != nil {
		t.Fatal(err)
	}
}

// Tests the check command when some keys have expired.
func Test_checkCmd_SomeExpired(t *testing.T) {
	_, err := writeTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	// to make sure only one key expires (namely "hello")
	time.Sleep(2 * time.Second)

	// code snippet from: https://stackoverflow.com/questions/10473800/in-go-how-do-i-capture-stdout-of-a-function-into-a-string
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	runCheckCmd(checkCmd, []string{})

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC

	if !(s.Contains(out, "Either renew or rotate them!") &&
		s.Contains(out, "hello") &&
		!s.Contains(out, "friend") &&
		!s.Contains(out, "leave")) {
		t.Fatal("the check command did not detect expiration correctly")
	}

	err = cleanupTestConfig()
	if err != nil {
		t.Fatal(err)
	}
}
