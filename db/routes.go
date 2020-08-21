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
		hashes := listItemsByHeight(db, config.ListSize, from, []byte(config.BlockHeightPrefix))
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

		var rawBlocks types.ItemList
		var tempBlock types.StoreBlock
		numBlocks := 0
		// Get blocks by hash, where block proposer matches proposer
		for ; numBlocks < config.ListSize && from > 0; from-- {
			log.Debugf("Getting block at height %d", from)

			hashes := listItemsByHeight(db, 1, from, []byte(config.BlockHeightPrefix))
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

// ListEnvelopesByProcessHandler writes a list of envelopes which share the given process
func ListEnvelopesByProcessHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		froms, ok := r.URL.Query()["from"]
		if !ok || len(froms[0]) < 1 {
			log.Errorf("Url Param 'from' is missing")
			http.Error(w, "Url Param 'from' missing", 400)
			return
		}
		processes, ok := r.URL.Query()["process"]
		if !ok || len(processes[0]) < 1 {
			log.Errorf("Url Param 'process' is missing")
			http.Error(w, "Url Param 'process' missing", 400)
			return
		}
		from, err := strconv.Atoi(froms[0])
		util.ErrPrint(err)

		procEnvHeightMap := getHeightMap(db, config.ProcessEnvelopeHeightMapKey)
		envHeight, ok := procEnvHeightMap.Heights[processes[0]]
		if !ok {
			envHeight = 0
		}
		from = util.Min(from, int(envHeight))

		// Get envelope heights by pid|height
		processBytes, err := hex.DecodeString(processes[0])
		util.ErrPrint(err)
		heights := listItemsByHeight(db, config.ListSize, from, append([]byte(config.EnvPIDPrefix), processBytes...))
		if len(heights) == 0 {
			log.Error("No hashes retrieved")
			http.Error(w, "No blocks available", 404)
			return
		}
		var rawEnvelopes types.ItemList
		// Get envelope packages by globalheight:package
		for _, rawHeight := range heights {
			height := new(types.Height)
			util.ErrPrint(proto.Unmarshal(rawHeight, height))
			rawEnvelope, err := db.Get(append([]byte(config.EnvPackagePrefix), []byte(util.IntToString(height.GetHeight()))...))
			util.ErrPrint(err)
			rawEnvelopes.Items = append(rawEnvelopes.GetItems(), rawEnvelope)
		}

		msg, err := proto.Marshal(&rawEnvelopes)
		util.ErrPrint(err)
		w.Write(msg)
		log.Debugf("Sent %d envelopes by process %s", len(rawEnvelopes.GetItems()), processes[0])
	}
}

// ListEnvelopesHandler writes a list of envelopes
func ListEnvelopesHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		froms, ok := r.URL.Query()["from"]
		if !ok || len(froms[0]) < 1 {
			log.Errorf("Url Param 'from' is missing")
			http.Error(w, "Url Param 'from' missing", 400)
			return
		}
		from, err := strconv.Atoi(froms[0])
		util.ErrPrint(err)

		envHeight := getHeight(db, config.LatestEnvelopeHeightKey, 1)
		from = util.Min(from, int(envHeight.GetHeight()))

		// Get envelope packages
		packages := listItemsByHeight(db, config.ListSize, from, []byte(config.EnvPackagePrefix))
		if len(packages) == 0 {
			log.Error("No envelopes retrieved")
			http.Error(w, "No envelopes available", 404)
			return
		}
		rawEnvelopes := &types.ItemList{Items: packages}

		msg, err := proto.Marshal(rawEnvelopes)
		util.ErrPrint(err)
		w.Write(msg)
		log.Debugf("Sent %d envelopes", len(rawEnvelopes.GetItems()))
	}
}

// GetEnvelopeHandler writes a single envelope
func GetEnvelopeHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		heights, ok := r.URL.Query()["height"]
		if !ok || len(heights[0]) < 1 {
			log.Errorf("Url Param 'id' is missing")
			http.Error(w, "Url Param 'id' missing", 400)
			return
		}
		height, err := strconv.Atoi(heights[0])
		util.ErrPrint(err)

		envHeight := getHeight(db, config.LatestEnvelopeHeightKey, 1)
		if height > int(envHeight.GetHeight()) {
			log.Errorf("Requested envelope does not exist")
			http.Error(w, "Requested envelope does not exist", 400)
			return
		}
		packageKey := append([]byte(config.EnvPackagePrefix), []byte(util.IntToString(height))...)
		rawPackage, err := db.Get(packageKey)
		util.ErrPrint(err)
		w.Write(rawPackage)
		log.Debugf("Sent envelope %d", height)
	}
}

