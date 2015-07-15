/*
Package bolt implements the store.Storer interface and allows the use of a BoltDB
database as persistent storage for try5.
*/
package bolt

import (
	"fmt"
	"time"

	"os"

	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"github.com/jllopis/try5/log"
	"github.com/jllopis/try5/store/manager"
)

// Store holds the resources used to acces a boltdb database.
type Store struct {
	C      *bolt.DB
	status int
	StoreOptions
}

// StoreOptions are the options needed to setup a boltdb database
type StoreOptions struct {
	Dbpath  string
	Timeout time.Duration
}

var (
	tables = [][]byte{[]byte(`accounts`), []byte(`keys`), []byte(`tokens`), []byte(`sessions`), []byte(`certs`)}
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stderr)
	manager.Register("boltdb", &Store{})
}

// NewStore returns an empyt Store object
func NewStore() *Store {
	return &Store{}
}

// Dial creates the database and the needed buckets. It also sets the Store.status
func (s *Store) Dial(options map[string]interface{}) error {
	log.LogD("Dial")
	if options == nil {
		log.LogE("dial options can not be nil", "p", "boltdb", "f", "Dial()")
		return fmt.Errorf("dial options can not be nil")
	}
	log.LogI("Dialing Store", "path", options["path"])
	s.StoreOptions = StoreOptions{Dbpath: options["path"].(string), Timeout: options["timeout"].(time.Duration)}
	db, err := bolt.Open(s.Dbpath, 0600, &bolt.Options{Timeout: s.Timeout * time.Second})
	if err != nil {
		return err
	}
	for _, v := range tables {
		err = db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists(v)
			return err
		})
		if err != nil {
			return nil
		}
	}
	s.C = db
	s.status = manager.CONNECTED
	log.LogI("Dialing done", "status", "CONNECTED", "p", "boltdb", "f", "Dial()")
	return nil
}

// Status returns the current status of the boltdb database. It returns two variables.
// The first one is an integer indicating the state and the second one is the
// string representation (printable) of the status.
func (s *Store) Status() (int, string) {
	return s.status, manager.StatusStr[s.status]
}

// Close effectively closes the database. It must be called when quitting to
// prevent data loss
// TODO(jllopis): call Close() automatically on exit if not called explicitly (context.WithCancel?)
func (s *Store) Close() error {
	s.status = manager.DISCONNECTED
	log.LogI("Store connection closed", "status", "DISCONNECTED", "p", "boltdb", "f", "Close()")
	return s.C.Close()
}
