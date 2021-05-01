# About
portunus is a CLI utility for managing the expiration of ssh keys, written in Go.
It acts as a daemon, periodically checking whether keys have expired and notifying the user, who is then encouraged to regenerate new keys.

## About the name...
[Portunus](https://en.wikipedia.org/wiki/Portunus_(mythology)) was the the Roman god of keys (among other things).

## Structure
- tracker/ -> this is where the code for tracking key lifetime lives
- locksmith/ -> this is where the code responsible for key substitution and generation lives

## Useful references
- [Generate SSH keypair in native Go](https://gist.github.com/devinodaniel/8f9b8a4f31573f428f29ec0e884e6673)
- [Four Steps to Daemonize Your Go Programs](https://ieftimov.com/post/four-steps-daemonize-your-golang-programs/)
