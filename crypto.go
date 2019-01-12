package main

import "crypto/aes"
import "crypto/cipher"
import "crypto/sha256"
import "crypto/rand"
import "fmt"
import "golang.org/x/crypto/hkdf"
import "io"
import "io/ioutil"

var server_secret []byte

const versionHeader = "\x5b\x00\x00\x01"

func init() {
	server_secret = readServerSecret()
}

func readServerSecret() []byte {
	secret, err := ioutil.ReadFile("keyfile")

	if err != nil {
		println("Could not read keyfile!")
		panic(err)
	}
	if len(secret) < 8 {
		panic("Key too short or empty")
	}

	return secret
}

// Generate a secure 32-byte key based on the server secret and the supplied
// client secret
func makeRecordKey(client_secret string) []byte {
	hkdf := hkdf.New(sha256.New, append(server_secret, []byte(client_secret)...), nil, nil)

	key := make([]byte, 32)
	if _, err := io.ReadFull(hkdf, key); err != nil {
		panic(err)
	}

	return key
}

// Initialise an AES-GCM cipher for a record using the supplied client_secret
func makeCipher(client_secret string) cipher.AEAD {
	block, err := aes.NewCipher(makeRecordKey(client_secret))
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	return aesgcm
}

func makeNonce() []byte {
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	return nonce
}

// Encrypt given plaintext using the supplied client_secret. The client secret
// is supplied by the client, which is responsible for generating a unique
// secret for each record.
// Returns an envelope that can be stored in the backend
func encrypt(plaintext, client_secret string) string {
	aesgcm := makeCipher(client_secret)
	nonce := makeNonce()
	ciphertext := aesgcm.Seal(nil, nonce, []byte(plaintext), nil)

	// envelope format:
	// start | len | content
	// 0     | 4   | versionHeader
	// 4     | 12  | nonce
	// 16    |     | ciphertext
	return fmt.Sprintf("%s%s%s", versionHeader, nonce, ciphertext)
}

// Deconstruct the envelope with the provided client_secret. Checks for
// envelope version compatibility and returns decrypted and authenticated
// plaintext
func decrypt(envelope, client_secret string) (error, string) {
	bytes := []byte(envelope)
	aesgcm := makeCipher(client_secret)

	if envelope[0:4] != versionHeader {
		panic("envelope version mismatch")
	}

	plaintext, err := aesgcm.Open(nil, bytes[4:16], bytes[16:], nil)
	if err != nil {
		return err, ""
	}

	return nil, string(plaintext)
}
