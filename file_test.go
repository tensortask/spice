package spice_test

import (
	"os"
	"testing"

	"github.com/tensortask/spice"
	"github.com/udhos/equalfile"
)

func TestFileEncryptionDecryption(t *testing.T) {
	hostKeys, err := spice.RandomKeyPair()
	if err != nil {
		t.Error(err)
	}
	peerKeys, err := spice.RandomKeyPair()
	if err != nil {
		t.Error(err)
	}
	sharedKey := hostKeys.SharedKey(peerKeys)
	nonce, err := spice.RandomNonce()
	if err != nil {
		t.Error(err)
	}

	encrypted_file, err := os.Create("testing/encrypted_shakespeare.txt")
	if err != nil {
		t.Error(err)
	}

	encryptor := spice.NewEncryptor(encrypted_file, nonce, sharedKey)
	err = encryptor.EncryptFile("testing/a_midsummers_night's_dream.txt")
	if err != nil {
		t.Error(err)
	}
	encrypted_file.Close()

	decrypted_file, err := os.Create("testing/decrypted_shakespeare.txt")
	if err != nil {
		t.Error(err)
	}

	decryptor := spice.NewDecryptor(decrypted_file, nonce, sharedKey)

	err = decryptor.DecryptFile("testing/encrypted_shakespeare.txt")
	if err != nil {
		t.Error(err)
	}
	decrypted_file.Close()
	cmp := equalfile.New(nil, equalfile.Options{})
	equal, err := cmp.CompareFile("testing/a_midsummers_night's_dream.txt", "testing/decrypted_shakespeare.txt")
	if err != nil {
		t.Error(err)
	}
	if !equal {
		t.Errorf("src file (testing/a_midsummers_night's_dream.txt) and decrypted file (testing/decrypted_shakespeare.txt) do not match.")
	}
}
