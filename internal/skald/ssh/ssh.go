package ssh

import (
	"fmt"
	"log"

	"golang.org/x/crypto/ssh"
)

func NewSSHConnection() {
	// Define the configuration for the SSH connection
	config := &ssh.ClientConfig{
		User: "valnix", // Replace with your username
		Auth: []ssh.AuthMethod{
			ssh.Password("valnix"), // Replace with your password
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Not recommended for production use
	}

	// Connect to the server
	client, err := ssh.Dial("tcp", "localhost:222", config) // Replace with your server address and port
	if err != nil {
		log.Fatalf("Failed to dial: %s", err)
	}
	defer client.Close()

	// Execute a command on the remote server
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %s", err)
	}
	defer session.Close()

	out, err := session.Output("mkdir valnix")
	if err != nil {
		fmt.Println("Cannot run command")
	}
	fmt.Println(string(out))
}
