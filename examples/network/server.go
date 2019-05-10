package main

import (
	"fmt"
	"net"
	"os"

	"github.com/tensortask/spice"
)

const (
	connHost = "localhost"
	connPort = "6666"
	connType = "tcp"
)

var (
	nonce     = [24]byte{58, 230, 79, 44, 187, 45, 107, 226, 245, 53, 169, 118, 218, 116, 235, 95, 132, 127, 166, 200, 203, 141, 251, 51}
	sharedKey = &[32]byte{199, 156, 103, 110, 157, 5, 107, 139, 94, 138, 53, 214, 74, 100, 211, 97, 106, 48, 11, 179, 200, 19, 244, 108, 138, 167, 49, 163, 156, 176, 66, 64}
)

func main() {
	// Listen for incoming connections.
	listener, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		panic(err)
	}
	// Close the listener when the application closes.
	defer listener.Close()
	fmt.Println("Listening on " + connHost + ":" + connPort)
	for {
		// Listen for an incoming connection.
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	decryptor := spice.NewDecryptor(os.Stdout, nonce, sharedKey)
	err := decryptor.DecryptFromReader(conn)
	if err != nil {
		panic(err)
	}
	conn.Close()
}
