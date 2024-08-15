# portunus

portunus is a CLI utility for managing the expiration of SSH keys, written in Go.

portunus automates the process of SSH key rotation, helping to maintain good security practices by regularly updating SSH keys. It's designed to be integrated into your shell startup script for seamless operation.

## Installation

To install portunus, clone this repository and build the binary:

```bash
git clone https://github.com/de-lachende-cavalier/portunus.git
cd portunus
go build
```

## Usage

Add the following line to your `~/.bashrc` (or equivalent) file:

```bash
/path/to/portunus [flags]
```

Replace `/path/to/portunus` with the actual path to your built binary, and `[flags]` with any necessary command-line flags.

## Project Structure

- `main.go`: Entry point of the application
- `cmd/`: Contains the Cobra command definitions
- `librarian/`: Houses the code for interfacing with the file system
- `locksmith/`: Contains the code responsible for key substitution and generation

## About the Name

[Portunus](https://en.wikipedia.org/wiki/Portunus_(mythology)) was the Roman god of keys, doors, and livestock. Our project borrows the name for its association with keys and security.

## Useful references

- [Generate SSH keypair in native Go](https://gist.github.com/devinodaniel/8f9b8a4f31573f428f29ec0e884e6673)
- [Cobra.dev](https://cobra.dev/)
- [Automate SSH key generation](https://nathanielhoag.com/blog/2014/05/26/automate-ssh-key-generation-and-deployment/)
