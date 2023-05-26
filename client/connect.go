package main

import (
	"fmt"
	"p2p-chat/common"
)

func SendMessage(client Client, privateKey, publicKey, recipientPublicKey, plainMessage string) error {
	generatedMessage, err := GenerateMessage(privateKey, publicKey, recipientPublicKey, plainMessage)
	if err != nil {
		return err
	}

	err = client.ForwardMessage(generatedMessage)
	if err != nil {
		return err
	}
	return nil
}

func SendMessageByUsername(client Client, privateKey, publicKey, username, plainMessage string) error {
	userWellKnown, err := common.GetPublicKeyForUsername(username)
	if err != nil {
		fmt.Println(err)
		return err
	}

	generatedMessage, err := GenerateMessage(privateKey, publicKey, userWellKnown.PublicKey, plainMessage)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = client.ForwardMessage(generatedMessage)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
