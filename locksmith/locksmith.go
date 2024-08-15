// Package locksmith is responsible of all the functionality associated with SSH keys (creating new ones, writing them to the correct location and with the correct encoding, etc.)
package locksmith

import (
	"time"

	"github.com/de-lachende-cavalier/portunus/librarian"
)

// Changes keys, deleting the old ones and substituting them with new ones (the names associated with the various key files are kept the same).
func RotateKeys(expiredPaths []string, cipher string, passwd string) (map[string]time.Time, error) {
	updatedData := make(map[string]time.Time)

	err := librarian.DeleteKeyFiles(expiredPaths)
	if err != nil {
		return nil, err
	}

	for _, path := range expiredPaths {
		updatedData[path] = time.Now() // set creation date

		err := genKeyPair(cipher, passwd, path)
		if err != nil {
			return nil, err
		}
	}

	return updatedData, nil
}
