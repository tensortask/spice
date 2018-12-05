package spice_test

import (
	"crypto/rand"
	"io"
	"os"
	"testing"

	"github.com/tensortask/spice"
	"golang.org/x/crypto/nacl/box"
)

func TestPipeEncryption(t *testing.T) {

	var nonce [24]byte
	_, err := io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		t.Error(err)
	}

	senderPublicKey, senderPrivateKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		t.Error(err)
	}
	recipientPublicKey, recipientPrivateKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		t.Error(err)
	}
	infile, _ := os.Open("big.txt")
	outfile, _ := os.Create("big_copy.txt")

	reader, writer := io.Pipe()
	defer reader.Close()

	sendPipe := spice.NewEncyptedPipe(infile, writer, senderPublicKey, senderPrivateKey)
	recPipe := spice.NewEncyptedPipe(reader, outfile, recipientPublicKey, recipientPrivateKey)

	go func() {
		defer writer.Close()
		defer infile.Close()
		sendPipe.Encrypt(&nonce, recipientPublicKey)
	}()

	recPipe.Decrypt(&nonce, senderPublicKey)
}
