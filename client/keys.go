package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"log"
	"strings"

	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"io"
	"math/big"
	mrand "math/rand"

	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

func GeneratePrivateKey() string {
	params := btcec.S256().Params()
	one := new(big.Int).SetInt64(1)

	b := make([]byte, params.BitSize/8+8)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return ""
	}

	k := new(big.Int).SetBytes(b)
	n := new(big.Int).Sub(params.N, one)
	k.Mod(k, n)
	k.Add(k, one)

	return hex.EncodeToString(k.Bytes())
}

func GetPublicKey(sk string) (string, error) {
	b, err := hex.DecodeString(sk)
	if err != nil {
		return "", err
	}

	_, pk := btcec.PrivKeyFromBytes(b)
	return hex.EncodeToString(schnorr.SerializePubKey(pk)), nil
}

// ComputeSharedSecret - ECDH
func ComputeSharedSecret(senderPrivKey string, receiverPubKey string) (sharedSecret []byte, err error) {
	privKeyBytes, err := hex.DecodeString(senderPrivKey)
	if err != nil {
		return nil, fmt.Errorf("Error decoding sender private key: %s. \n", err)
	}
	privKey, _ := btcec.PrivKeyFromBytes(privKeyBytes)

	// adding 02 to signal that this is a compressed public key (33 bytes)
	pubKeyBytes, err := hex.DecodeString("02" + receiverPubKey)
	if err != nil {
		return nil, fmt.Errorf("Error decoding hex string of receiver public key: %s. \n", err)
	}
	pubKey, err := btcec.ParsePubKey(pubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("Error parsing receiver public key: %s. \n", err)
	}

	return btcec.GenerateSharedSecret(privKey, pubKey), nil
}

// Encrypt - aes-256-cbc
func Encrypt(message []byte, key []byte) (string, error) {
	// block size is 16 bytes
	iv := make([]byte, 16)
	// can probably use a less expensive lib since IV has to only be unique; not perfectly random; math/rand?
	_, err := rand.Read(iv)
	if err != nil {
		return "", fmt.Errorf("Error creating initization vector: %s. \n", err.Error())
	}

	// automatically picks aes-256 based on key length (32 bytes)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("Error creating block cipher: %s. \n", err.Error())
	}
	mode := cipher.NewCBCEncrypter(block, iv)

	plaintext := message

	// add padding
	base := len(plaintext)
	extra := 0
	var padding int
	if base < 100 {
		// add some random padding to this message since it is too small
		extra = mrand.Intn(230) // the total padding will be padding + extra, which can't be more than 256
	}

	// this will be a number between 1 and 16 (including), never 0
	padding = block.BlockSize() - (base+extra)%block.BlockSize()

	// encode the padding in all the padding bytes themselves
	padtext := bytes.Repeat([]byte{byte(padding + extra)}, padding+extra)

	paddedMsgBytes := append(plaintext, padtext...)

	ciphertext := make([]byte, len(paddedMsgBytes))
	mode.CryptBlocks(ciphertext, paddedMsgBytes)

	return base64.StdEncoding.EncodeToString(ciphertext) + "?iv=" + base64.StdEncoding.EncodeToString(iv), nil
}

// Decrypt - aes-256-cbc
func Decrypt(content string, key []byte) (string, error) {
	parts := strings.Split(content, "?iv=")
	if len(parts) < 2 {
		return "", fmt.Errorf("Error parsing encrypted message: no initilization vector. \n")
	}

	ciphertext, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return "", fmt.Errorf("Error decoding ciphertext from base64: %s. \n", err)
	}

	iv, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("Error decoding iv from base64: %s. \n", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("Error creating block cipher: %s. \n", err.Error())
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// remove padding
	padding := int(plaintext[len(plaintext)-1]) // the padding amount is encoded in the padding bytes themselves
	message := string(plaintext[0 : len(plaintext)-padding])

	return message, nil
}

func GenerateSeedWords() (string, error) {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return "", err
	}

	words, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}

	return words, nil
}

func SeedFromWords(words string) []byte {
	return bip39.NewSeed(words, "")
}

func PrivateKeyFromSeed(seed []byte) (string, error) {
	key, err := bip32.NewMasterKey(seed)
	if err != nil {
		return "", err
	}

	derivationPath := []uint32{
		bip32.FirstHardenedChild + 44,
		bip32.FirstHardenedChild + 1237,
		bip32.FirstHardenedChild + 0,
		0,
		0,
	}

	next := key
	for _, idx := range derivationPath {
		var err error
		next, err = next.NewChildKey(idx)
		if err != nil {
			return "", err
		}
	}

	return hex.EncodeToString(next.Key), nil
}

func ValidateWords(words string) bool {
	return bip39.IsMnemonicValid(words)
}

func KeyGen() string {
	seedWords, err := GenerateSeedWords()
	if err != nil {
		log.Println(err)
		return ""
	}

	seed := SeedFromWords(seedWords)

	sk, err := PrivateKeyFromSeed(seed)
	if err != nil {
		log.Println(err)
		return ""
	}

	return sk
}

func GetPubKey(privateKey string) string {
	if keyb, err := hex.DecodeString(privateKey); err != nil {
		log.Printf("Error decoding key from hex: %s\n", err.Error())
		return ""
	} else {
		_, pubkey := btcec.PrivKeyFromBytes(keyb)
		return hex.EncodeToString(schnorr.SerializePubKey(pubkey))
	}
}
