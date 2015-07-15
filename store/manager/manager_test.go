package manager

import (
	"os"
	"testing"
	"time"
)

var (
	storeConfig = make(map[string]interface{})
)

func TestMain(m *testing.M) {
	setup()
	res := m.Run()
	cleanup()
	os.Exit(res)
}

func TestManager(t *testing.T) {
	m, err := NewManager(storeConfig, "boltdb")
	if err != nil {
		t.Fatalf("Could not create Manager: %s", err)
	}
	err = m.Init()
	if err != nil {
		t.Fatalf("Could not Init manager: %s", err)
	}
}

func setup() {
	storeConfig["path"] = "./test.db"
	storeConfig["timeout"] = 5 * time.Second
}

func cleanup() {

}
