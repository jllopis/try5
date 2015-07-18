package main

import (
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/jllopis/try5/account"
	logger "github.com/jllopis/try5/log"
	_ "github.com/jllopis/try5/store/backend/boltdb"
	"github.com/jllopis/try5/store/manager"
)

var (
	storeConfig = make(map[string]interface{})
)

func init() {
	logger.SetLevel(logrus.DebugLevel)
	storeConfig["path"] = "./test.db"
	storeConfig["timeout"] = 5 * time.Second
}

func main() {
	logger.LogI("Starting app", "pkg", "main", "func", "main()", "storeConfig", storeConfig)
	m, err := manager.NewManager(storeConfig, "boltdb")
	if err != nil {
		logger.LogE("Could not create Manager", "error", err)
	}
	defer m.Close()

	err = m.Init()
	if err != nil {
		logger.LogF("Could not Init manager", "pkg", "main", "func", "Init()", "error", err)
	}

	// Create new account
	uid, err := m.CreateAccount("test@email.com", "Test Account", "secret")
	if err != nil {
		logger.LogI("error creating account", "pkg", "main", "func", "main()", "error", err)
	} else {
		logger.LogI("new account created", "pkg", "main", "func", "main()", "uid", uid)
	}
	// Load it
	acc1, err := m.LoadAccount(uid)
	if err != nil {
		logger.LogE("error loading account", "pkg", "main", "func", "main()", "error", err)
	} else {
		logger.LogI("new account loaded", "pkg", "main", "func", "main()", "uid", acc1)
	}
	// Save another
	em := "test2@my.domain"
	n := "Second Test Account"
	p := "secret"
	acc2 := &account.Account{Email: &em, Name: &n, Password: &p}
	err = m.SaveAccount(acc2)
	if err != nil {
		logger.LogE("error saving account", "pkg", "main", "func", "main()", "error", err)
	} else {
		logger.LogI("account saved", "pkg", "main", "func", "main()", "uid", *acc2.UID)
	}

	// Get all accounts
	accList, err := m.LoadAllAccounts()
	if err != nil {
		logger.LogE("error loading all accounts", "pkg", "main", "func", "main()", "error", err)
	} else {
		logger.LogI("all accounts loaded", "pkg", "main", "func", "main()", "len", len(accList))
	}

	// Delete Account
	err = m.DeleteAccount(*acc2.UID)
	if err != nil {
		logger.LogE("error deleting acc2", "pkg", "main", "func", "main()", "uid", *acc2.UID, "error", err)
	} else {
		logger.LogI("acc2 deleted", "pkg", "main", "func", "main()", "uid", *acc2.UID)
	}

	// Check if it exists
	if exist, err := m.ExistAccount("uid=" + *acc2.UID); err != nil {
		logger.LogE("error checking acc2 existence", "pkg", "main", "func", "main()", "uid", *acc2.UID, "error", err)
	} else {
		if exist {
			logger.LogE("acc2 still exist", "pkg", "main", "func", "main()", "uid", *acc2.UID)
		} else {
			logger.LogI("acc2 does not exist", "pkg", "main", "func", "main()", "uid", *acc2.UID)
		}
	}

	// Free resources
	if _, err := os.Stat(storeConfig["path"].(string)); err == nil {
		os.Remove(storeConfig["path"].(string))
	}
}
