package main

import "bytes"
import "testing"

func TestMakeRecordKey(t *testing.T) {
	key := makeRecordKey("test")
	key_len := len(key)

	if key_len != 32 {
		t.Errorf("Generated key was %d in length, instead of expected 32", key_len)
	}

	if bytes.Equal(key, make([]byte, 32)) {
		t.Error("Generated key was empty")
	}
}

func TestMakeNonce(t *testing.T) {
	nonce := makeNonce()
	nonce_len := len(nonce)

	if nonce_len != 12 {
		t.Errorf("Generated nonce was %d in length, instead of expected 12", nonce_len)
	}

	if bytes.Equal(nonce, make([]byte, 12)) {
		t.Error("Generated nonce was empty")
	}
}

func TestEncrypt(t *testing.T) {
	e := encrypt("plaintext", "client-secret")
	e_len := len(e)

	if e_len != 41 {
		t.Errorf("Expected envelope to be 25 bytes long, got %d", e_len)
	}
}

func TestDecrypt(t *testing.T) {
	e := []byte{91, 0, 0, 1, 48, 187, 74, 223, 50, 134, 211, 188, 64, 223, 233, 43, 128, 244, 163, 25, 107, 254, 180, 22, 230, 250, 202, 38, 168, 63, 207, 198, 114, 117, 97, 166, 21, 76, 151, 12, 6}

	_, pt := decrypt(string(e), "client-secret")
	if pt != "plaintext" {
		t.Error("Decrypt with valid keys failed")
	}

	err, _ := decrypt(string(e), "wrong-secret")
	if err == nil {
		t.Error("Decrypt with wrong client-secret did not return error")
	}

	server_secret = []byte("wrong-secret")
	err, _ = decrypt(string(e), "client-secret")
	if err == nil {
		t.Error("Decrypt with wrong server-secret did not return error")
	}
}
