package main

import (
	"fmt"
	"testing"
)

func TestEncryptAndDecryptAES(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef")
	plaintext := "hello, world!"
	ciphertext, err := EncryptAES(key, plaintext)
	if err != nil {
		t.Errorf("Error encrypting: %s", err.Error())
		return
	}
	decrypted, err := DecryptAES(key, ciphertext)
	if err != nil {
		t.Errorf("Error decrypting: %s", err.Error())
		return
	}
	if decrypted != plaintext {
		t.Errorf("Decrypted message does not match original message: %s != %s", decrypted, plaintext)
	}
}

func TestFullFlow(t *testing.T) {
	recipientPrivateKey, recipientPublicKey, err := GenerateKeyPair()
	fmt.Println(recipientPrivateKey)
	fmt.Println(recipientPublicKey)

	if err != nil {
		t.Errorf("Error generating key pair: %s", err.Error())
		return
	}
	senderPrivateKey, senderPublicKey, err := GenerateKeyPair()
	if err != nil {
		t.Errorf("Error generating key pair: %s", err.Error())
		return
	}

	plainMessage := "hello world!"

	generatedMessage, err := GenerateMessage(senderPrivateKey, senderPublicKey, recipientPublicKey, plainMessage)
	if err != nil {
		return
	}
	decryptMessage, err := DecryptAndVerifyMessage(recipientPrivateKey, generatedMessage)
	if err != nil {
		return
	}

	if decryptMessage != plainMessage {
		t.Errorf("Decrypted message does not match original message: %s != %s", decryptMessage, plainMessage)
	}
}
