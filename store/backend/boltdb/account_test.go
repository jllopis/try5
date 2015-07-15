package bolt

import (
	"reflect"
	"testing"

	"github.com/jllopis/try5/account"
	"github.com/jllopis/try5/log"
	"github.com/jllopis/try5/tryerr"
)

var (
	acc *account.Account
)

func TestSaveAccount(t *testing.T) {
	if testStore == nil {
		log.LogF("testStore is nil")
	}

	var err error
	acc, err = account.NewAccount("testaccount@dom.local", "Test account", "SuperDifficultPass")
	if err != nil {
		t.Fatal("Error creating account: ", err)
	}

	err = testStore.SaveAccount(acc)
	if err != nil {
		t.Fatal("Error saving account to boltdb store:", err)
	}
}

func TestLoadAccount(t *testing.T) {
	u, err := testStore.LoadAccount(*acc.UID)
	if err != nil {
		t.Fatalf("Error from boltdb store: %v", err)
	}
	if u == nil {
		t.Fatal("Error getting account from boltdb store")
	}
	t.Logf("Got from store: %#v\n", u)

	_, err = testStore.LoadAccount("")
	if err != nil && err != tryerr.ErrAccountNotFound {
		t.Fatalf("Error from boltdb store: %v (%v)", err, reflect.TypeOf(err))
	}
}

func TestLoadAllAccounts(t *testing.T) {
	a := []map[string]string{
		{"m": "a1@a.b", "n": "a1", "p": "pass1"},
		{"m": "a2@a.b", "n": "a2", "p": "pass2"},
		{"m": "a3@a.b", "n": "a3", "p": "pass3"},
		{"m": "a4@a.b", "n": "a4", "p": "pass4"},
	}
	for _, e := range a {
		x, err := account.NewAccount(e["m"], e["n"], e["p"])
		if err != nil {
			t.Fatal("error creating account in TestLoadAccounts")
		}
		testStore.SaveAccount(x)
	}
	all, err := testStore.LoadAllAccounts()
	if err != nil {
		t.Fatalf("error in LoadAllAccounts: %v", err)
	}
	if len(all) != 5 {
		t.Fatalf("LoadAllAccounts: Expected=%d Got=%d", 5, len(all))
	}
}

func TestGetAccountByEmail(t *testing.T) {
	if _, err := testStore.GetAccountByEmail("a1@a.b"); err != nil {
		t.Fatalf("error getting account by email: %v", err)
	}
}

func TestDeleteAccount(t *testing.T) {
	if err := testStore.DeleteAccount(*acc.UID); err != nil {
		t.Fatalf("Error deleting account (%v): %v", *acc.UID, err)
	}
}
