package locksmith

import (
	"crypto"
	"crypto/ssh"
	"crypto/x509"
	"encoding/pem"
)

// Encodes the RSA private key to valid PEM.
func encodePrivPEM(privKey crypto.PrivateKey) ([]byte, error) {
	privDER, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return nil, err
	}

	privBlock := pem.Block{
		Type:    "OPENSSH PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	privPEM := pem.EncodeToMemory(&privBlock)

	return privPEM, nil
}

// Encodes the RSA public key in the format SSH expects.
func encodePubSSH(pubBytes crypto.PublicKey) ([]byte, error) {
}
