package locksmith

import (
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"reflect"
	"testing"

	"golang.org/x/crypto/ssh"
)

// Tests encoding of a private RSA key to PEM format.
func Test_encodePrivPEM_RSA(t *testing.T) {
	privK, _, err := genKeyPair("rsa")
	if err != nil {
		t.Fatal(err)
	}

	privBytes, err := encodePrivPEM(privK)
	if err != nil {
		t.Fatal(err)
	}

	block, _ := pem.Decode(privBytes)
	if block == nil {
		t.Fatal("the PEM block shouldn't be empty")
	}

	if len(block.Headers) != 0 {
		t.Fatalf("the block shouldn't have any headers: expected len = 0, got %d", len(block.Headers))
	}

	if block.Type != "OPENSSH PRIVATE KEY" {
		t.Fatalf("the block should have the correct type: expected OPENSSH PRIVATE KEY, got %q", block.Type)
	}

	privDER, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		t.Fatal(err)
	}

	// this part is tested in key_gen_test.go
	privKRSA, _ := privK.(*rsa.PrivateKey)

	privDERRSA, ok := privDER.(*rsa.PrivateKey)
	if !ok {
		t.Fatalf("the private key obtained from PEM file is of incorrect type: expected *rsa.PrivateKey, got %T", privDERRSA)
	}

	if !reflect.DeepEqual(privKRSA, privDERRSA) {
		t.Fatalf("the private RSA key generated initially and the one obtained from the PEM file are not the same")
	}
}

// Tests encoding of a private Ed25519 key to PEM format.
func Test_encodePrivPEM_Ed25519(t *testing.T) {
	privK, _, err := genKeyPair("ed25519")
	if err != nil {
		t.Fatal(err)
	}

	privBytes, err := encodePrivPEM(privK)
	if err != nil {
		t.Fatal(err)
	}

	block, _ := pem.Decode(privBytes)
	if block == nil {
		t.Fatal("the PEM block shouldn't be empty")
	}

	if len(block.Headers) != 0 {
		t.Fatalf("the block shouldn't have any headers: expected len = 0, got %d", len(block.Headers))
	}

	if block.Type != "OPENSSH PRIVATE KEY" {
		t.Fatalf("the block should have the correct type: expected OPENSSH PRIVATE KEY, got %q", block.Type)
	}

	privDER, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		t.Fatal(err)
	}

	// this part is tested in key_gen_test.go
	privKEd25519, _ := privK.(ed25519.PrivateKey)

	privDEREd25519, ok := privDER.(ed25519.PrivateKey)
	if !ok {
		t.Fatalf("the private key obtained from PEM file is of incorrect type: expected ed25519.PrivateKey, got %T", privDEREd25519)
	}

	if !privKEd25519.Equal(privDEREd25519) {
		t.Fatalf("the private Ed25519 key generated initially and the one obtained from the PEM file are not the same")
	}
}

// Tests encoding of a public RSA key to a compatible SSH format.
func Test_encodePubSSH_RSA(t *testing.T) {
	_, pubK, err := genKeyPair("rsa")
	if err != nil {
		t.Fatal(err)
	}

	pubBytes, err := encodePubSSH(pubK)
	if err != nil {
		t.Fatal(err)
	}

	pubSSH, _, _, _, err := ssh.ParseAuthorizedKey(pubBytes)
	if err != nil {
		t.Fatal(err)
	}

	if pubSSH.Type() != "ssh-rsa" {
		t.Fatalf("public key type incorrect: expected ssh-rsa, got %s", pubSSH.Type())
	}
}

// Tests encoding of a public Ed25519 key to a compatible SSH format.
func Test_encodePubSSH_Ed25519(t *testing.T) {
	_, pubK, err := genKeyPair("ed25519")
	if err != nil {
		t.Fatal(err)
	}

	pubBytes, err := encodePubSSH(pubK)
	if err != nil {
		t.Fatal(err)
	}

	pubSSH, _, _, _, err := ssh.ParseAuthorizedKey(pubBytes)
	if err != nil {
		t.Fatal(err)
	}

	if pubSSH.Type() != "ssh-ed25519" {
		t.Fatalf("public key type incorrect: expected ssh-ed25519, got %s", pubSSH.Type())
	}
}
