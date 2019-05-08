package spice_test

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/tensortask/spice"
)

func TestDecryption(t *testing.T) {
	testString := []byte("this is a test")
	encryptedData := []byte{193, 169, 253, 213, 94, 44, 240, 65, 32, 194, 251, 86, 65, 75, 104, 6, 205, 132, 60, 67, 10, 68, 228, 244, 232, 29, 54, 141, 97, 141}
	nonce := [24]byte{58, 230, 79, 44, 187, 45, 107, 226, 245, 53, 169, 118, 218, 116, 235, 95, 132, 127, 166, 200, 203, 141, 251, 51}
	sharedKey := &[32]byte{199, 156, 103, 110, 157, 5, 107, 139, 94, 138, 53, 214, 74, 100, 211, 97, 106, 48, 11, 179, 200, 19, 244, 108, 138, 167, 49, 163, 156, 176, 66, 64}

	var resultBuffer bytes.Buffer
	resultWriter := bufio.NewWriter(&resultBuffer)

	decryptor := spice.NewDecryptor(resultWriter, nonce, sharedKey)
	err := decryptor.Decrypt(encryptedData)
	if err != nil {
		t.Error(err)
	}
	resultWriter.Flush()

	if !bytes.Equal(testString, resultBuffer.Bytes()) {
		t.Errorf("want: %v \n got: %v", testString, resultBuffer.Bytes())
	}
}
