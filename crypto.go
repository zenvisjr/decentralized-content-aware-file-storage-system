package main

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
)



// generateID generates a random ID
func generateID() string {
	buf := make([]byte, 32)
	io.ReadFull(rand.Reader, buf)
	return hex.EncodeToString(buf)
}

// hashKey hashes the given key using SHA-1
func hashKey(key string) string {
	encKey := sha1.Sum([]byte(key))
	return hex.EncodeToString(encKey[:])
}

// newEncryptionKey generates a random encryption key
func newEncryptionKey() []byte {
	enckeyBuf := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, enckeyBuf); err != nil {
		log.Panic(err)
	}
	return enckeyBuf
}

// copyStream copies data from src to dst using the given stream
func copyStream(stream cipher.Stream, blockSize int, src io.Reader, dst io.Writer) (int, error) {

	var (
		buf         = make([]byte, 32*1024)
		encyDataLen = blockSize
	)

	// Reads a chunk from src
	// Encrypts the chunk in-place using stream.XORKeyStream(...)
	// Writes the encrypted chunk to dist
	// Continues until the end of the input stream

	for {
		n, err := src.Read(buf)
		if n > 0 {
			// Encrypt in-place using CTR
			stream.XORKeyStream(buf, buf[:n])
			//Writes the encrypted chunk to dist
			nn, err := dst.Write(buf[:n])
			if err != nil {
				return 0, err
			}
			encyDataLen += nn
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}
	}
	return encyDataLen, nil
}

// copyEncrypt encrypts data from an input stream (src) and writes the encrypted result to an output stream (dist),
// using AES encryption in CTR mode (stream cipher style).
func copyEncrypt(dist io.Writer, src io.Reader, key []byte) (int, error) {

	// Initializes an AES cipher block using the given key.
	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, err
	}

	// Allocates an IV (Initialization Vector) of the block size
	//Purpose: Ensures uniqueness of encryption per session even with the same key.
	iv := make([]byte, block.BlockSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return 0, err
	}

	//Writes the IV to the output stream in start
	// The receiver must read and use this IV to decrypt.
	if _, err := dist.Write(iv); err != nil {
		return 0, err
	}

	// Creates a CTR mode stream cipher using the AES block and IV.
	stream := cipher.NewCTR(block, iv)
	return copyStream(stream, block.BlockSize(), src, dist)
}

// copyDecrypt decrypts data from src to dst using the given key
func copyDecrypt(dist io.Writer, src io.Reader, key []byte) (int, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, err
	}

	//read the iv from the src
	iv := make([]byte, block.BlockSize())
	if _, err := src.Read(iv); err != nil {
		return 0, err
	}

	// fmt.Println("block size", block.BlockSize())

	stream := cipher.NewCTR(block, iv)

	return copyStream(stream, block.BlockSize(), src, dist)
}

//Encrypted Size = Original Size + IV (Initilization Vector) Size

func hashFileContent(r io.Reader) (string, int64, error) {
	hasher := sha256.New()
	n, err := io.Copy(hasher, r)
	if err != nil {
		return "", 0, err
	}
	hashedFile := hasher.Sum(nil)
	fmt.Println("hashed file", hashedFile)
	return hex.EncodeToString(hashedFile), n, nil
}



func signSignature(data []byte, priv *rsa.PrivateKey) ([]byte, error) {
	hash := sha256.Sum256(data)
	return rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hash[:])
}

func verifySignature(data, sig []byte, pub *rsa.PublicKey) error {
	hash := sha256.Sum256(data)
	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, hash[:], sig)
}
