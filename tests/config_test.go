package tests

import (
	"fmt"
	"p2p-chat/common"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	users, err := common.LoadConfig("./data/users.json")
	if err != nil {
		return
	}

	user1 := common.UserWellKnown{
		Username:  "me",
		PublicKey: "0fbbead7194be93d6fc659dd0fe9bdf8b459febfbbfa602322d081eef19688c2",
		Meta:      "Just some more information",
	}

	if user1 != users[0] {
		fmt.Printf("%+v\n", users)
		t.Errorf("user not found")
	}
}
