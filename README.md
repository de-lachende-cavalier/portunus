# About
portunus is a CLI utility for managing the expiration of ssh keys, written in Go.
To use it as intended, it should be built and then run with the correct CLI flags directly from your ~/.bashrc (or equivalent) file.

## About the name...
[Portunus](https://en.wikipedia.org/wiki/Portunus_(mythology)) was the the Roman god of keys (among other things).

## Structure
- librarian/ -> this is where the code for interfacing with the file system lives
- locksmith/ -> this is where the code responsible for key substitution and generation lives

## Useful references
- [Generate SSH keypair in native Go](https://gist.github.com/devinodaniel/8f9b8a4f31573f428f29ec0e884e6673)
- [Cobra.dev](https://cobra.dev/)
- [Automate SSH key generation](https://nathanielhoag.com/blog/2014/05/26/automate-ssh-key-generation-and-deployment/)
