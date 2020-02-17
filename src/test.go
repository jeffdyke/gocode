package main

import (
  "github.com/jeffdyke/ssh"
  "log"
)


func main() {
  var jump = ssh.PublicKeyConnection{ssh.ConnectionInfo{"jeff", "jump.bondlink.org"}}
  c, e := jump.connect()
  if e != nil {
    log.Panicf("what the fuck is %v", e)
  }
  runCommand(*c, "ls -la ")
}
