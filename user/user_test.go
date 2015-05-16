package user

import (
	"fmt"
	"testing"
)

func TestAccount(t *testing.T) {
	user, err := NewUser("testuser@dom.local", "Test User", "SuperDifficultPass")
	if err != nil {
		t.Fatal("Error creating account: ", err)
	}
	fmt.Printf("%#v\n", user)
}
