# portunus

portunus is a CLI utility for managing the expiration of SSH keys, written in Go.

portunus automates the process of SSH key rotation, helping to maintain good security practices by regularly updating SSH keys. It's designed to be integrated into your shell startup script for seamless operation.

## Features

- **Key Rotation**: Automatically rotate SSH keys when they expire
- **Key Renewal**: Extend the expiration date of existing keys
- **Expiration Tracking**: Track and manage key expiration dates
- **Multiple Cipher Support**: Support for ed25519, RSA, and ECDSA keys

## Installation

### From Source

```bash
git clone https://github.com/de-lachende-cavalier/portunus.git
cd portunus
go build
```

### Using Go Install

```bash
go install github.com/de-lachende-cavalier/portunus@latest
```

## Usage

### Basic Commands

```bash
# Check for expired keys
portunus check

# Rotate expired keys
portunus rotate -t 30d -p "your-password"

# Renew expired keys
portunus renew -t 30d
```

### Shell Integration

Add the following line to your `~/.bashrc`, `~/.zshrc`, or equivalent:

```bash
/path/to/portunus check
```

This will check for expired keys each time you open a new shell.

### Command Options

#### Rotate Command

```
portunus rotate [flags]

Flags:
  -c, --cipher string       specifies which cipher to use for key generation (default "ed25519")
  -p, --password string     specifies the password to use with ssh-keygen
  -s, --subset strings      specifies the subset of keys you want to act on
  -t, --time string         specifies for how much longer the key should be valid
```

#### Renew Command

```
portunus renew [flags]

Flags:
  -s, --subset strings      specifies the subset of keys you want to act on
  -t, --time string         specifies for how much longer the key should be valid
```

#### Global Flags

```
      --config string       config file (default is $HOME/.portunus.json)
      --log-level string    log level (debug, info, warn, error) (default "info")
      --pretty-logs         enable pretty logging (default true)
```

## Project Structure

- `main.go`: Entry point of the application
- `cmd/`: Contains the Cobra command definitions
- `pkg/config/`: Configuration management
- `pkg/keys/`: SSH key management
- `pkg/logger/`: Structured logging

## About the Name

[Portunus](https://en.wikipedia.org/wiki/Portunus_(mythology)) was the Roman god of keys, doors, and livestock.