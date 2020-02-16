package main

import (
	"github.com/jeffdyke/jssh"
	"golang.org/x/crypto/ssh"
	"log"
	"flag"
	"os/user"
)

func connect(ci jssh.ConnectionInfo, bastion string, cmds string) (*ssh.Client, error) {
	var conn = jssh.BastionConnectInfo{
		ConnectionInfo: ci,
		Bastion: bastion,
	}
	return conn.Connect()

}

func parseArgs(u string, h, string, b string) *ssh.Client {
	var auth jssh.ConnectionInfo
	if b != "" {
		var auth = jssh.BastionConnectInfo{
			ConnectionInfo: jssh.ConnectionInfo{User: u, Host: h},
			Bastion:        b,
		}
		auth.Connect()
	} else {
		var auth = jssh.PublicKeyConnection{ConnectionInfo: jssh.ConnectionInfo{User: u, Host: h}
		auth.connect()
	}


}
func main() {
	var u, _ = user.Current()
	var usr = flag.String("user", u.Username, "Defaults to your login name" )
	var host = flag.String("host", "fuckyou.foo.com", "Required must specify host, if using bastion see that help")
	var bastion = flag.String("bastion", "", "Required if host")
	flag.Parse()
	var connObj =

	if *bastion != "" {
		log.Printf("Using %v to run command on %v", bastion, host)
		bastion()
	} else {
		var conn = jssh.PublicKeyConnection{ConnectionInfo:jssh.ConnectionInfo{User: *usr, Host: *host}}
	}

	if conn != nil {
		log.Panicf("what the fuck is %v", e)
	}
	c, e := conn.Connect()
	if e != nil {
		log.Panicf("what the fuck is %v", e)
	}
	result := jssh.RunCommand(*c, "ls -la")
	log.Printf("home dir from staging %v\n", result)
}

