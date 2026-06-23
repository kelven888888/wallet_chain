package aes

import (
	"fmt"
	"testing"
)

const (
	key = "IgkibX71IEf382PT"
	iv  = "IgkibX71IEf382PT"
)

func TestEncrypt(t *testing.T) {
	t.Log(New(key, iv).Encrypt("123456"))
}

func TestDecrypt(t *testing.T) {
	t.Log(New(key, iv).Decrypt("GO-ri84zevE-z1biJwfQPw=="))
}

func BenchmarkEncryptAndDecrypt(b *testing.B) {
	b.ResetTimer()
	aes := New(key, iv)
	for i := 0; i < b.N; i++ {
		encryptString, _ := aes.Encrypt("123456")
		fmt.Println(encryptString)
		aes.Decrypt(encryptString)
	}
}
