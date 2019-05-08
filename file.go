package spice

import (
	"io"
	"os"
)

// EncryptFile ....
func (e *Encryptor) EncryptFile(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	for {
		buffer := make([]byte, DefaultChunkSize)
		bytesRead, err := file.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		e.Encrypt(buffer[:bytesRead])
	}
	return nil
}

// DecryptFile ...
func (d *Decryptor) DecryptFile(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	for {
		buffer := make([]byte, DefaultChunkSize+DefaultOverhead)
		bytesRead, err := file.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		d.Decrypt(buffer[:bytesRead])
	}
	return nil
}
