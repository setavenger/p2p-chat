package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/setavenger/p2p-chat/common"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

// Client struct
type Client struct {
	BaseURL    string
	PrivateKey string
	PublicKey  string
}

func (c *Client) SetPrivateKey(key string) {
	c.PrivateKey = key
	c.SetPublicKey()
}
func (c *Client) SetPublicKey() {
	publicKey, err := common.GetPublicKey(c.PrivateKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	if publicKey == "" {
		log.Fatal("public key is empty")
	}
	c.PublicKey = publicKey
}

// makeRequest makes an HTTP request to the API server
func (c *Client) makeRequest(method, endpoint string, queryParams url.Values, requestBody []byte, private bool) (*http.Response, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Add query parameters to the request URL
	q := u.Query()
	for key, values := range queryParams {
		for _, value := range values {
			q.Add(key, value)
		}
	}
	u.RawQuery = q.Encode()

	urlComplete := fmt.Sprintf("%s%s", c.BaseURL, u.String())

	var req *http.Request

	if method == "GET" {
		req, err = http.NewRequest("GET", urlComplete, nil)
	} else if method == "POST" {
		req, err = http.NewRequest("POST", urlComplete, bytes.NewBuffer(requestBody))
	} else {
		return nil, fmt.Errorf("unsupported HTTP method: %s", method)
	}

	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Nonce", generateNonce())
	req.Header.Set("Public-Key", c.PublicKey)

	if private {
		message := createMessage(req, requestBody)

		var signature []byte
		signature, err = common.Sign(c.PrivateKey, message)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		// Add signature and nonce to the request headers
		req.Header.Set("Signature", hex.EncodeToString(signature))

		// Perform any additional operations for private requests
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return resp, nil
}

func createMessage(req *http.Request, requestBody []byte) []byte {
	// Construct the message by concatenating relevant request data
	message := []byte(fmt.Sprintf("%s%s%s%s%s", req.Method, req.URL.Path, req.Header.Get("Nonce"), req.URL.RawQuery, requestBody))
	return message
}

func generateNonce() string {
	// Generate a unique nonce for each request
	// You can use a continuously increasing integer or a timestamp as the nonce
	nonce := time.Now().UnixNano()
	return fmt.Sprintf("%d", nonce)
}

// ForwardMessage forwards a message to the designated server
func (c *Client) ForwardMessage(message *common.Message) error {
	endpoint := "/api/forward"
	requestBody, err := json.Marshal(message)
	if err != nil {
		return err
	}

	resp, err := c.makeRequest(http.MethodPost, endpoint, nil, requestBody, false)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to forward message")
	}
	return nil
}

// GetAllMessages fetches all messages
func (c *Client) GetAllMessages() ([]common.Message, error) {
	endpoint := fmt.Sprintf("/api/users/%s/messages", c.PublicKey)

	resp, err := c.makeRequest(http.MethodGet, endpoint, nil, nil, true)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println(err)
		return nil, fmt.Errorf("failed to fetch all messages")
	}

	var messages []common.Message
	if err = json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return messages, nil
}

// GetEveryMessage fetches all messages (sent/received)
func (c *Client) GetEveryMessage() ([]common.Message, error) {
	endpoint := fmt.Sprintf("/api/users/%s/messages/every", c.PublicKey)

	resp, err := c.makeRequest(http.MethodGet, endpoint, nil, nil, true)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		_, err = io.Copy(os.Stdout, resp.Body)
		if err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}

		return nil, fmt.Errorf("failed to fetch every message")
	}

	var messages []common.Message
	if err = json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return messages, nil
}

// GetUnreadMessages fetches all unread messages
func (c *Client) GetUnreadMessages() ([]common.Message, error) {
	endpoint := fmt.Sprintf("/api/users/%s/messages/unread", c.PublicKey)

	resp, err := c.makeRequest(http.MethodGet, endpoint, nil, nil, true)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println(err)
		return nil, fmt.Errorf("failed to fetch unread messages")
	}

	var unreadMessages []common.Message
	if err = json.NewDecoder(resp.Body).Decode(&unreadMessages); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return unreadMessages, nil
}

// GetReadMessages fetches all read messages
func (c *Client) GetReadMessages() ([]common.Message, error) {
	endpoint := fmt.Sprintf("/api/users/%s/messages/read", c.PublicKey)

	resp, err := c.makeRequest(http.MethodGet, endpoint, nil, nil, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch read messages")
	}

	var readMessages []common.Message
	if err = json.NewDecoder(resp.Body).Decode(&readMessages); err != nil {
		return nil, err
	}
	return readMessages, nil
}

// MarkMessageAsRead marks a message as read
func (c *Client) MarkMessageAsRead(messageID string) error {
	endpoint := fmt.Sprintf("/api/messages/%s/read", messageID)

	resp, err := c.makeRequest(http.MethodPut, endpoint, nil, nil, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to mark message as read")
	}
	return nil
}
