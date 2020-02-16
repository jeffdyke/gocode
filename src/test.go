package main

import (
  "github.com/jeffdyke/jssh"
  "log"
)


func main() {
  var jump = jssh.PublicKeyConnection{jssh.ConnectionInfo{"jeff", "jump.bondlink.org"}}
  c, e := jump.connect()
  if e != nil {
    log.Panicf("what the fuck is %v", e)
  }
  runCommand(*c, "ls -la ")
}
