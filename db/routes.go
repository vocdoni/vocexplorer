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

// HeightHandler writes the int64 height value corresponding to given key
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
			log.Debug("Could not get height")
		}

		msg, err := json.Marshal(height)
		if err != nil {
			log.Error(err)
		}
		fmt.Fprintf(w, string(msg))
		log.Debugf("Sent %d bytes", len(msg))
	}
}

// // HashHandler writes the hash value corresponding to the given key
// func HashHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "text/plain")
// 		// w.Header().Set("Content-Length", util.IntToString(len("Key not found")))
// 		keys, ok := r.URL.Query()["key"]
// 		if !ok || len(keys[0]) < 1 {
// 			log.Errorf("Url Param 'key' is missing")
// 			fmt.Fprintf(w, "Key not found")
// 			// http.Error(w, "Url Param 'key' missing", 400)
// 			return
// 		}
// 		hash, err := db.Get([]byte(keys[0]))
// 		if err != nil {
// 			log.Error(err)
// 			fmt.Fprintf(w, "Key not found")
// 			// http.Error(w, "Key not found", 400)
// 			return
// 		}
// 		w.Header().Set("Content-Length", util.IntToString(len(tmbytes.HexBytes(hash).String())))
// 		log.Infof("Hash: %s", tmbytes.HexBytes(hash).String())
// 		fmt.Fprintf(w, tmbytes.HexBytes(hash).String())
// 		log.Debugf("Sent %d bytes", len(hash))
// 	}
// }

// ListBlocksHandler writes a list of blocks by height
func ListBlocksHandler(db *dvotedb.BadgerDB, cdc *amino.Codec) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		froms, ok := r.URL.Query()["from"]
		if !ok || len(froms[0]) < 1 {
			log.Errorf("Url Param 'from' is missing")
			http.Error(w, "Url Param 'from' missing", 400)
			return
		}
		from, err := strconv.Atoi(froms[0])
		util.ErrPrint(err)
		hashes := listHashesByHeight(db, config.ListSize, from, config.BlockHeightPrefix)
		if len(hashes) == 0 {
			http.Error(w, "No blocks available", 404)
			return
		}
		var blocks [config.ListSize]types.StoreBlock
		for i, hash := range hashes {
			raw, err := db.Get(append([]byte(config.BlockHashPrefix), hash...))
			util.ErrPrint(err)

			err = cdc.UnmarshalBinaryLengthPrefixed(raw, &blocks[i])
			util.ErrPrint(err)
		}

		msg, err := json.Marshal(blocks)
		if err != nil {
			log.Error(err)
		}
		fmt.Fprintf(w, string(msg))
		log.Debugf("Sent %d blocks", len(blocks))
	}
}

// ListTxsHandler writes the tx corresponding to given key
func ListTxsHandler(db *dvotedb.BadgerDB, cdc *amino.Codec) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		froms, ok := r.URL.Query()["from"]
		if !ok || len(froms[0]) < 1 {
			log.Errorf("Url Param 'from' is missing")
			http.Error(w, "Url Param 'from' missing", 400)
			return
		}
		from, err := strconv.Atoi(froms[0])
		util.ErrPrint(err)
		hashes := listHashesByHeight(db, config.ListSize, from, config.TxHeightPrefix)
		if len(hashes) == 0 {
			http.Error(w, "No txs available", 404)
			return
		}
		var txs []types.SendTx
		for i, hash := range hashes {
			raw, err := db.Get(append([]byte(config.TxHashPrefix), hash...))
			util.ErrPrint(err)

			var tx types.StoreTx
			err = cdc.UnmarshalBinaryLengthPrefixed(raw, &tx)
			util.ErrPrint(err)

			send := types.SendTx{
				Hash:   hash,
				Height: from - i,
				Store:  tx,
			}
			txs = append(txs, send)
		}

		msg, err := json.Marshal(txs)
		if err != nil {
			log.Error(err)
		}
		fmt.Fprintf(w, string(msg))
		log.Debugf("Sent %d txs", len(txs))
	}
}
