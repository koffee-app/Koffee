package config

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"koffee/pkg/logger"
	"os"
)

var pasetoPublic ed25519.PublicKey
var pasetoPrivate ed25519.PrivateKey

var pasetoPublicStreaming ed25519.PublicKey

// PasetoInit inits the paseto configuration for every key in the system
func PasetoInit() {
	key, ok := os.LookupEnv("PASETO_PUBLIC_KEY")
	key2, ok2 := os.LookupEnv("PASETO_PRIVATE_KEY")

	if !ok || !ok2 {
		panic("Paseto Key Public not retrieved correctly...")
	}

	fmt.Println(key, key2)

	b, _ := hex.DecodeString(key)
	pasetoPublic = ed25519.PublicKey(b)
	b, _ = hex.DecodeString(key2)
	pasetoPrivate = ed25519.PrivateKey(b)

	// Get the PKCS8 Key
	f, _ := os.Open("key.txt")
	b = make([]byte, 128)
	_, err := f.Read(b)
	if err != nil {
		logger.Log("paseto", "Error retrieving key.txt")
		panic(err)
	}

	// Parse it
	bytesPriv, err := x509.ParsePKCS8PrivateKey(b)
	if err != nil {
		panic(err)
	}
	// Get the Private Key as ed25519 (For paseto)
	keyPriv := ed25519.PrivateKey(bytesPriv.(ed25519.PrivateKey))
	// Parse it as a public key
	pub := keyPriv.Public().(ed25519.PublicKey)
	// Finally in the data type that we want
	pasetoPublicStreaming = ed25519.PublicKey(pub)
}

// PrivateKeyParsed  Returns the parsed key for use in paseto
func PrivateKeyParsed() ed25519.PrivateKey {
	return pasetoPrivate
}

// PublicKeyParsed Returns the parsed key for use in paseto
func PublicKeyParsed() ed25519.PublicKey {
	return pasetoPublic
}

// StreamingServiceKey Returns the parsed ParsePKCS8PrivateKey as a public key for use in paseto to parse tokens from the streaming Rust service
func StreamingServiceKey() ed25519.PublicKey {
	return pasetoPublicStreaming
}
