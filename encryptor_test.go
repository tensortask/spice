package spice_test

import (
	"fmt"
	"io"
	"os"
	"testing"
	"time"

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
	info, err := file.Stat()
	if err != nil {
		// Could not obtain stat, handle error
	}
	fileSizeMB := float64(info.Size()) / float64(1000000)
	start := time.Now()

	for {
		buffer := make([]byte, spice.DefaultChunkSize+spice.DefaultOverhead)
		bytesRead, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		encryptor.Encrypt(buffer[:bytesRead])
	}
	elapsed := time.Since(start).Seconds()
	fmt.Printf("encrypted %f MB in %f seconds (%f MB/s)\n", fileSizeMB, elapsed, (fileSizeMB / elapsed))

}
