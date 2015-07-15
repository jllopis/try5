package bolt

import (
	"os"
	"testing"
	"time"

	"github.com/jllopis/try5/log"
)

var (
	storeConfig = make(map[string]interface{})
	testStore   *Store
)

func TestMain(m *testing.M) {
	setup()
	res := m.Run()
	cleanup()
	os.Exit(res)
}

func setup() {
	log.SetLevel(2) // Info level
	storeConfig["path"] = "./test.db"
	storeConfig["timeout"] = 5 * time.Second

	testStore = NewStore()
	if err := testStore.Dial(storeConfig); err != nil {
		log.LogF("error creating store", "error", err.Error())
	}
}

func cleanup() {
	if _, err := os.Stat(storeConfig["path"].(string)); err == nil {
		os.Remove(storeConfig["path"].(string))
	}
}
