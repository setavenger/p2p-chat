package common

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"log"
)

// GenerateKeyPair - generate RSA public and private keys for the client
func GenerateKeyPair() (string, string, error) {
	privateKey := KeyGen()
	publicKey, err := GetPublicKey(privateKey)
	if err != nil {
		return "", "", err
	}
	return privateKey, publicKey, nil
}

func GetPublicKey(sk string) (string, error) {
	b, err := hex.DecodeString(sk)
	if err != nil {
		return "", err
	}

	_, pk := btcec.PrivKeyFromBytes(b)
	return hex.EncodeToString(schnorr.SerializePubKey(pk)), nil
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
