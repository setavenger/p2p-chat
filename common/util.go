package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type UserWellKnown struct {
	Username  string `json:"username,omitempty" gorm:"column:username;index"`
	PublicKey string `json:"public_key,omitempty" gorm:"column:public_key;index"`
	Meta      string `json:"meta,omitempty" gorm:"column:meta"`
}

func ParseUserName(username string) (string, string, error) {
	parts := strings.Split(username, "@")
	if len(parts) != 2 {
		return "", "", errors.New("username bad format")
	}
	return parts[0], parts[1], nil
}

func GenerateWellKnownAddress(username string) (string, error) {
	userNameShort, domain, err := ParseUserName(username)
	if err != nil {
		return "", err
	}
	urlBase := fmt.Sprintf("http://%s/.well-known/p2pchat/%s", domain, userNameShort)
	return urlBase, nil
}

func GetPublicKeyForUsername(username string) (*UserWellKnown, error) {

	address, err := GenerateWellKnownAddress(username)
	if err != nil {
		return nil, err
	}

	var req *http.Request

	req, err = http.NewRequest("GET", address, nil)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

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

		return nil, fmt.Errorf("failed to fetch user data")
	}

	var userWellKnown UserWellKnown
	if err = json.NewDecoder(resp.Body).Decode(&userWellKnown); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &userWellKnown, nil
}
