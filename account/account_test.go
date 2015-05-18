package account

import (
	"fmt"
	"testing"
)

func TestAccount(t *testing.T) {
	account, err := NewAccount("testuser@dom.local", "Test Account", "SuperDifficultPass")
	if err != nil {
		t.Fatal("Error creating account: ", err)
	}
	fmt.Printf("%#v\n", account)
}
