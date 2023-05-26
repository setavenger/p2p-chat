package common

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
)

func VerifySignature(message []byte, signature, publicKey string) error {
	pubKeyBytes, err := hex.DecodeString("02" + publicKey)
	if err != nil {
		return fmt.Errorf("Error decoding hex string of receiver public key: %s. \n", err)
	}
	pubKey, err := btcec.ParsePubKey(pubKeyBytes)
	if err != nil {
		return fmt.Errorf("Error parsing receiver public key: %s. \n", err)
	}

	s, err := hex.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("signature '%s' is invalid hex: %w", s, err)
	}
	sig, err := schnorr.ParseSignature(s)
	if err != nil {
		return fmt.Errorf("failed to parse signature: %w", err)
	}

	// check signature
	hash := sha256.Sum256(message)
	isValid := sig.Verify(hash[:], pubKey)

	if isValid {
		return nil
	} else {
		return errors.New("could not verify sender")
	}
}

func Sign(privateKey string, message []byte) ([]byte, error) {
	h := sha256.Sum256(message)

	s, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("sign called with invalid private key '%s': %w", privateKey, err)
	}
	sk, _ := btcec.PrivKeyFromBytes(s)

	sig, err := schnorr.Sign(sk, h[:])
	if err != nil {
		return nil, err
	}
	return sig.Serialize(), err
}
