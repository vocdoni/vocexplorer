package db

import (
	"encoding/hex"
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

		var height types.Height
		err = proto.Unmarshal(val, &height)
		util.ErrPrint(err)
		log.Debug("Sent height " + util.IntToString(height.GetHeight()) + " for key " + keys[0])

		w.Write(val)
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
		var rawBlocks types.ItemList
		for _, hash := range hashes {
			rawBlock, err := db.Get(append([]byte(config.BlockHashPrefix), hash...))
			rawBlocks.Items = append(rawBlocks.GetItems(), rawBlock)
			util.ErrPrint(err)
		}

		msg, err := proto.Marshal(&rawBlocks)
		util.ErrPrint(err)
		w.Write(msg)
		log.Debugf("Sent %d blocks", len(rawBlocks.GetItems()))
	}
}

// ListBlocksByValidatorHandler writes a list of blocks which share the given proposer
func ListBlocksByValidatorHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("VALBLOCKS")
		froms, ok := r.URL.Query()["from"]
		if !ok || len(froms[0]) < 1 {
			log.Errorf("Url Param 'from' is missing")
			http.Error(w, "Url Param 'from' missing", 400)
			return
		}
		proposers, ok := r.URL.Query()["proposer"]
		if !ok || len(proposers[0]) < 1 {
			log.Errorf("Url Param 'proposer' is missing")
			http.Error(w, "Url Param 'proposer' missing", 400)
			return
		}
		from, err := strconv.Atoi(froms[0])
		util.ErrPrint(err)

		// latestBlockHeight := &types.Height{}
		// rawValHeight, err := db.Get([]byte(config.LatestBlockHeightKey))
		// util.ErrPrint(err)
		// err = proto.Unmarshal(rawValHeight, latestBlockHeight)
		// util.ErrPrint(err)

		var rawBlocks types.ItemList
		var tempBlock types.StoreBlock
		numBlocks := 0
		// Get blocks by hash, where block proposer matches proposer
		for ; numBlocks < config.ListSize && from > 0; from-- {
			log.Debugf("Getting block at height %d", from)

			hashes := listHashesByHeight(db, 1, from, config.BlockHeightPrefix)
			if len(hashes) == 0 {
				log.Error("No hashes retrieved")
				http.Error(w, "No blocks available", 404)
				return
			}
			for _, hash := range hashes {
				rawBlock, err := db.Get(append([]byte(config.BlockHashPrefix), hash...))
				util.ErrPrint(err)
				err = proto.Unmarshal(rawBlock, &tempBlock)
				util.ErrPrint(err)
				log.Debugf("Found block with proposer %s", tempBlock.GetProposer())
				if util.HexToString(tempBlock.GetProposer()) == proposers[0] {
					rawBlocks.Items = append(rawBlocks.GetItems(), rawBlock)
					numBlocks++
				}
			}
		}

		msg, err := proto.Marshal(&rawBlocks)
		util.ErrPrint(err)
		w.Write(msg)
		log.Debugf("Sent %d blocks by validator %s", len(rawBlocks.GetItems()), proposers[0])
	}
}

// GetBlockHandler writes a block by height
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
		var rawTxs types.ItemList
		for _, hash := range hashes {
			rawTx, err := db.Get(append([]byte(config.TxHashPrefix), hash...))
			util.ErrPrint(err)
			var tx types.StoreTx
			err = proto.Unmarshal(rawTx, &tx)
			util.ErrPrint(err)
			send := types.SendTx{
				Hash:  hash,
				Store: &tx,
			}
			rawSend, err := proto.Marshal(&send)
			util.ErrPrint(err)
			rawTxs.Items = append(rawTxs.GetItems(), rawSend)
			util.ErrPrint(err)
		}
		msg, err := proto.Marshal(&rawTxs)
		util.ErrPrint(err)
		w.Write(msg)
		log.Debugf("Sent %d txs", len(rawTxs.GetItems()))
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
			http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
			return
		}
		id := ids[0]
		hash, err := hex.DecodeString(id)
		if util.ErrPrint(err) {
			http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
			return
		}
		key := append([]byte(config.TxHashPrefix), hash...)
		has, err := db.Has(key)
		if util.ErrPrint(err) {
			http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
			return
		}
		if !has {
			log.Errorf("Tx hash key not found")
			http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
			return
		}
		raw, err := db.Get(key)
		if util.ErrPrint(err) {
			http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
			return
		}

		var tx types.StoreTx
		err = proto.Unmarshal(raw, &tx)
		if util.ErrPrint(err) {
			http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
			return
		}

		http.Redirect(w, r, "/txs/"+util.IntToString(tx.TxHeight), http.StatusPermanentRedirect)
	}
}

// GetValidatorHandler writes the validator corresponding to given address key
func GetValidatorHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ids, ok := r.URL.Query()["id"]
		if !ok || len(ids[0]) < 1 {
			log.Errorf("Url Param 'id' is missing")
			http.Error(w, "Url Param 'id' missing", 400)
			return
		}
		id := ids[0]
		addressBytes, err := hex.DecodeString(id)
		util.ErrPrint(err)
		key := append([]byte(config.ValidatorPrefix), addressBytes...)
		raw, err := db.Get(key)
		util.ErrPrint(err)
		w.Write(raw)
		log.Debugf("Sent validator")
	}
}
