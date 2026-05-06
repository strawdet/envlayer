package encryptor_test

import (
	"fmt"
	"testing"

	"github.com/user/envlayer/internal/encryptor"
)

func BenchmarkEncrypt(b *testing.B) {
	e := encryptor.New("bench-passphrase")
	plaintext := "some-sensitive-value-123"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := e.Encrypt(plaintext); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecrypt(b *testing.B) {
	e := encryptor.New("bench-passphrase")
	enc, _ := e.Encrypt("some-sensitive-value-123")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := e.Decrypt(enc); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncryptMap(b *testing.B) {
	e := encryptor.New("bench-passphrase")
	env := make(map[string]string, 20)
	for i := 0; i < 20; i++ {
		env[fmt.Sprintf("KEY_%d", i)] = fmt.Sprintf("value-%d", i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := e.EncryptMap(env); err != nil {
			b.Fatal(err)
		}
	}
}
