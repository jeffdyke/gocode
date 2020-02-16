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

func RunCommand(client ssh.Client, cmds []string) string {
	var e error
	sess, e := client.NewSession()
	if e != nil {
		log.Fatalf("Could not create new Session %v", e)
	}
	defer sess.Close()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	sess.Stdout = &stdout
	sess.Stderr = &stderr

	stdin, err := sess.StdinPipe()
	if err != nil {
		log.Fatalf("Failed to open StdInPipe. %v ",err)
	}
	e = sess.Shell()
	if e != nil {
		log.Fatalf("Failed to create shell.  %v", e)
	}
	for _, cmd := range cmds {
		log.Printf("Running %v", cmd)
		_, e = fmt.Fprintf(stdin, "%s\n", cmd)
		if e != nil {
			log.Printf("Failed to run %s. Error: %v", cmd, e)
		}
	}
	log.Printf("StdOut %v", stdout.String())
	log.Printf("StdErr %v", stderr.String())
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

func BastionConnect(usr string, host string, bastion string)  (client *ssh.Client, err error){
	var conn = BastionConnectInfo{
		ConnectionInfo: ConnectionInfo{User: usr, Host: host},
		Bastion: bastion,
	}
	return conn.Connect()
}

func PublicKeyConnect(usr string, host string) (*ssh.Client, error) {
	var conn = PublicKeyConnection{ConnectionInfo:ConnectionInfo{User: usr, Host: host}}
	return conn.Connect()
}


func (info *BastionConnectInfo) Connect() (*ssh.Client, error) {
	var localAgent = sshAgentConnect()
	_ = clientAuth(info.User, ssh.PublicKeysCallback(localAgent.Signers))
	var sshAgent = sshAgentConnect()
	var config = clientAuth(info.User, ssh.PublicKeysCallback(sshAgent.Signers))
	sshc, err := ssh.Dial(TCP, formatHost(info.Bastion), config)
	if err != nil {
		log.Fatalf("Failed to connect to Bastion host %v\nError: %v", info.Bastion, err)
		return nil, err
	}
	lanConn, err := sshc.Dial(TCP, formatHost(info.Host))
	if err != nil {
		log.Fatalf("Failed to connect to %v\nError: %v", info.Host, err)
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
