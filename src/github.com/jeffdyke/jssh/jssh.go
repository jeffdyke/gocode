package jssh


import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/user"
	"strings"
)

const TCP  = "tcp"
const PORT = "22"
const SOCKET = "SSH_AUTH_SOCK"
const UNIX = "unix"


func sshAgentConnect() agent.ExtendedAgent {
	socket := os.Getenv(SOCKET)
	conn, err := net.Dial(UNIX, socket)
	if err != nil {
		log.Fatalf("Failed to connect to %v", SOCKET)
	}
	agentClient := agent.NewClient(conn)
	return agentClient
}

func clientAuth(usr string, auth ssh.AuthMethod) *ssh.ClientConfig {
	config := &ssh.ClientConfig{
		User: usr,
		Auth: []ssh.AuthMethod{
			auth,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return config
}

func RunCommand(client ssh.Client, cmds string) string {
	var ls = strings.Split(cmds, ",")
	var session *ssh.Session
	session, _ = client.NewSession()



	for _, cmd := range ls {

		var o bytes.Buffer
		var e bytes.Buffer
		session.Stdout = &o
		session.Stderr = &e



		if err:= session.Run(cmd); err != nil {
			panic("Failed to run " + err.Error())

		} else {
			log.Printf("Successfully ran %v", cmd)
			log.Printf("StdOut if Any: %v", session.Stdout)
		}
	}
	defer session.Close()

	return fmt.Sprintf("%v completed", cmds)
}
func formatHost(host string) string {
	return fmt.Sprintf("%s:%s", host, PORT)
}


type ConnectionInfo struct {
	User string
	Host string
}
type PublicKeyConnection struct {
	ConnectionInfo
}

type BastionConnectInfo struct {
	ConnectionInfo
	Bastion string

}

func (info *BastionConnectInfo) Connect() (*ssh.Client, error) {
	var localAgent = sshAgentConnect()
	clientAuth(info.User, ssh.PublicKeysCallback(localAgent.Signers))
	var sshAgent = sshAgentConnect()
	var config = clientAuth(info.User, ssh.PublicKeysCallback(sshAgent.Signers))
	sshc, err := ssh.Dial(TCP, formatHost(info.Bastion), config)
	if err != nil {
		log.Panicf("Failed to connect to Bastopm host %v\nError: %v", info.Bastion, err)
		return nil, err
	}
	lanConn, err := sshc.Dial(TCP, formatHost(info.Host))
	if err != nil {
		log.Panicf("Failed to connect to %v\nError: %v", info.Host, err)
		return nil, err
	}
	ncc, chans, reqs, err := ssh.NewClientConn(lanConn, formatHost(info.Host), config)
	if err != nil {
		fmt.Printf("got error trying to get new client connection %v\n -- %v\n", formatHost(info.Host), err)
		return nil, err
	}

	sClient := ssh.NewClient(ncc, chans, reqs)
	return sClient, nil
}

func (info PublicKeyConnection) Connect() (*ssh.Client, error)  {
	var usr, _ = user.Current()
	key, err := ioutil.ReadFile(fmt.Sprintf("%v/.ssh/id_rsa", usr.HomeDir))
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
		return nil, err
	}
	var config = clientAuth(info.User, ssh.PublicKeys(signer))

	client, err := ssh.Dial(TCP, formatHost(info.Host), config)

	return client, err

}
