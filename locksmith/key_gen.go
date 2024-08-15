package locksmith

import (
	"fmt"
	"os"
	"os/exec"
)

// Wrapper around ssh-keygen.
// Generates private/public key pair with the given password.
func genKeyPair(cipher string, passwd string, path string) error {
	// removing the previous files is crucial, becuase ssh-keygen prompts an overwrite otherwise
	os.Remove(path)
	os.Remove(path + ".pub")

	args := []string{"-q", "-t", cipher, "-N", passwd, "-f", path, "-a", "50"}

	switch cipher {
	case "rsa":
		args = append(args, "-b", "4096")
	case "ed25519":
	default:
		return fmt.Errorf("cipher not supported: %s", cipher)
	}

	cmd := exec.Command("ssh-keygen", args...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}
	if _, err := os.Stat(path + ".pub"); os.IsNotExist(err) {
		return err
	}

	return nil
}
