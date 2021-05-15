package cmd

import (
	"testing"
	"os"
	"bytes"
	"io"
  "time"

  "github.com/mowzhja/portunus/librarian"
)

// Tests the check command.
//
// code snippet from: https://stackoverflow.com/questions/10473800/in-go-how-do-i-capture-stdout-of-a-function-into-a-string
func Test_checkCmd(t *testing.T) {
	curConfig := make(map[string][2]time.Time)

	curConfig["hello"] = [2]time.Time{time.Now().Round(0),
		time.Now().Add(1 * time.Second).Round(0)}
	curConfig["friend"] = [2]time.Time{time.Now().Round(0),
		time.Now().Add(22 * time.Hour).Round(0)}
	curConfig["leave"] = [2]time.Time{time.Now().Round(0),
		time.Now().Add(2 * time.Minute).Round(0)}

	err := librarian.WriteConfig(curConfig)
	if err != nil {
		t.Fatal(err)
	}

	// to make sure only "hello" expires (so only "hello" is returned by GetExpiredKeys())
	time.Sleep(2 * time.Second)
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

	// clean up
	err = os.Remove(os.Getenv("HOME") + "/.portunus_data.gob")
	if err != nil {
		t.Fatal(err)
	}
}
