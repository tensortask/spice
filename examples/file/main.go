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
