// Package tracker contains all the functionality that allows tracking files in various directories.
package librarian

// Reads encoded data relating to files (config data is encrypted and then hex encoded).
// what about key? => derived from a PBKDF => each time the keys are rotated a new passphrase is given to the user (obvs this passphrase is not stored anywhere)
func readConfig() {
}

// Checks that the data makes sense
func checkIntegrity() {
}

