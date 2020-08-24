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
			http.Error(w, "Url Param 'key' missing", http.StatusBadRequest)
			return
		}
		val, err := db.Get([]byte(keys[0]))
		if err != nil {
			log.Error(err)
			http.Error(w, "Key not found", http.StatusInternalServerError)
			return
		}

		var height types.Height
		err = proto.Unmarshal(val, &height)
		if err != nil {
			log.Error(err)
			http.Error(w, "Height not found", http.StatusInternalServerError)
			return
		}
		log.Debug("Sent height " + util.IntToString(height.GetHeight()) + " for key " + keys[0])

		w.Write(val)
	}
}

func buildItemByHeightHandler(db *dvotedb.BadgerDB, heightKey, key string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		heights, ok := r.URL.Query()["height"]
		if !ok || len(heights[0]) < 1 {
			log.Errorf("Url Param 'id' is missing")
			http.Error(w, "Url Param 'id' missing", http.StatusBadRequest)
			return
		}
		height, err := strconv.Atoi(heights[0])
		util.ErrPrint(err)

		envHeight := getHeight(db, heightKey, 0)
		if height > int(envHeight.GetHeight()) {
			log.Errorf("Requested item does not exist")
			http.Error(w, "Requested item does not exist", http.StatusInternalServerError)
			return
		}
		processKey := append([]byte(key), []byte(util.IntToString(height))...)
		rawProcess, err := db.Get(processKey)
		if err != nil {
			log.Error(err)
			http.Error(w, "item not found", http.StatusInternalServerError)
			return
		}
		w.Write(rawProcess)
		log.Debugf("sent item with key %s, height %d", key, height)
	}
}

func buildListItemsHandler(db *dvotedb.BadgerDB, key string, getItem func(key []byte) ([]byte, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		froms, ok := r.URL.Query()["from"]
		if !ok || len(froms[0]) < 1 {
			log.Errorf("Url Param 'from' is missing")
			http.Error(w, "Url Param 'from' missing", http.StatusBadRequest)
			return
		}
		from, err := strconv.Atoi(froms[0])
		util.ErrPrint(err)
		items := listItemsByHeight(db, config.ListSize, from, []byte(key))
		if len(items) == 0 {
			http.Error(w, "No items available", http.StatusInternalServerError)
			return
		}
		var itemList types.ItemList
		for _, rawItem := range items {
			newItem, err := getItem(rawItem)
			util.ErrPrint(err)
			if len(newItem) > 0 {
				rawItem = newItem
			}
			itemList.Items = append(itemList.GetItems(), rawItem)
		}

		msg, err := proto.Marshal(&itemList)
		if err != nil {
			log.Error(err)
			http.Error(w, "Unable to encode data", http.StatusInternalServerError)
			return
		}
		w.Write(msg)
		log.Debugf("Sent %d items", len(itemList.GetItems()))
	}
}

func buildListItemsByParent(db *dvotedb.BadgerDB, parentName, heightMapKey, getHeightPrefix, itemPrefix string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		froms, ok := r.URL.Query()["from"]
		if !ok || len(froms[0]) < 1 {
			log.Errorf("Url Param 'from' is missing")
			http.Error(w, "Url Param 'from' missing", http.StatusBadRequest)
			return
		}
		parents, ok := r.URL.Query()[parentName]
		if !ok || len(parents[0]) < 1 {
			log.Errorf("Url Param '" + parentName + "' is missing")
			http.Error(w, "Url Param '"+parentName+"' is missing", http.StatusBadRequest)
			return
		}
		from, err := strconv.Atoi(froms[0])
		util.ErrPrint(err)

		heightMap := getHeightMap(db, heightMapKey)
		itemHeight, ok := heightMap.Heights[parents[0]]
		if !ok {
			itemHeight = 0
		}
		from = util.Max(util.Min(from, int(itemHeight)), config.ListSize)

		// Get keys
		parentBytes, err := hex.DecodeString(parents[0])
		util.ErrPrint(err)
		keys := listItemsByHeight(db, config.ListSize, from, append([]byte(getHeightPrefix), parentBytes...))
		if len(keys) == 0 {
			log.Error("No keys retrieved")
			http.Error(w, "No items available", 404)
			return
		}
		var rawItems types.ItemList
		// Get packages by height:package
		for _, rawKey := range keys {
			height := new(types.Height)
			util.ErrPrint(proto.Unmarshal(rawKey, height))
			rawPackage, err := db.Get(append([]byte(itemPrefix), []byte(util.IntToString(height.GetHeight()))...))
			util.ErrPrint(err)
			rawItems.Items = append(rawItems.GetItems(), rawPackage)
		}

		msg, err := proto.Marshal(&rawItems)
		if err != nil {
			log.Error(err)
			http.Error(w, "Unable to encode data", http.StatusInternalServerError)
			return
		}
		util.ErrPrint(err)
		w.Write(msg)
		log.Debugf("Sent %d items by %s %s", len(rawItems.GetItems()), parentName, parents[0])
	}
}

