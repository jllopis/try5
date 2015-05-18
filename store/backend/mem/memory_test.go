package mem

import (
	"fmt"
	"testing"

	"github.com/jllopis/try5/account"
)

func TestAccount(t *testing.T) {
	account, err := account.NewAccount("testaccount@dom.local", "Test account", "SuperDifficultPass")
	if err != nil {
		t.Fatal("Error creating account: ", err)
	}

	m := NewMemStore()
	if m == nil {
		t.Fatal("Error creating mem store")
	}

	savedAccount, err := m.SaveAccount(account)
	if err != nil {
		t.Fatal("Error saving account to mem store:", err)
	}
	fmt.Printf("Saved account: %#v\n", savedAccount)

	u := m.LoadAccount(savedAccount.UID)
	if u == nil {
		t.Fatal("Error getting account from mem store")
	}
	fmt.Printf("Got from store: %#v\n", u)

	u2 := m.LoadAccount("")
	if u2 != nil {
		t.Fatal("Got inexistent account from mem store")
		fmt.Printf("Got from store: %#v\n", u2)
	}

}
