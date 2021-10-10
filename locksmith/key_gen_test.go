package locksmith

import (
	"os"
	"reflect"
	"testing"

	"golang.org/x/crypto/ssh"
)

// Tests the generation of RSA keys.
func Test_genKeyPair_RSA(t *testing.T) {
	var paths = []string{"/tmp/key1", "/tmp/key2"}
	passwd := []byte("rsa_test")

	for _, path := range paths {
		err := genKeyPair("rsa", string(passwd), path)
		if err != nil {
			t.Fatal(err)
		}
	}

	// check that the keys have been created correctly
	err := checkPathsExist(paths)
	if err != nil {
		t.Fatal(err)
	}

	// check correct algorithm usage
	for _, path := range paths {
		privBytes, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}

		signer, err := ssh.ParsePrivateKeyWithPassphrase(privBytes, passwd)
		if err != nil {
			t.Fatal(err)
		}
		// the public key derived from the private one
		pubSign := signer.PublicKey()

		pubBytes, err := os.ReadFile(path + ".pub")
		if err != nil {
			t.Fatal(err)
		}

		pubKey, _, _, _, err := ssh.ParseAuthorizedKey(pubBytes)
		if err != nil {
			t.Fatal(err)
		}

		if pubKey.Type() != "ssh-rsa" {
			t.Fatalf("type error for public key: expected ssh-rsa, got %s", pubKey.Type())
		}

		if !reflect.DeepEqual(pubKey, pubSign) {
			t.Fatal("the public key derived from the private on and the one read from the public key file don't match")
		}
	}

	// cleanup
	err = cleanupPaths(paths)
	if err != nil {
		t.Fatal(err)
	}
}

// Tests the generation of Ed25519 keys.
func Test_genKeyPair_Ed25519(t *testing.T) {
	var paths = []string{"/tmp/key1", "/tmp/key2"}
	passwd := []byte("ed25519_test")

	for _, path := range paths {
		err := genKeyPair("ed25519", string(passwd), path)
		if err != nil {
			t.Fatal(err)
		}
	}

	// check that the keys have been created correctly
	err := checkPathsExist(paths)
	if err != nil {
		t.Fatal(err)
	}

	// check correct algorithm usage
	for _, path := range paths {
		privBytes, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}

		signer, err := ssh.ParsePrivateKeyWithPassphrase(privBytes, passwd)
		if err != nil {
			t.Fatal(err)
		}
		// the public key derived from the private one
		pubSign := signer.PublicKey()

		pubBytes, err := os.ReadFile(path + ".pub")
		if err != nil {
			t.Fatal(err)
		}

		pubKey, _, _, _, err := ssh.ParseAuthorizedKey(pubBytes)
		if err != nil {
			t.Fatal(err)
		}

		if pubKey.Type() != "ssh-ed25519" {
			t.Fatalf("type error for public key: expected ssh-ed25519, got %s", pubKey.Type())
		}

		if !reflect.DeepEqual(pubKey, pubSign) {
			t.Fatal("the public key derived from the private one and the one read from the public key file don't match")
		}
	}

	// cleanup
	err = cleanupPaths(paths)
	if err != nil {
		t.Fatal(err)
	}
}
