package main

import (
	"net"

	"github.com/tensortask/spice"
)

const (
	connHost = "localhost"
	connPort = "6666"
	connType = "tcp"
	fileName = "testing/a_midsummers_night's_dream.txt"
)

var (
	nonce     = [24]byte{58, 230, 79, 44, 187, 45, 107, 226, 245, 53, 169, 118, 218, 116, 235, 95, 132, 127, 166, 200, 203, 141, 251, 51}
	sharedKey = &[32]byte{199, 156, 103, 110, 157, 5, 107, 139, 94, 138, 53, 214, 74, 100, 211, 97, 106, 48, 11, 179, 200, 19, 244, 108, 138, 167, 49, 163, 156, 176, 66, 64}
)

func main() {

	// connect to this socket
	conn, err := net.Dial(connType, connHost+":"+connPort)
	if err != nil {
		panic(err)
	}

	encryptor := spice.NewEncryptor(conn, nonce, sharedKey)
	encryptor.EncryptFile(fileName)

}
