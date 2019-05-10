package spice

import (
	"io"
)

// EncryptFromReader ....
func (e *Encryptor) EncryptFromReader(reader io.Reader) error {
	for {
		buffer := make([]byte, DefaultChunkSize)
		bytesRead, err := reader.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		e.Encrypt(buffer[:bytesRead])
	}
	return nil
}

// DecryptFromReader ...
func (d *Decryptor) DecryptFromReader(reader io.Reader) error {
	for {
		buffer := make([]byte, DefaultChunkSize+DefaultOverhead)
		bytesRead, err := reader.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		d.Decrypt(buffer[:bytesRead])
	}
	return nil
}
