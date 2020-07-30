package db

import (
	"encoding/json"
	"net/http"

	dvotedb "gitlab.com/vocdoni/go-dvote/db"
)

// DumpHandler dumpts the ocntents of db
func DumpHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		val, err := db.Get([]byte("test:12345"))
		if err != nil {
			panic(err)
		}
		if err := json.NewEncoder(w).Encode(string(val)); err != nil {
			panic(err)
		}
	}
}
