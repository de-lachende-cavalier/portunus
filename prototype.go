package main

import (
  "os"
  "os/signal"
  "syscall"
  "time"
  "log"
)

// XXX
// XXX the following code is just bad code throughout, but it's useful to have a very high level reference for the future
// XXX
func main() {
  info, _ := os.Stat("/Users/d0larhyde/.ssh/id_ed25519.pub")
  // fmt.Println(info.ModTime()) // this gets last modified time (i suppose that people don't modify their keys (why would they?))
  // ideally something that could be done is have portunus have its own user who's the only one to have access to the keys (apart from ssh-agent, obviously) so that 
  // other users are not even allowed to modify keys (but i can think about that later)
  log.SetOutput(os.Stdout) // step 1 from Four Steps to Daemonize... (check out README.md)

  sigChan := make(chan os.Signal, 1) // steps 2 and 3
  signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
  defer signal.Stop(sigChan)

  // XXX remember to go build instead of go run to be able to correctly handle signals
  go func() {
    for { // wrapping the goroutine in this for will make sure it will only exit on sigint/sigterm
      select {
      case s := <-sigChan:
        switch s {
        case syscall.SIGINT, syscall.SIGTERM:
          log.Println("SIGINT/SIGTERM => terminating...")
          os.Exit(1)
        case syscall.SIGHUP:
          log.Println("SIGHUP => reloading config (as per directions...)")
          // we actually do nothing here
        }
      }
    }
  }()

  // simulate portunus telling the user the key expired
  for {
    log.Println("The SSH key expired! Last modified:", info.ModTime())
    time.Sleep(2 * time.Second)
  }


  // step 4 requires use to use an init system to control our daemon => left for when the code is correctly implemented (shouldn't be too hard)
}
