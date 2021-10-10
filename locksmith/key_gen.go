package locksmith

import (
	"errors"
	"os/exec"
)

// Wrapper around ssh-keygen, generates private/public key pair with the given password.
func genKeyPair(cipher string, passwd string, path string) error {
	args := []string{"-q", "-t", cipher, "-N", passwd, "-f", path, "-a", "50"}

	switch cipher {
	case "rsa":
		args = append(args, "-b", "4096")
		return execCommand(args...)
	case "ed25519":
		return execCommand(args...)
	default:
		return errors.New("cipher not supported")
	}
}

// Helper function, to separate command execution from the definition of the various flags.
func execCommand(args ...string) error {
	return exec.Command("ssh-keygen", args...).Run()
}
