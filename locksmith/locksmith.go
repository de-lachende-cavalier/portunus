// Package locksmith is responsible of all the functionality associated with SSH keys (creating new ones,
// writing them to the correct location and with the correct encoding, etc.)
package locksmith

import (
	"time"

	"github.com/mowzhja/portunus/librarian"
)

// Changes keys, deleting the old ones and substituting them with new ones (the names associated
// with the various key files are kept the same).
//
// expired is an array of strings containing the filenames of the various private keys (pubkey files
// have the same name with the .pub extension added).
//
// N.B.: RotateKeys() assumes that all the files it receives as input have expired,
// checking whether keys have expired happens elsewhere.
func RotateKeys(expiredPaths []string, cipher string) (map[string]time.Time, error) {
	updatedData := make(map[string]time.Time)
	var err error

	err = librarian.DeleteKeyFiles(expiredPaths)
	if err != nil {
		return nil, err
	}

	for _, path := range expiredPaths {
		updatedData[path] = time.Now() // set creation date
		privKey, pubKey, err := genKeyPair(cipher)
		if err != nil {
			return nil, err
		}

		privBytes, err := encodePrivPEM(privKey)
		pubBytes, err := encodePubSSH(pubKey)

		err = librarian.WritePrivKey(privBytes, path)
		if err != nil {
			return nil, err
		}

		err = librarian.WritePubKey(pubBytes, path+".pub")
		if err != nil {
			return nil, err
		}
	}

	return updatedData, nil
}
