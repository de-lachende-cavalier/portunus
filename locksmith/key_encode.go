package locksmith

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"golang.org/x/crypto/ssh"
)

// Encodes the private key to valid PEM.
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

// Encodes the public key in the format SSH expects.
func encodePubSSH(pubKey crypto.PublicKey) ([]byte, error) {
	pubSSHKey, err := ssh.NewPublicKey(pubKey)
	if err != nil {
		return nil, err
	}

	pubBytes := ssh.MarshalAuthorizedKey(pubSSHKey)

	return pubBytes, nil
}
