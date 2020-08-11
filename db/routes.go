package db

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"

	dvotedb "gitlab.com/vocdoni/go-dvote/db"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/types"
	"gitlab.com/vocdoni/vocexplorer/util"
	"google.golang.org/protobuf/proto"
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
		w.Write(val)
		log.Debugf("Sent height")
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
func ListBlocksHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
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
		var rawBlocks [config.ListSize][]byte
		for i, hash := range hashes {
			rawBlocks[i], err = db.Get(append([]byte(config.BlockHashPrefix), hash...))
			util.ErrPrint(err)
		}

		msg, err := json.Marshal(rawBlocks)
		w.Write(msg)
		log.Debugf("Sent %d blocks", len(rawBlocks))
	}
}

// GetBlockHandler writes a list of blocks by height
func GetBlockHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ids, ok := r.URL.Query()["id"]
		if !ok || len(ids[0]) < 1 {
			log.Errorf("Url Param 'id' is missing")
			http.Error(w, "Url Param 'id' missing", 400)
			return
		}
		id := ids[0]

		key := []byte(config.BlockHeightPrefix + id)
		hash, err := db.Get(key)
		if err != nil {
			log.Error(err)
		}
		raw, err := db.Get(append([]byte(config.BlockHashPrefix), hash...))
		util.ErrPrint(err)

		w.Write(raw)
		log.Debugf("Sent block %s", id)
	}
}

// ListTxsHandler writes a list of txs starting with the given height key
func ListTxsHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
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
			log.Errorf("No txs available at height %d", from)
			http.Error(w, "No txs available", 404)
			return
		}
		var rawTxs [][]byte
		for _, hash := range hashes {
			raw, err := db.Get(append([]byte(config.TxHashPrefix), hash...))
			util.ErrPrint(err)

			var tx types.StoreTx
			err = proto.Unmarshal(raw, &tx)
			util.ErrPrint(err)

			send := types.SendTx{
				Hash:  hash,
				Store: &tx,
			}
			rawTx, err := proto.Marshal(&send)
			util.ErrPrint(err)
			rawTxs = append(rawTxs, rawTx)
		}

		msg, err := json.Marshal(rawTxs)
		util.ErrPrint(err)
		w.Write(msg)
		log.Debugf("Sent %d txs", len(rawTxs))
	}
}

// GetTxHandler writes the tx corresponding to given height key
func GetTxHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ids, ok := r.URL.Query()["id"]
		if !ok || len(ids[0]) < 1 {
			log.Errorf("Url Param 'id' is missing")
			http.Error(w, "Url Param 'id' missing", 400)
			return
		}
		id := ids[0]
		key := []byte(config.TxHeightPrefix + id)
		hash, err := db.Get(key)
		if err != nil {
			log.Error(err)
		}
		raw, err := db.Get(append([]byte(config.TxHashPrefix), hash...))
		util.ErrPrint(err)

		var tx types.StoreTx
		err = proto.Unmarshal(raw, &tx)
		util.ErrPrint(err)
		height, err := strconv.Atoi(id)
		util.ErrPrint(err)

		send := types.SendTx{
			Hash:  hash,
			Store: &tx,
		}

		rawTx, err := proto.Marshal(&send)
		util.ErrPrint(err)
		w.Write(rawTx)
		log.Debugf("Sent tx %d", height)
	}
}

// TxHashRedirectHandler redirects to the tx corresponding to given height key
func TxHashRedirectHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ids, ok := r.URL.Query()["hash"]
		if !ok || len(ids[0]) < 1 {
			log.Errorf("Url Param 'id' is missing")
			http.Error(w, "Url Param 'id' missing", 400)
			http.Redirect(w, r, r.Header.Get("Referer"), 302)
			return
		}
		id := ids[0]
		hash, err := hex.DecodeString(id)
		if util.ErrPrint(err) {
			http.Redirect(w, r, r.Header.Get("Referer"), 302)
			return
		}
		key := append([]byte(config.TxHashPrefix), hash...)
		has, err := db.Has(key)
		if util.ErrPrint(err) {
			http.Redirect(w, r, r.Header.Get("Referer"), 302)
			return
		}
		if !has {
			log.Errorf("Tx hash key not found")
			http.Redirect(w, r, r.Header.Get("Referer"), 302)
			return
		}
		raw, err := db.Get(key)
		if util.ErrPrint(err) {
			http.Redirect(w, r, r.Header.Get("Referer"), 302)
			return
		}

		var tx types.StoreTx
		err = proto.Unmarshal(raw, &tx)
		if util.ErrPrint(err) {
			http.Redirect(w, r, r.Header.Get("Referer"), 302)
			return
		}

		http.Redirect(w, r, "/txs/"+util.IntToString(tx.TxHeight), 301)
	}
}
