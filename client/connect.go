package main

import (
	"fmt"
	"github.com/setavenger/p2p-chat/common"
)

func SendMessage(client Client, privateKey, publicKey, recipientPublicKey, plainMessage, parentID string) error {
	generatedMessage, err := GenerateMessage(privateKey, publicKey, recipientPublicKey, plainMessage)
	if err != nil {
		return err
	}

	generatedMessage.ParentID = parentID
	err = client.ForwardMessage(generatedMessage)
	if err != nil {
		return err
	}
	return nil
}

func SendMessageByUsername(client Client, privateKey, publicKey, username, plainMessage, parentID string) error {
	userWellKnown, err := common.GetPublicKeyForUsername(username)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = SendMessage(client, privateKey, publicKey, userWellKnown.PublicKey, plainMessage, parentID)
	if err != nil {
		return err
	}

	return nil
}
