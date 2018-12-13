package spice_test

import (
	"io"
	"os"
	"testing"

	"github.com/tensortask/spice"
)

func TestEncrypt(t *testing.T) {
	hostKeys, err := spice.RandomKeypair()
	if err != nil {
		t.Error(err)
	}
	peerKeys, err := spice.RandomKeypair()
	if err != nil {
		t.Error(err)
	}
	nonce, err := spice.RandomNonce()
	if err != nil {
		t.Error(err)
	}
	sharedKey := hostKeys.SharedKey(peerKeys.Public)

	outfile, _ := os.Create("big_copy.txt")

	encryptor := spice.NewEncryptor(outfile, nonce, sharedKey)
	file, err := os.Open("big.txt")
	if err != nil {
		t.Error(err)
	}

	for {
		buffer := make([]byte, spice.ChunkSize+spice.Overhead)
		bytesRead, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		encryptor.Encrypt(buffer[:bytesRead])
	}

}
