package db

import (
	"time"

	dvotedb "gitlab.com/vocdoni/go-dvote/db"
)

// NewDB initializes a badger database at the given path
func NewDB(path string) (*dvotedb.BadgerDB, error) {
	return dvotedb.NewBadgerDB(path)
}

// UpdateDB continuously updates the database by calling dvote & tendermint apis
func UpdateDB(db *dvotedb.BadgerDB) {
	for {
		time.Sleep(1 * time.Second)
		db.Put([]byte("test:12345"), []byte("99999999999"))
	}
}
