// Package locksmith is responsible of all the functionality associated with SSH keys (creating new ones,
// writing them to the correct location and with the correct encoding, etc.)
package locksmith

import (
	"time"
)

// Changes keys, deleting the old ones and substituting them with new ones (the names associated
// with the various key files are kept the same).
//
// expired is an array of strings containing the filenames of the various private keys (pubkey files
// have the same name with the .pub extension added).
//
// N.B.: ChangeKeys() assumes that all the files it receives as input have expired,
// checking whether keys have expired or not is up to the tracker package.
func ChangeKeys(expired []string, cipher string) (map[string]time.Time, error) {
	var err error

	updatedData := make(map[string]time.Time)

	expiredPaths, err := buildAbsPaths(expired)
	if err != nil {
		return nil, err
	}

	err = deleteKeyFiles(expiredPaths)
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

		err = writePrivKey(privBytes, path)
		if err != nil {
			return nil, err
		}

		err = writePubKey(pubBytes, path)
		if err != nil {
			return nil, err
		}
	}

	return updatedData, nil
}
