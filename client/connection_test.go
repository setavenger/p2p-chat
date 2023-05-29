package main

import (
	"fmt"
	"github.com/setavenger/p2p-chat/common"
	"testing"
)

const examplePrivateKey1 = "c687bf0179f2a038c04fa1d30abf23c7797c0bd65316ca2fee021a1f9241ad52"
const examplePublicKey1 = "0f8abf593879bca06435993639130ea2caf017196b3379dc208e003b9b411e48"

const examplePrivateKey2 = "b5ecbb76d605b0d9025bf7cdd830bf9c01a0a1967d89462aa4016d7fe897f63e"
const examplePublicKey2 = "0fbbead7194be93d6fc659dd0fe9bdf8b459febfbbfa602322d081eef19688c2"

const examplePrivateKey3 = "27f414f1836bf6a2df7056d3881f2bce74848aea31ce18facb0b02e7d107aa99"
const examplePublicKey3 = "d623bd7bae8a256c35ba9db430795c501532b75a4dc8aa4e1ef36fb00b37da3a"

const baseURL = "http://localhost:8889"

// TestSendMessage - requires that server is running on Port 8000 locally
func TestSendMessage(t *testing.T) {
	privateKey, publicKey, err := common.GenerateKeyPair()
	if err != nil {
		t.Errorf("Error generating key pair: %s", err.Error())
		return
	}

	_, recipientPublicKey, err := common.GenerateKeyPair()
	if err != nil {
		t.Errorf("Error generating key pair: %s", err.Error())
		return
	}

	myMessage1 := "Hello, what is happening."
	err = SendMessage(Client{BaseURL: baseURL}, privateKey, publicKey, recipientPublicKey, myMessage1, "")
	if err != nil {
		t.Errorf("Error sending message: %s", err.Error())
		return
	}
	myMessage2 := "This should be working. Moreover this message is going to be way longer than the other message. Let's see how big the difference is." + "This should be working. Moreover this message is going to be way longer than the other message. Let's see how big the difference is." + "This should be working. Moreover this message is going to be way longer than the other message. Let's see how big the difference is."
	err = SendMessage(Client{BaseURL: baseURL}, privateKey, publicKey, recipientPublicKey, myMessage2, "")
	if err != nil {
		t.Errorf("Error sending message: %s", err.Error())
		return
	}
}

// TestSendMessageByUsername - requires that server is running on Port 8000 locally
func TestSendMessageByUsername(t *testing.T) {
	privateKey, publicKey, err := common.GenerateKeyPair()
	if err != nil {
		t.Errorf("Error generating key pair: %s", err.Error())
		return
	}

	username := "me@localhost:8889"

	message := "Hello, what is happening."
	err = SendMessageByUsername(Client{BaseURL: baseURL}, privateKey, publicKey, username, message, "")
	if err != nil {
		t.Errorf("Error sending message: %s", err.Error())
		return
	}
}

// TestSendMessage - requires that server is running on Port 8000 locally
func TestSendAndRetrieve(t *testing.T) {

	myMessage1 := "I love it so much."
	err := SendMessage(Client{BaseURL: baseURL}, examplePrivateKey1, examplePublicKey1, examplePublicKey2, myMessage1, "")
	if err != nil {
		t.Errorf("Error sending message: %s", err.Error())
		return
	}
	myMessage2 := "This should be working. Moreover this message is going to be way longer than the other message. Let's see how big the difference is."
	err = SendMessage(Client{BaseURL: baseURL}, examplePrivateKey3, examplePublicKey3, examplePublicKey2, myMessage2, "")
	if err != nil {
		t.Errorf("Error sending message: %s", err.Error())
		return
	}

	retrieveClient := Client{
		BaseURL:    baseURL,
		PrivateKey: examplePrivateKey2,
		PublicKey:  examplePublicKey2,
	}
	unreadMessages, err := retrieveClient.GetUnreadMessages()
	if err != nil {
		t.Errorf("Error receiving messages: %s", err.Error())
		return
	}

	for _, unreadMessage := range unreadMessages {
		var message string
		message, err = DecryptAndVerifyMessage(examplePrivateKey2, &unreadMessage)
		if err != nil {
			t.Errorf("Error decrypting and veryfying the message: %s", err.Error())
			return
		}
		fmt.Println(message)
	}
}

