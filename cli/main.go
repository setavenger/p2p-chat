package main

import (
	"fmt"
	"github.com/setavenger/p2p-chat/common"
	"log"
)

func main() {
	privateKey, publicKey, err := common.GenerateKeyPair()
	if err != nil {
		log.Fatalf("Error generating keys: %v\n", err)
	}

	fmt.Printf("Private Key:	%s\n", privateKey)
	fmt.Printf("Public Key:	%s\n", publicKey)
}
