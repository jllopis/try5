package mem

import (
	"fmt"
	"testing"

	"github.com/jllopis/try5/user"
)

func TestAccount(t *testing.T) {
	user, err := user.NewUser("testuser@dom.local", "Test User", "SuperDifficultPass")
	if err != nil {
		t.Fatal("Error creating account: ", err)
	}

	m := NewMemStore()
	if m == nil {
		t.Fatal("Error creating mem store")
	}

	savedUser, err := m.SaveUser(user)
	if err != nil {
		t.Fatal("Error saving user to mem store:", err)
	}
	fmt.Printf("Saved user: %#v\n", savedUser)

	u := m.LoadUser(savedUser.UID)
	if u == nil {
		t.Fatal("Error getting user from mem store")
	}
	fmt.Printf("Got from store: %#v\n", u)

	u2 := m.LoadUser("")
	if u2 != nil {
		t.Fatal("Got inexistent user from mem store")
		fmt.Printf("Got from store: %#v\n", u2)
	}

}
