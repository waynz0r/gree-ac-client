package client

import (
	"crypto/aes"

	"github.com/andreburgaud/crypt2go/ecb"
	"github.com/andreburgaud/crypt2go/padding"
)

func encrypt(plaintext, key []byte) (encypted []byte, err error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	mode := ecb.NewECBEncrypter(block)
	padder := padding.NewPkcs7Padding(mode.BlockSize())
	plaintext, err = padder.Pad(plaintext) // padd last block of plaintext if block size less than block cipher size
	if err != nil {
		return
	}

	encypted = make([]byte, len(plaintext))
	mode.CryptBlocks(encypted, plaintext)

	return
}

func decrypt(encrypted, key []byte) (plaintext []byte, err error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	mode := ecb.NewECBDecrypter(block)
	plaintext = make([]byte, len(encrypted))
	mode.CryptBlocks(plaintext, encrypted)
	padder := padding.NewPkcs7Padding(mode.BlockSize())
	plaintext, err = padder.Unpad(plaintext) // unpad plaintext after decryption

	return
}
