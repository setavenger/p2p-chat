package common

import (
	"encoding/json"
	"fmt"
	"os"
)

func LoadConfig(filename string) ([]UserWellKnown, error) {

	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error when opening file: ", err)
		return nil, err
	}

	var config []UserWellKnown

	err = json.Unmarshal(content, &config)
	if err != nil {
		fmt.Println("Error when opening unmarshalling file: ", err)
		return nil, err
	}

	return config, nil
}
