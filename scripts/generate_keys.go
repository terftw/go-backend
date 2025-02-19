package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"log"
)

func main() {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}

	// Convert private key to PEM format
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	// Convert public key to PEM format
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
	})

	// Encode as base64 for .env
	privateKeyBase64 := base64.StdEncoding.EncodeToString(privateKeyPEM)
	publicKeyBase64 := base64.StdEncoding.EncodeToString(publicKeyPEM)

	// Print keys for .env
	println("PRIVATE_KEY=" + privateKeyBase64)
	println("PUBLIC_KEY=" + publicKeyBase64)
}
