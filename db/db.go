package db

import (
	"strings"
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
		db.Put([]byte("test:12345"), []byte("abcdef"))
	}
}

// List returns a list of keys matching a given prefix
func list(d *dvotedb.BadgerDB, max int, from, prefix string) (list []string) {
	iter := d.NewIterator().(*dvotedb.BadgerIterator)
	if len(from) > 0 {
		iter.Seek([]byte(from))
	}
	for iter.Next() {
		if max < 1 {
			break
		}
		if strings.HasPrefix(string(iter.Key()), prefix) {
			list = append(list, string(iter.Key()))
			max--
		}
	}
	iter.Release()
	return
}
