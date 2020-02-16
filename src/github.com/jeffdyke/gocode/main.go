package main

import (
	"github.com/jeffdyke/jssh"
	"log"
)



func main() {
	var jump = jssh.PublicKeyConnection{ConnectionInfo:jssh.ConnectionInfo{User: "jeff", Host: "jump.bondlink.com"}}
	c, e := jump.Connect()
	if e != nil {
		log.Panicf("what the fuck is %v", e)
	}
	result := jssh.RunCommand(*c, "ls -la")
	log.Printf("home dir from staging %v\n", result)
}

