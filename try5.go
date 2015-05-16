package try5

import "github.com/jllopis/try5/store"

type Try5 struct {
	store.Storer
}

// New initialize the storage with a compliant database and
// return a Try5 struct to work with
func New(s store.Storer) *Try5 {
	return &Try5{s}
}
