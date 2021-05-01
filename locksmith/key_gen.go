package locksmith

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"errors"
)

func genKeyPair(cipher string) (crypto.PrivateKey, crypto.PublicKey, error) {
	switch cipher {
	case "rsa":
		privKey, err := RSAPrivKey()
		if err != nil {
			return nil, nil, err
		}
		pubKey := privKey.PublicKey
		return privKey, pubKey, nil
	case "ed25519":
		privKey, pubKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, nil, err
		}
		return privKey, pubKey, nil
	default:
		return nil, nil, errors.New("Cipher not supported.")
	}
}

func RSAPrivKey() (*rsa.PrivateKey, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	return privKey, nil
}
