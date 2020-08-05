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
	"gitlab.com/vocdoni/vocexplorer/util"
)

// ListBlocksHandler writes a list of blocks corresponding to keys which match {prefix}
func ListBlocksHandler(db *dvotedb.BadgerDB, cdc *amino.Codec) func(w http.ResponseWriter, r *http.Request) {
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
		keys := listKeysByHeight(db, config.ListSize, from, prefs[0])
		log.Debugf("Found %d keys", len(keys))

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

// HeightHandler writes the int64 value corresponding to given key
func HeightHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		keys, ok := r.URL.Query()["key"]
		if !ok || len(keys[0]) < 1 {
			log.Errorf("Url Param 'key' is missing")
			http.Error(w, "Url Param 'key' missing", 400)
			return
		}
		val, err := db.Get([]byte(keys[0]))
		if err != nil {
			log.Error(err)
			http.Error(w, "Key not found", 404)
			return
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

// ListTxsHandler writes the tx corresponding to given key
func ListTxsHandler(db *dvotedb.BadgerDB, cdc *amino.Codec) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		heights, ok := r.URL.Query()["height"]
		if !ok || len(heights[0]) < 1 {
			log.Errorf("Url Param 'height' is missing")
			http.Error(w, "Url Param 'height' missing", 400)
			return
		}
		indexes, ok := r.URL.Query()["index"]
		if !ok || len(indexes[0]) < 1 {
			log.Errorf("Url Param 'index' is missing")
			http.Error(w, "Url Param 'index' missing", 400)
			return
		}
		height, err := strconv.Atoi(heights[0])
		util.ErrPrint(err)
		index, err := strconv.Atoi(indexes[0])
		util.ErrPrint(err)
		hashes := listTxKeys(db, config.ListSize, height, index)
		if len(hashes) == 0 {
			http.Error(w, "No txs available", 404)
			return
		}
		var txs []types.SendTx
		for _, hash := range hashes {
			raw, err := db.Get(append([]byte(config.TxHashPrefix), hash...))
			if err != nil {
				log.Error(err)
			}
			var tx types.StoreTx
			err = cdc.UnmarshalBinaryLengthPrefixed(raw, &tx)
			if err != nil {
				log.Error(err)
			}
			send := types.SendTx{
				Hash:  hash,
				Store: tx,
			}
			txs = append(txs, send)
		}

		msg, err := json.Marshal(txs)
		if err != nil {
			log.Error(err)
		}
		fmt.Fprintf(w, string(msg))
		log.Debugf("Sent %d bytes", len(msg))
	}
}