// ListBlocksHandler writes a list of blocks by height
func ListBlocksHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildListItemsHandler(db,
		config.BlockHeightPrefix,
		func(key []byte) ([]byte, error) {
			return db.Get(append([]byte(config.BlockHashPrefix), key...))
		})
}

// ListBlocksByValidatorHandler writes a list of blocks which share the given proposer
func ListBlocksByValidatorHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildListItemsByParent(db, "proposer", config.ValidatorHeightMapKey, config.BlockByValidatorPrefix, config.BlockHashPrefix)
}

// ListEnvelopesByProcessHandler writes a list of envelopes which share the given process
func ListEnvelopesByProcessHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildListItemsByParent(db, "process", config.ProcessEnvelopeHeightMapKey, config.EnvPIDPrefix, config.EnvPackagePrefix)
}

// ListEnvelopesHandler writes a list of envelopes
func ListEnvelopesHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildListItemsHandler(db, config.EnvPackagePrefix, func(key []byte) ([]byte, error) {
		return []byte{}, nil
	})
}

// GetEnvelopeHandler writes a single envelope
func GetEnvelopeHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildItemByHeightHandler(db, config.LatestEnvelopeHeightKey, config.EnvPackagePrefix)
}

// EnvelopeHeightByProcessHandler writes the number of envelopes which share the given processID
func EnvelopeHeightByProcessHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		processes, ok := r.URL.Query()["process"]
		if !ok || len(processes[0]) < 1 {
			log.Errorf("Url Param 'process' is missing")
			http.Error(w, "Url Param 'process' missing", http.StatusBadRequest)
			return
		}
		var heightMap types.HeightMap
		valMapKey := []byte(config.ProcessEnvelopeHeightMapKey)
		has, err := db.Has(valMapKey)
		if err != nil || !has {
			log.Error("No envelope height not found")
			http.Error(w, "No envelopes available", http.StatusInternalServerError)
			return
		}
		if has {
			rawValMap, err := db.Get(valMapKey)
			util.ErrPrint(err)
			err = proto.Unmarshal(rawValMap, &heightMap)
			util.ErrPrint(err)
		}
		height, ok := heightMap.Heights[processes[0]]
		if !ok {
			height = 0
		}
		envHeight := &types.Height{Height: int64(height)}
		msg, err := proto.Marshal(envHeight)
		if err != nil {
			log.Error(err)
			http.Error(w, "Unable to encode data", http.StatusInternalServerError)
			return
		}
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

// GetBlockHandler writes a block by height
func GetBlockHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		heights, ok := r.URL.Query()["height"]
		if !ok || len(heights[0]) < 1 {
			log.Errorf("Url Param 'height' is missing")
			http.Error(w, "Url Param 'height' missing", 400)
			return
		}
		height := heights[0]
		key := []byte(config.BlockHeightPrefix + height)
		hash, err := db.Get(key)
		if err != nil {
			log.Error(err)
		}
		raw, err := db.Get(append([]byte(config.BlockHashPrefix), hash...))
		util.ErrPrint(err)

		w.Write(raw)
		log.Debugf("Sent block %s", height)
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
			http.Error(w, "Url Param 'id' missing", http.StatusBadRequest)
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

// GetEntityHandler writes a single entity
func GetEntityHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildItemByHeightHandler(db, config.LatestEntityHeight, config.EntityIDPrefix)
}

// GetProcessHandler writes a single process
func GetProcessHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildItemByHeightHandler(db, config.LatestProcessHeight, config.ProcessIDPrefix)
}

// ListEntitiesHandler writes a list of entities from 'from'
func ListEntitiesHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildListItemsHandler(db, config.EntityIDPrefix, func(key []byte) ([]byte, error) {
		return []byte{}, nil
	})
}

// ListProcessesHandler writes a list of processes from 'from'
func ListProcessesHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildListItemsHandler(db, config.ProcessIDPrefix, func(key []byte) ([]byte, error) {
		return []byte{}, nil
	})
}

// ListProcessesByEntityHandler writes a list of processes belonging to 'entity'
func ListProcessesByEntityHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildListItemsByParent(db, "entity", config.EntityProcessHeightMapKey, config.ProcessByEntityPrefix, config.ProcessIDPrefix)
}
