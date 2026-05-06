package encryptor_test

import (
	"strings"
	"testing"

	"github.com/user/envlayer/internal/encryptor"
)

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	e := encryptor.New("supersecret")
	plaintext := "my-database-password"

	enc, err := e.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt error: %v", err)
	}
	if enc == plaintext {
		t.Error("encrypted value should differ from plaintext")
	}

	dec, err := e.Decrypt(enc)
	if err != nil {
		t.Fatalf("Decrypt error: %v", err)
	}
	if dec != plaintext {
		t.Errorf("expected %q, got %q", plaintext, dec)
	}
}

func TestEncrypt_NonDeterministic(t *testing.T) {
	e := encryptor.New("key")
	a, _ := e.Encrypt("value")
	b, _ := e.Encrypt("value")
	if a == b {
		t.Error("two encryptions of the same value should differ (random nonce)")
	}
}

func TestDecrypt_InvalidBase64(t *testing.T) {
	e := encryptor.New("key")
	_, err := e.Decrypt("!!!not-base64!!!")
	if err == nil {
		t.Error("expected error for invalid base64 input")
	}
}

func TestDecrypt_TamperedCiphertext(t *testing.T) {
	e := encryptor.New("key")
	enc, _ := e.Encrypt("secret")
	tampered := strings.ToUpper(enc)
	_, err := e.Decrypt(tampered)
	if err == nil {
		t.Error("expected error for tampered ciphertext")
	}
}

func TestDecrypt_WrongKey(t *testing.T) {
	e1 := encryptor.New("correct-key")
	e2 := encryptor.New("wrong-key")
	enc, _ := e1.Encrypt("secret")
	_, err := e2.Decrypt(enc)
	if err == nil {
		t.Error("expected error when decrypting with wrong key")
	}
}

func TestEncryptMap_DecryptMap(t *testing.T) {
	e := encryptor.New("mapkey")
	original := map[string]string{
		"DB_PASS": "hunter2",
		"API_KEY": "abc123",
	}
	encrypted, err := e.EncryptMap(original)
	if err != nil {
		t.Fatalf("EncryptMap error: %v", err)
	}
	decrypted, err := e.DecryptMap(encrypted)
	if err != nil {
		t.Fatalf("DecryptMap error: %v", err)
	}
	for k, v := range original {
		if decrypted[k] != v {
			t.Errorf("key %s: expected %q, got %q", k, v, decrypted[k])
		}
	}
}
