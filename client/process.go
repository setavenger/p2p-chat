package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/setavenger/p2p-chat/common"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
)

// EncryptAES - encrypt plaintext using AES encryption with the generated key
func EncryptAES(key []byte, plaintext string) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err = rand.Read(iv); err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))
	return ciphertext, nil
}

// DecryptAES - decrypt ciphertext using the decrypted AES key
func DecryptAES(key []byte, ciphertext []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext), nil
}

// EncryptRSA - encrypt plaintext using the recipient's public RSA key
func EncryptRSA(pubkey *rsa.PublicKey, plaintext []byte) ([]byte, error) {
	hashed := sha256.Sum256(plaintext)
	ciphertext, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		pubkey,
		plaintext,
		hashed[:],
	)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

// DecryptRSA - decrypt ciphertext using the client's private RSA key
func DecryptRSA(encrypted []byte, privateKey *rsa.PrivateKey) (string, error) {
	// Decrypt the message using RSA decryption with PKCS1v15 padding
	decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encrypted)
	if err != nil {
		return "", err
	}

	// Extract the AES key and IV from the decrypted message
	keySize := aes.BlockSize
	if len(decrypted) < keySize {
		return "", errors.New("invalid decrypted message")
	}
	key := decrypted[:keySize]
	iv := decrypted[keySize : 2*keySize]

	// Decrypt the message using AES-CBC decryption with the extracted key and IV
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	padded := decrypted[2*keySize:]
	mode.CryptBlocks(padded, padded)

	// Unpad the decrypted message using PKCS7 padding
	unpadded, err := pkcs7Unpad(padded)
	if err != nil {
		return "", err
	}

	return string(unpadded), nil
}

func pkcs7Unpad(padded []byte) ([]byte, error) {
	length := len(padded)
	unpadding := int(padded[length-1])
	if unpadding > length {
		return nil, errors.New("invalid padding")
	}
	return padded[:length-unpadding], nil
}

// EncryptMessage Encrypt a message using RSA and AES
func EncryptMessage(privateKey string, receiverKey string, message []byte) (string, error) {
	sharedSecret, err := ComputeSharedSecret(privateKey, receiverKey)
	if err != nil {
		log.Printf("Error computing shared key: %s. \n", err.Error())
		return "", err
	}

	encryptedMessage, err := Encrypt(message, sharedSecret)
	if err != nil {
		log.Printf("Error encrypting message: %s. \n", err.Error())
		return "", err
	}
	return encryptedMessage, nil
}

// DecryptMessage Decrypt a message using RSA and AES
func DecryptMessage(privateKey, senderPublicKey string, encryptedMessage string) (string, error) {
	sharedSecret, err := ComputeSharedSecret(privateKey, senderPublicKey)
	if err != nil {
		log.Printf("Error computing shared key: %s. \n", err.Error())
		return "", err
	}
	txt, err := Decrypt(encryptedMessage, sharedSecret)
	if err != nil {
		log.Printf("Error decrypting message: %s. \n", err.Error())
		return "", err
	}
	//fmt.Print(txt)

	return txt, nil
}

// FetchMessages - Supposed to get messages from server
func FetchMessages(publicKey string, privateKey *rsa.PrivateKey) ([]common.Message, error) {
	// Fetch the unread messages from the server
	resp, err := http.Get("http://example.com/fetch-messages/" + publicKey)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Decode the response body from base64
	decoded, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		return nil, err
	}

	// Decrypt the response body using RSA and AES decryption
	decrypted, err := DecryptRSA(decoded, privateKey)
	if err != nil {
		return nil, err
	}

	// Unmarshal the decrypted response body into a slice of messages
	var messages []common.Message
	err = json.Unmarshal([]byte(decrypted), &messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

// GenerateMessage - generates a new encrypted message and signature from the sender to the recipient
func GenerateMessage(senderPrivateKey, senderPublicKey, recipientPublicKey, message string) (*common.Message, error) {
	// Encrypt the message using AES and generate a signature using SHA256
	ciphertext, err := EncryptMessage(senderPrivateKey, recipientPublicKey, []byte(message))
	if err != nil {
		return nil, err
	}
	signature, err := common.Sign(senderPrivateKey, []byte(message))
	if err != nil {
		return nil, err
	}

	msg := common.Message{
		Sender:    senderPublicKey,
		Recipient: recipientPublicKey,
		Encrypted: ciphertext,
		Timestamp: uint64(time.Now().Unix()),
		Signature: hex.EncodeToString(signature),
	}

	err = msg.GetId()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Return the message object with the encrypted AES key
	return &msg, nil
}

// DecryptAndVerifyMessage - verifies the signature and decrypts the message from the sender
func DecryptAndVerifyMessage(ownPrivateKey string, msg *common.Message) (string, error) {

	ownPublicKey, err := common.GetPublicKey(ownPrivateKey)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	// todo both directions sent/received
	// adding 02 to signal that this is a compressed public key (33 bytes)
	pubKeyBytes, err := hex.DecodeString("02" + msg.Sender)
	if err != nil {
		return "", fmt.Errorf("Error decoding hex string of receiver public key: %s. \n", err)
	}
	pubKey, err := btcec.ParsePubKey(pubKeyBytes)
	if err != nil {
		return "", fmt.Errorf("Error parsing receiver public key: %s. \n", err)
	}

	var decryptMessage string
	if ownPublicKey == msg.Sender {
		// Decode the encrypted message and signature
		decryptMessage, err = DecryptMessage(ownPrivateKey, msg.Recipient, msg.Encrypted)
		if err != nil {
			return "", err
		}
	} else {
		// Decode the encrypted message and signature
		decryptMessage, err = DecryptMessage(ownPrivateKey, msg.Sender, msg.Encrypted)
		if err != nil {
			return "", err
		}
	}

	//fmt.Println(decryptMessage)
	// Verify the signature
	// read signature
	s, err := hex.DecodeString(msg.Signature)
	if err != nil {
		return "", fmt.Errorf("signature '%s' is invalid hex: %w", s, err)
	}
	sig, err := schnorr.ParseSignature(s)
	if err != nil {
		return "", fmt.Errorf("failed to parse signature: %w", err)
	}

	// check signature
	hash := sha256.Sum256([]byte(decryptMessage))
	isValid := sig.Verify(hash[:], pubKey)

	if !isValid {
		return "", fmt.Errorf("signature does not check out")
	}

	return decryptMessage, nil
}

func DecryptAllMessages(privateKey string, messages []common.Message) []common.MessagePlain {
	var plains []common.MessagePlain
	for _, message := range messages {
		plainText, err := DecryptAndVerifyMessage(privateKey, &message)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		plains = append(plains, common.MessagePlain{
			ID:             message.ID,
			Sender:         message.Sender,
			SenderUsername: message.SenderUsername,
			Content:        plainText,
			Recipient:      message.Recipient,
			Timestamp:      message.Timestamp,
			Read:           message.Read,
			ParentID:       message.ParentID,
		})
	}

	return plains
}
