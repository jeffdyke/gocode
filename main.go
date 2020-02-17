package main

import (
	"flag"
	"github.com/jeffdyke/jssh"
	"golang.org/x/crypto/ssh"
	"log"
	"os/user"
	"strings"
)



func main() {
	var u, _ = user.Current()
	var usr = flag.String("user", u.Username, "Defaults to your login name" )
	var host = flag.String("host", "", "Required must specify host, if using bastion see that help")
	var bastion = flag.String("bastion", "", "Required if host")

	flag.Parse()
	var sClient *ssh.Client
	var err  error

	if *bastion != ""  && *host != "" {
		if strings.EqualFold(*bastion, *host) {
			log.Fatalf("-host(%v) and -bastion(%v) can not be the same", *host, *bastion)
		}
		log.Printf("Using %v to run command on %v", *bastion, *host)
		sClient, err = jssh.BastionConnect(*usr, *host, *bastion)
	} else if *host != "" {
		sClient, err = jssh.PublicKeyConnect(*usr, *host)
	} else {
		log.Fatal("Usage go main.go -host [-user] -bastion")
	}


	if err != nil {
		log.Panicf("what the fuck is %v", err)
	}

	cmds := []string{"pwd", "pwd", "hostname", "echo 'You'"}
	result := jssh.RunCommands(*sClient, cmds)
	log.Printf("home dir from staging %v\n", result)
}