// EnvelopeHeightByProcessHandler writes the number of envelopes which share the given processID
func EnvelopeHeightByProcessHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		processes, ok := r.URL.Query()["process"]
		if !ok || len(processes[0]) < 1 {
			log.Errorf("Url Param 'process' is missing")
			http.Error(w, "Url Param 'process' missing", 400)
			return
		}
		var heightMap types.HeightMap
		valMapKey := []byte(config.ProcessEnvelopeHeightMapKey)
		has, err := db.Has(valMapKey)
		util.ErrPrint(err)
		if has {
			rawValMap, err := db.Get(valMapKey)
			util.ErrPrint(err)
			proto.Unmarshal(rawValMap, &heightMap)
		}
		height, ok := heightMap.Heights[processes[0]]
		if !ok {
			height = 0
		}
		envHeight := &types.Height{Height: int64(height)}
		msg, err := proto.Marshal(envHeight)
		util.ErrPrint(err)
		w.Write(msg)
		log.Debugf("Found %d envelopes by process %s", envHeight.GetHeight(), processes[0])
	}
}

// NumBlocksByValidatorHandler writes the number of blocks which share the given proposer
func NumBlocksByValidatorHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		proposers, ok := r.URL.Query()["proposer"]
		if !ok || len(proposers[0]) < 1 {
			log.Errorf("Url Param 'proposer' is missing")
			http.Error(w, "Url Param 'proposer' missing", 400)
			return
		}
		var heightMap types.HeightMap
		valMapKey := []byte(config.ValidatorHeightMapKey)
		has, err := db.Has(valMapKey)
		util.ErrPrint(err)
		if has {
			rawValMap, err := db.Get(valMapKey)
			util.ErrPrint(err)
			proto.Unmarshal(rawValMap, &heightMap)
		}
		height, ok := heightMap.Heights[proposers[0]]
		if !ok {
			height = 0
		}
		blockHeight := &types.Height{Height: int64(height)}
		msg, err := proto.Marshal(blockHeight)
		util.ErrPrint(err)
		w.Write(msg)
		log.Debugf("Found %d blocks by validator %s", blockHeight.GetHeight(), proposers[0])
	}
}

// // NumBlocksByValidatorHandler writes the number of blocks which share the given proposer
// func NumBlocksByValidatorHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		proposers, ok := r.URL.Query()["proposer"]
// 		if !ok || len(proposers[0]) < 1 {
// 			log.Errorf("Url Param 'proposer' is missing")
// 			http.Error(w, "Url Param 'proposer' missing", 400)
// 			return
// 		}

// 		latestBlockHeight := getHeight(db, config.LatestBlockHeightKey, 1)
// 		numBlocks := int64(0)
// 		complete := make(chan struct{})
// 		// Get blocks by hash, where block proposer matches proposer
// 		height := 1
// 		for ; height < int(latestBlockHeight.GetHeight()); height++ {
// 			go countBlockByValidator(&numBlocks, db, height, proposers[0], complete)
// 		}
// 		// Sync all countBlock routines
// 		for range complete {
// 			if height <= 2 {
// 				break
// 			}
// 			height--
// 		}

// 		blockHeight := &types.Height{Height: int64(numBlocks)}
// 		msg, err := proto.Marshal(blockHeight)
// 		util.ErrPrint(err)
// 		w.Write(msg)
// 		log.Debugf("Found %d blocks by validator %s", blockHeight.GetHeight(), proposers[0])
// 	}
// }

// func countBlockByValidator(numBlocks *int64, db *dvotedb.BadgerDB, height int, proposer string, complete chan struct{}) {
// 	defer func() { complete <- struct{}{} }()
// 	key := []byte(config.BlockHeightPrefix + util.IntToString(height))
// 	has, err := db.Has(key)
// 	if !has || util.ErrPrint(err) {
// 		return
// 	}
// 	hash, err := db.Get(key)
// 	if err != nil {
// 		log.Error(err)
// 	}
// 	var tempBlock types.StoreBlock
// 	rawBlock, err := db.Get(append([]byte(config.BlockHashPrefix), hash...))
// 	util.ErrPrint(err)
// 	err = proto.Unmarshal(rawBlock, &tempBlock)
// 	util.ErrPrint(err)
// 	if util.HexToString(tempBlock.GetProposer()) == proposer {
// 		atomic.AddInt64(numBlocks, 1)
// 	}
// }

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
		hashes := listItemsByHeight(db, config.ListSize, from, []byte(config.TxHeightPrefix))
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
