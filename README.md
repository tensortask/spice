<p align="center">
<a href="https://tensortask.com">
<img width="200" alt="TensorTask Logo" src="https://storage.googleapis.com/tensortask-static/tensortask_transparent.png">
</a>
</p>

[![GoDoc][1]][2] [![Go Report Card][3]][4] [![Keybase Chat][5]][6] [![Cloud Build][7]][8]

[1]: https://godoc.org/github.com/tensortask/spice?status.svg
[2]: https://godoc.org/github.com/tensortask/spice
[3]: https://goreportcard.com/badge/github.com/tensortask/spice
[4]: https://goreportcard.com/report/github.com/tensortask/spice
[5]: https://img.shields.io/badge/keybase%20chat-tensortask.public-blue.svg
[6]: https://keybase.io/team/tensortask.public
[7]: https://storage.googleapis.com/tensortask-static/build/spice.svg
[8]: https://github.com/sbsends/cloud-build-badge

[9]: http://nacl.cr.yp.to/
[10]: https://godoc.org/golang.org/x/crypto/nacl
[11]: https://crypto.stackexchange.com/questions/22435/public-key-encryption-and-big-files-with-nacl#answer-31554
[12]: https://www.imperialviolet.org/2014/06/27/streamingencryption.html
[13]: https://godoc.org/golang.org/x/crypto/nacl/secretbox

# 🌶 Spice: Blazingly Fast Chunk Encryption

```diff
- #############################
- #   NOT PRODUCTION READY    #
- #############################
```

This package has not undergone a security audit. Please DO NOT use this for mission critical comms just yet.

## Applications:
* Big Data
* P2P

## Encryption Library

The encryption library used is [NaCl][9], which leverages Curve25519, XSalsa20 and Poly1305. The implementation is from the high quality golang crypto package ([golang.org/x/crypto/nacl][10]).

**Why NaCl?**

* Super fast encryption/decryption
* Super fast signatures
* Authentication at the chunk level (signature verfication)
* Aysemmetric Encryption (Ideal for P2P)
* Large Nonce

**NaCl allows each chunk of data to be authenticated by the receiver, and encrypted with the recepient's public key. This enables true end-to-end encryption with guranteed authenticity.**

## Chunk-Level Replay Prevention
Chunked encryption is succeptable to chunk-level replay attacks. [Jack O'Connor][11] explains this well: 

* Bob sends two messages to Sarah. Let's say the first message comes in chunks A1 and A2, and the second message comes in chunks B1 and B2.
* Mallory intercepts both messages. She then constructs the new message A1+B2.
* Mallory sends the forged message to Sarah. The boxes open with Bob's key, and the sequence numbers look good.


In order to mitigate chunk-level replay, some groups have proposed using a incremented nonce. This ensures that each block must be decrypted in the correct order or else the decryption will **not** be verfied.

Instead of using a counter, spice links chunks of data by using the first 24 bytes of the last encrypted block. This is analogous to how a blockchain references the previous hash. Chunk referencing is also very efficient.

**A CRYPTOGRAPHICALLY RANDOM 24 BYTE ARRAY MUST BE GENERATED FOR THE FIRST NONCE!**

## 16kb Encrypted Chunk Size
Rationale: 

When sending data over the network, chunking is pretty much a given. TLS has a maximum record size of 16KB and this fits neatly with authenticated encryption APIs which all operate on an entire message at once.
>-[Adam Langley][11]



1. The whole message needs to be held in memory to be processed.
2. Using large messages pressures implementations on small machines to decrypt and process plaintext before authenticating it. This is very dangerous, and this API does not allow it, but a protocol that uses excessive message sizes might present some implementations with no other choice.
3. Fixed overheads will be sufficiently amortised by messages as small as 8KB.
4. Performance may be improved by working with messages that fit into data caches.
>-[Golang Developer's][13]

**Each data block is chunked in 16368 byte sections and sealed with a 16 byte signature. This results in encrypted chunks that are exactly 16kb. Optimal for web transmission and safe for small computers.**

## Usage

#### Encrypted Over The Wire

RUN THIS EXAMPLE (from spice root directory)!

1st terminal: `go run examples/network/server.go`

2nd terminal: `go run examples/network/client.go`

```golang
# server.go

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
```

```golang
# client.go

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
```

#### File Encryption


RUN THIS EXAMPLE (from spice root directory)!

1st terminal: `go run examples/file/main.go`

```golang
# main.go

package main

import (
	"os"

	"github.com/tensortask/spice"
)

func main() {
	hostKeys, err := spice.RandomKeyPair()
	if err != nil {
		panic(err)
	}
	peerKeys, err := spice.RandomKeyPair()
	if err != nil {
		panic(err)
	}
	sharedKey := hostKeys.SharedKey(peerKeys)
	nonce, err := spice.RandomNonce()
	if err != nil {
		panic(err)
	}

	encryptedFile, err := os.Create("testing/encrypted_shakespeare.txt")
	if err != nil {
		panic(err)
	}

	encryptor := spice.NewEncryptor(encryptedFile, nonce, sharedKey)
	err = encryptor.EncryptFile("testing/a_midsummers_night's_dream.txt")
	if err != nil {
		panic(err)
	}
	encryptedFile.Close()

	decryptedFile, err := os.Create("testing/decrypted_shakespeare.txt")
	if err != nil {
		panic(err)
	}

	decryptor := spice.NewDecryptor(decryptedFile, nonce, sharedKey)

	err = decryptor.DecryptFile("testing/encrypted_shakespeare.txt")
	if err != nil {
		panic(err)
	}
	decryptedFile.Close()
}
```

## Performance

NOTE: this is a very old metric...  recently spice has been clocking in much hotter. Full benchmarks to come.

2017 Macbook, 1.3Ghz i5, 8GB ram = 146.69 MB/s

Well above average consumer internet connection speeds ✅
