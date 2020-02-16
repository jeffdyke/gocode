package main

import (
	"github.com/jeffdyke/jssh"
	"log"
	"flag"
	"os/user"
)



func main() {
	var u, _ = user.Current()
	var usr = flag.String("user", u.Username, "Defaults to your login name" )
	var host = flag.String("host", "", "Required must specify host, if using bastion see that help")
	//var bastion = flag.String("bastion", "", "Required if host")
	flag.Parse()
	var jump = jssh.PublicKeyConnection{ConnectionInfo:jssh.ConnectionInfo{User: *usr, Host: *host}}
	c, e := jump.Connect()
	if e != nil {
		log.Panicf("what the fuck is %v", e)
	}
	result := jssh.RunCommand(*c, "ls -la")
	log.Printf("home dir from staging %v\n", result)
}

