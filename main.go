package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	s := "gopher"
	fmt.Println("Hello and welcome, %s!", s)

	args := os.Args

	fmt.Println(args)

	if len(args) < 2 {
		os.Exit(0)
	}

	switch args[1] {
	case "-h", "--help":
		fmt.Println("HELP COMMAND")
		os.Exit(0)
	case "forward":
		if len(args) != 4 {
			fmt.Println("Incorrect number of arguments")
			os.Exit(1)
		}

		o := args[2]
		d := args[3]

		//helpers.ParseForwardingArgs(o, d)

		err := Forward(o, d)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}

		os.Exit(0)

	default:
		fmt.Println("Command not found")
		os.Exit(0)
	}
}

func Forward(c string, d string) error {
	listener, err := net.Listen("tcp", c)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", c, err)
		return err
	}
	defer func() {
		err = listener.Close()
		if err != nil {
			log.Fatalf("Failed to close %s: %v", c, err)
		}
	}()
	log.Printf("Port forwarding server started on %s, forwarding to %s", c, d)

	for {
		// Accept incoming connections
		var conn net.Conn
		conn, err = listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		// Handle the connection in a new goroutine
		go handleConnection(conn, d)
	}
}

func handleConnection(localConn net.Conn, remoteAddr string) {
	// Connect to the remote server
	remoteConn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		log.Printf("Failed to connect to remote address %s: %v", remoteAddr, err)
		localConn.Close()
		return
	}

	// Start bidirectional copy
	go func() {
		defer localConn.Close()
		defer remoteConn.Close()
		_, err := io.Copy(remoteConn, localConn)
		if err != nil {
			log.Printf("Error while copying from local to remote: %v", err)
		}
	}()

	go func() {
		defer localConn.Close()
		defer remoteConn.Close()
		_, err := io.Copy(localConn, remoteConn)
		if err != nil {
			log.Printf("Error while copying from remote to local: %v", err)
		}
	}()
}
