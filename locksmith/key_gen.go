package locksmith

import (
	"os/exec"
	"errors"
)

// Wrapper around ssh-keygen, generates private/public key pair with the given password.
//
// N.B.: It is imperative that the file identified by path doesn't exist!
func genKeyPair(cipher string, passwd string, path string) error {
	// XXX the number of rounds was set to a mostly arbitrary value (i tested a
	// XXX bunch of them and found that 50 is reasonably fast)
	args := []string{"-q", "-t", cipher, "-N", passwd, "-f", path, "-a", "50"}

	switch cipher {
	case "rsa":
		args = append(args, "-b", "4096")
		return execCommand(args...)
	case "ed25519":
		return execCommand(args...)
	default:
		return errors.New("Cipher not supported.")
	}
}

// Helper function, to separate command execution from the definition of the various flags.
func execCommand(args ...string) error {
    return exec.Command("ssh-keygen", args...).Run()
}
