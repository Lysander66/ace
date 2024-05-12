package cryptogo

import (
	"fmt"
	"testing"
)

func TestAESEncryptCBC(t *testing.T) {
	var (
		message = "hello world"
		key     = "57A891D97E332A9D"
		iv      = "3f56d1d36af225bf"
	)
	encrypted := AESEncryptCBC(message, key, iv)
	fmt.Println("加密数据:", encrypted)

	decrypted := AESDecryptCBC(encrypted, key, iv)
	fmt.Println("解密数据:", decrypted)
}