// TestSendMessage - requires that server is running on Port 8000 locally
func TestSendAndRetrieve2(t *testing.T) {

	sender1PrivateKey, sender1PublicKey, err := common.GenerateKeyPair()
	if err != nil {
		t.Errorf("Error generating key pair: %s", err.Error())
		return
	}
	sender2PrivateKey, sender2PublicKey, err := common.GenerateKeyPair()
	if err != nil {
		t.Errorf("Error generating key pair: %s", err.Error())
		return
	}

	recipientPrivateKey, recipientPublicKey, err := common.GenerateKeyPair()
	if err != nil {
		t.Errorf("Error generating key pair: %s", err.Error())
		return
	}

	myMessage1 := "Hello!"
	err = SendMessage(Client{BaseURL: baseURL}, sender1PrivateKey, sender1PublicKey, recipientPublicKey, myMessage1, "")
	if err != nil {
		t.Errorf("Error sending message: %s", err.Error())
		return
	}
	myMessage2 := "This should be working. Moreover this message is going to be way longer than the other message. Let's see how big the difference is."
	err = SendMessage(Client{BaseURL: baseURL}, sender2PrivateKey, sender2PublicKey, recipientPublicKey, myMessage2, "")
	if err != nil {
		t.Errorf("Error sending message: %s", err.Error())
		return
	}

	retrieveClient := Client{
		BaseURL:    baseURL,
		PrivateKey: recipientPrivateKey,
		PublicKey:  recipientPublicKey,
	}
	unreadMessages, err := retrieveClient.GetUnreadMessages()
	if err != nil {
		t.Errorf("Error receiving messages: %s", err.Error())
		return
	}

	if len(unreadMessages) != 2 {
		t.Errorf("not all messages were received")
	}
}

// TestSendMessage - requires that server is running on Port 8000 locally
func TestSendAndReply(t *testing.T) {

	myMessage1 := "Thanks for the reply."
	err := SendMessage(Client{BaseURL: baseURL}, examplePrivateKey1, examplePublicKey1, examplePublicKey2, myMessage1, "")
	if err != nil {
		t.Errorf("Error sending message: %s", err.Error())
		return
	}
	retrieveClient := Client{
		BaseURL:    baseURL,
		PrivateKey: examplePrivateKey2,
		PublicKey:  examplePublicKey2,
	}
	unreadMessages, err := retrieveClient.GetUnreadMessages()
	if err != nil {
		t.Errorf("Error receiving messages: %s", err.Error())
		return
	}
	myMessage2 := "Yeah, sure no worries."
	err = SendMessage(Client{BaseURL: baseURL}, examplePrivateKey2, examplePublicKey2, examplePublicKey1, myMessage2, unreadMessages[0].ID)
	if err != nil {
		t.Errorf("Error sending message: %s", err.Error())
		return
	}
}

func TestRetrieveMessages(t *testing.T) {
	retrieveClient := Client{
		BaseURL:    baseURL,
		PrivateKey: examplePrivateKey1,
		PublicKey:  examplePublicKey1,
	}
	unreadMessages, err := retrieveClient.GetUnreadMessages()
	if err != nil {
		t.Errorf("Error receiving messages: %s", err.Error())
		return
	}

	for _, unreadMessage := range unreadMessages {
		var message string
		message, err = DecryptAndVerifyMessage(examplePrivateKey1, &unreadMessage)
		if err != nil {
			t.Errorf("Error decrypting and veryfying the message: %s", err.Error())
			return
		}
		fmt.Println(message)
	}
}
