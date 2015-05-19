package bolt

import (
	"fmt"
	"testing"
	"time"

	"github.com/jllopis/try5/account"
)

func TestAccount(t *testing.T) {
	opts := &BoltStoreOptions{
		Dbpath:  "/tmp/test.db",
		Timeout: 5 * time.Second,
	}

	account, err := account.NewAccount("testaccount@dom.local", "Test account", "SuperDifficultPass")
	if err != nil {
		t.Fatal("Error creating account: ", err)
	}

	m := NewBoltStore(opts)
	if m == nil {
		t.Fatal("Error creating boltdb store")
	}

	savedAccount, err := m.SaveAccount(account)
	if err != nil {
		t.Fatal("Error saving account to boltdb store:", err)
	}
	fmt.Printf("Saved account: %#v\n", savedAccount)

	u, err := m.LoadAccount(*savedAccount.UID)
	if err != nil {
		t.Fatal("Error from boltdb store: %v", err)
	}
	if u == nil {
		t.Fatal("Error getting account from boltdb store")
	}
	fmt.Printf("Got from store: %#v\n", u)

	u2, err := m.LoadAccount("")
	if err != nil {
		t.Fatal("Error from boltdb store: %v", err)
	}
	if u2 != nil {
		t.Fatal("Got inexistent account from boltdb store")
		fmt.Printf("Got from store: %#v\n", u2)
	}

}
