package locksmith

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"errors"
)

// Generates the private key and the corresponding public key for either RSA (4096 bits) or Ed25519.
func genKeyPair(cipher string) (crypto.PrivateKey, crypto.PublicKey, error) {
	switch cipher {
	case "rsa":
		privKey, err := rsaPrivKey()
		if err != nil {
			return nil, nil, err
		}
		pubKey := privKey.Public()
		return privKey, pubKey, nil
	case "ed25519":
		pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, nil, err
		}
		return privKey, pubKey, nil
	default:
		return nil, nil, errors.New("Cipher not supported.")
	}
}

// Helper function, to keep the genKeyPair() function cleaner.
func rsaPrivKey() (*rsa.PrivateKey, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	return privKey, nil
}
