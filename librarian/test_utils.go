package librarian

import (
	"fmt"
	"os"
)

// Helper function, creates three random files in /tmp.
func createTestFiles() []string {
	var privPaths []string
	names := []string{"gonomolo", "hyperion", "super_private"}
	base := "/tmp/"

	for _, name := range names {
		privPaths = append(privPaths, base+name)
	}

	for _, path := range privPaths {
		_, err := os.Create(path)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		_, err = os.Create(path + ".pub")
		if err != nil {
			fmt.Println(err)
			return nil
		}
	}

	// check if the files have actually been created
	for _, path := range privPaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Printf("error creating %s", path)
			return nil
		}

		if _, err := os.Stat(path + ".pub"); os.IsNotExist(err) {
			fmt.Printf("error creating %s", path)
			return nil
		}
	}

	return privPaths
}
