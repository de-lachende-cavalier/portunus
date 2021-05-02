package locksmith

import (
	"testing"
	"crypto/rsa"
	"crypto/ed25519"
)

func Test_genKeyPairRSA(t *testing.T) {
	privK, pubK, err := genKeyPair("rsa")
	if err != nil {
		t.Fatal(err)
	}

	// necessary steps => genKeyPair returns interfaces{}
	privRSA, ok := privK.(*rsa.PrivateKey)
	if !ok {
		t.Fatalf("the private RSA key is not of the correct type: expected *rsa.PrivateKey, got %T", privK)
	}

	pubRSA, ok := pubK.(*rsa.PublicKey)
	if !ok {
		t.Fatalf("the public RSA key is not of the correct type: expected *rsa.PublicKey, got %T", pubK)
	}

	if !pubRSA.Equal(privRSA.Public()) {
		t.Fatal("the public RSA key generated is not associated with the generated private key")
	}

	if (pubRSA.Size() * 8) != 4096 {
		t.Fatalf("the public RSA key is of the wrong size: expected 4096, got %d", pubRSA.Size() * 8)
	}
}

func Test_genKeyPairEd25519(t *testing.T) {
	privK, pubK, err := genKeyPair("ed25519")
	if err != nil {
		t.Fatal(err)
	}

	privEd25519, ok := privK.(ed25519.PrivateKey)
	if !ok {
		t.Fatalf("the private Ed25519 key is not of the correct type: expected ed25519.PrivateKey, got %T", privK)
	}

	pubEd25519, ok := pubK.(ed25519.PublicKey)
	if !ok {
		t.Fatalf("the public Ed25519 key is not of the correct type: expected ed25519.PublicKey, got %T", pubK)
	}

	if !pubEd25519.Equal(privEd25519.Public()) {
		t.Fatal("the public Ed25519 key generated is not associated with the generated private key")
	}
}

func Test_genKeyPairNonExisting(t *testing.T) {
	privK, pubK, err := genKeyPair("made up cipher")

	if privK != nil {
		t.Fatalf("private key should be nil when using a made up cipher: expected nil, got %T", privK)
	}

	if pubK != nil {
		t.Fatalf("public key should be nil when using a made up cipher: expected nil, got %T", pubK)
	}

	if err == nil {
		t.Fatalf("error should not be nil when using a made up cipher: expected 'Cipher no supported', got %T", err)
	}
}
