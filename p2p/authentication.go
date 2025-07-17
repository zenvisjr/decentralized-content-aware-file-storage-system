package p2p

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func EnsureKeyPair() error {
	if _, err := os.Stat("private.pem"); os.IsNotExist(err) {
		// fmt.Println("No key pair found, generating...")
		return generateKeyPair()
	}
	// fmt.Println("Key pair already exists.")
	return nil
}

func generateKeyPair() error {
	// fmt.Println("Generating key pair...")
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	privFile, err := os.Create("private.pem")
	if err != nil {
		return err
	}

	defer privFile.Close()
	if err := pem.Encode(privFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}); err != nil {
		return err
	}

	pubFile, err := os.Create("public.pem")
	if err != nil {
		return err
	}
	defer pubFile.Close()
	pubBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}
	if err := pem.Encode(pubFile, &pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes}); err != nil {
		return err
	}

	return nil
}

// func LoadPublicKey() ([]byte, error) {
// 	publicKey, err := os.ReadFile("public.pem")
// 	if err != nil {
// 		return nil, err
// 	}
// 	return publicKey, nil

// }

func LoadPublicKey() (*rsa.PublicKey, error) {
	// Read the PEM file
	pubData, err := os.ReadFile("public.pem")
	if err != nil {
		return nil, fmt.Errorf("failed to read public.pem: %w", err)
	}

	// Decode PEM block
	block, _ := pem.Decode(pubData)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("invalid PEM block or wrong type")
	}

	// Parse the public key bytes
	pubAny, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	// Type assert to *rsa.PublicKey
	pubKey, ok := pubAny.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}

	return pubKey, nil
}

func LoadPrivateKey() (*rsa.PrivateKey, error) {
	privData, err := os.ReadFile("private.pem")
	if err != nil {
		return nil, fmt.Errorf("failed to read private.pem: %w", err)
	}
	block, _ := pem.Decode(privData)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("invalid private key PEM block")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

