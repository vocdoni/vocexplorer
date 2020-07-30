package db

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/config"

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

// ListHandler writes a list of values corresponding to keys which match {prefix}
func ListHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		prefs, ok := r.URL.Query()["prefix"]
		if !ok || len(prefs[0]) < 1 {
			log.Errorf("Url Param 'prefix' is missing")
			return
		}
		froms, ok := r.URL.Query()["from"]
		if !ok || len(froms[0]) < 1 {
			log.Errorf("Url Param 'from' is missing")
			return
		}
		keys := list(db, config.ListSize, froms[0], prefs[0])

		var vals []string
		for _, key := range keys {
			val, err := db.Get([]byte(key))

			if err != nil {
				log.Error(err)
			} else {
				vals = append(vals, string(val))
			}
		}

		msg, err := json.Marshal(vals)
		if err != nil {
			log.Error(err)
		}
		log.Infof("Posting message: " + string(msg))

		fmt.Fprintf(w, string(msg))
	}
}
