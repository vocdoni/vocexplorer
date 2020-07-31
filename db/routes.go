package db

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/tendermint/go-amino"
	dvotedb "gitlab.com/vocdoni/go-dvote/db"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/types"
)

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
		from, err := strconv.Atoi(froms[0])
		if err != nil {
			log.Error(err)
		}
		keys := list(db, config.ListSize, from, prefs[0])
		log.Debugf("Found %d keys", len(keys))

		var cdc = amino.NewCodec()
		cdc.RegisterConcrete(types.StoreBlock{}, "storeBlock", nil)

		var blocks [config.ListSize]types.StoreBlock
		for i, key := range keys {
			val, err := db.Get([]byte(key))
			if err != nil {
				log.Error(err)
			}
			err = cdc.UnmarshalBinaryLengthPrefixed(val, &blocks[i])
			if err != nil {
				log.Error(err)
			}
		}

		msg, err := json.Marshal(blocks)
		if err != nil {
			log.Error(err)
		}
		fmt.Fprintf(w, string(msg))
		log.Debugf("Sent %d bytes", len(msg))
	}
}

// ValHandler writes the value of the matched key
func ValHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		keys, ok := r.URL.Query()["key"]
		if !ok || len(keys[0]) < 1 {
			log.Errorf("Url Param 'key' is missing")
			return
		}
		val, err := db.Get([]byte(keys[0]))
		if err != nil {
			log.Error(err)
		}
		height, num, err := amino.DecodeInt64(val)
		if err != nil {
			log.Error(err)
		}
		if num <= 1 {
			log.Debug("Could not get block height")
		}

		msg, err := json.Marshal(height)
		if err != nil {
			log.Error(err)
		}
		fmt.Fprintf(w, string(msg))
		log.Debugf("Sent %d bytes", len(msg))
	}
}
