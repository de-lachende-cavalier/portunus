package locksmith

import (
	"time"
)

// Changes keys, deleting the old ones and substituting them with new ones (the names associated
// with the various key files are kept the same).
//
// expired is an array of strings containing the filenames of the various private keys (pubkey files
// have the same name with the .pub extension added).
func ChangeKeys(expired []string, cipher string) (map[string]time.Time, error) {
	var updatedData map[string]time.Time
	var err error

	expiredPaths, err := buildAbsPaths(expired)
	if err != nil {
		return nil, err
	}

	err = deleteKeyFiles(expiredPaths)
	if err != nil {
		return nil, err
	}

	for _, file := range expiredPaths {
		updatedData[file] = time.Now() // set creation date
		privKey, pubKey, err := genKeyPair(cipher)
		if err != nil {
			return nil, err
		}

		privBytes, err := encodePrivPEM(privKey)
		pubBytes, err := encodePubSSH(pubKey)

		err = writePrivKey(privBytes, file)
		if err != nil {
			return nil, err
		}

		err = writePubKey(pubBytes, file)
		if err != nil {
			return nil, err
		}
	}

	return updatedData, nil
}
