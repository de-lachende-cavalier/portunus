package librarian

import (
	"os"
	"reflect"
	"time"

	"testing"
)

// Tests reading and writing of the config file, if we use valid data.
func Test_readWriteConfig(t *testing.T) {
	curConfig := make(map[string][2]time.Time)

	curConfig["hello"] = [2]time.Time{time.Now().Round(0),
		time.Now().Add(44 * time.Minute).Round(0)}
	curConfig["friend"] = [2]time.Time{time.Now().Round(0),
		time.Now().Add(10 * time.Second).Round(0)}
	curConfig["leave"] = [2]time.Time{time.Now().Round(0),
		time.Now().Add(3 * time.Hour).Round(0)}

	err := WriteConfig(curConfig)
	if err != nil {
		t.Fatal(err)
	}

	readConfig, err := ReadConfig()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(curConfig, readConfig) {
		t.Fatalf("the config data written and the data read don't match up: expected %q, got %q", curConfig, readConfig)
	}

	// clean up
	err = os.Remove(configFile)
	if err != nil {
		t.Fatal(err)
	}
}

// Tests reading and writing of the config file, if we use invalid data.
func Test_readWriteConfig_InvalidData(t *testing.T) {
	invalidConfig := make(map[string][2]time.Time)

	invalidConfig["lol"] = [2]time.Time{time.Now().Add(1 * time.Second), time.Now()}

	err := WriteConfig(invalidConfig)
	if err == nil {
		t.Fatal("expected error on invalid config data, but no error was returned")
	}
}
