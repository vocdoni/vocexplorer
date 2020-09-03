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

//PingHandler responds to a ping
func PingHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	}
}

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

// HeightMapHandler writes the string:int64 height map corresponding to given key
func HeightMapHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
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

		var heightMap types.HeightMap
		err = proto.Unmarshal(val, &heightMap)
		if err != nil {
			log.Error(err)
			http.Error(w, "Height map not found", http.StatusInternalServerError)
			return
		}
		log.Debug("Sent height map for key " + keys[0])
		w.Write(val)
	}
}

func buildItemByIDHandler(db *dvotedb.BadgerDB, IDName, itemPrefix string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ids, ok := r.URL.Query()[IDName]
		if !ok || len(ids[0]) < 1 {
			log.Errorf("Url Param '" + IDName + "' is missing")
			http.Error(w, "Url Param '"+IDName+"' missing", http.StatusBadRequest)
			return
		}
		id := ids[0]
		addressBytes, err := hex.DecodeString(id)
		if err != nil {
			log.Error(err)
		}
		key := append([]byte(itemPrefix), addressBytes...)
		raw, err := db.Get(key)
		if err != nil {
			log.Error(err)
		}
		w.Write(raw)
		log.Debugf("Sent Item")
	}
}

func buildItemByHeightHandler(db *dvotedb.BadgerDB, heightKey, key string, getItem func(key []byte) ([]byte, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		heights, ok := r.URL.Query()["height"]
		if !ok || len(heights[0]) < 1 {
			log.Errorf("Url Param 'height' is missing")
			http.Error(w, "Url Param 'height' missing", http.StatusBadRequest)
			return
		}
		height, err := strconv.Atoi(heights[0])
		if err != nil {
			log.Error(err)
			http.Error(w, "Cannot decode height", http.StatusInternalServerError)
			return
		}

		envHeight := getHeight(db, heightKey, 0)
		if height > int(envHeight.GetHeight()) {
			log.Errorf("Requested item does not exist")
			http.Error(w, "Requested item does not exist", http.StatusInternalServerError)
			return
		}
		itemKey := append([]byte(key), util.EncodeInt(height)...)
		rawItem, err := db.Get(itemKey)
		if err != nil {
			log.Error(err)
			http.Error(w, "item not found", http.StatusInternalServerError)
			return
		}
		if getItem != nil {
			rawItem, err = getItem(rawItem)
			if err != nil {
				log.Error(err)
				http.Error(w, "item not found", http.StatusInternalServerError)
				return
			}
		}
		w.Write(rawItem)
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
		if err != nil {
			log.Error(err)
		}
		items := listItemsByHeight(db, config.ListSize, from, []byte(key))
		if len(items) == 0 {
			log.Error("Retrieved no items")
			http.Error(w, "No items available", http.StatusInternalServerError)
			return
		}
		var itemList types.ItemList
		for _, rawItem := range items {
			if getItem != nil {
				rawItem, err = getItem(rawItem)
				if err != nil {
					log.Error(err)
					http.Error(w, "item not found", http.StatusInternalServerError)
					return
				}
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

func buildListItemsByParent(db *dvotedb.BadgerDB, parentName, heightMapKey, getHeightPrefix, itemPrefix string, marshalHeight bool) func(w http.ResponseWriter, r *http.Request) {
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
		if err != nil {
			log.Error(err)
		}

		heightMap := getHeightMap(db, heightMapKey)
		itemHeight, ok := heightMap.Heights[parents[0]]
		if !ok {
			log.Error("Parent does not exist")
			http.Error(w, "No items available", http.StatusInternalServerError)
			return
		}
		from = util.Min(from, int(itemHeight))

		// Get keys
		parentBytes, err := hex.DecodeString(parents[0])
		if err != nil {
			log.Error(err)
		}
		keys := listItemsByHeight(db, config.ListSize, from, append([]byte(getHeightPrefix), parentBytes...))
		if len(keys) == 0 {
			log.Error("No keys retrieved")
			http.Error(w, "No items available", http.StatusInternalServerError)
			return
		}
		var rawItems types.ItemList
		// Get packages by height:package
		for _, rawKey := range keys {
			var rawPackage []byte
			if marshalHeight {
				height := new(types.Height)
				if err := proto.Unmarshal(rawKey, height); err != nil {
					log.Error(err)
				}
				log.Debugf("Getting item from height %d", height.GetHeight())
				rawPackage, err = db.Get(append([]byte(itemPrefix), util.EncodeInt(height.GetHeight())...))
			} else {
				rawPackage, err = db.Get(append([]byte(itemPrefix), rawKey...))
			}
			if err != nil {
				log.Error(err)
			} else {
				rawItems.Items = append(rawItems.GetItems(), rawPackage)
			}
		}
		if len(rawItems.Items) <= 0 {
			http.Error(w, "No items found", http.StatusInternalServerError)
			return
		}
		msg, err := proto.Marshal(&rawItems)
		if err != nil {
			log.Error(err)
			http.Error(w, "Unable to encode data", http.StatusInternalServerError)
			return
		}
		if err != nil {
			log.Error(err)
		}
		w.Write(msg)
		log.Debugf("Sent %d items by %s %s", len(rawItems.GetItems()), parentName, parents[0])
	}
}

func buildHeightByParentHandler(db *dvotedb.BadgerDB, parentName, heightMapKey string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		parents, ok := r.URL.Query()[parentName]
		if !ok || len(parents[0]) < 1 {
			log.Errorf("Url Param '" + parentName + "' is missing")
			http.Error(w, "Url Param 'process' missing", http.StatusBadRequest)
			return
		}
		var heightMap types.HeightMap
		has, err := db.Has([]byte(heightMapKey))
		if err != nil || !has {
			log.Error("No item height not found")
			http.Error(w, "No items available", http.StatusInternalServerError)
			return
		}
		if has {
			rawHeightMap, err := db.Get([]byte(heightMapKey))
			if err != nil {
				log.Error(err)
			}
			err = proto.Unmarshal(rawHeightMap, &heightMap)
			if err != nil {
				log.Error(err)
			}
		}
		height, ok := heightMap.Heights[parents[0]]
		if !ok {
			height = 0
		}
		sendHeight := &types.Height{Height: int64(height)}
		msg, err := proto.Marshal(sendHeight)
		if err != nil {
			log.Error(err)
			http.Error(w, "Unable to encode data", http.StatusInternalServerError)
			return
		}
		w.Write(msg)
		log.Debugf("Found %d items by %s %s", sendHeight.GetHeight(), parentName, parents[0])
	}
}

// GetEnvelopeHandler writes a single envelope
func GetEnvelopeHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildItemByHeightHandler(db, config.LatestEnvelopeHeightKey, config.EnvPackagePrefix, nil)
}

// GetBlockHandler writes a block by height
func GetBlockHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildItemByHeightHandler(db,
		config.LatestBlockHeightKey,
		config.BlockHeightPrefix,
		func(key []byte) ([]byte, error) {
			return db.Get(append([]byte(config.BlockHashPrefix), key...))
		})
}

// GetValidatorHandler writes the validator corresponding to given address key
func GetValidatorHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildItemByIDHandler(db, "id", config.ValidatorPrefix)
}

// GetEntityHandler writes a single entity
func GetEntityHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildItemByHeightHandler(db, config.LatestEntityHeight, config.EntityIDPrefix, nil)
}

// GetProcessHandler writes a single process
func GetProcessHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildItemByHeightHandler(db, config.LatestProcessHeight, config.ProcessIDPrefix, nil)
}

// ListEntitiesHandler writes a list of entities from 'from'
func ListEntitiesHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildListItemsHandler(db, config.EntityIDPrefix, nil)
}

// ListProcessesHandler writes a list of processes from 'from'
func ListProcessesHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildListItemsHandler(db, config.ProcessIDPrefix, nil)
}

// ListProcessesByEntityHandler writes a list of processes belonging to 'entity'
func ListProcessesByEntityHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildListItemsByParent(db, "entity", config.EntityProcessHeightMapKey, config.ProcessByEntityPrefix, config.ProcessIDPrefix, true)
}

// ListValidatorsHandler writes a list of validators from 'from'
func ListValidatorsHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildListItemsHandler(db,
		config.ValidatorHeightPrefix,
		func(key []byte) ([]byte, error) {
			return db.Get(append([]byte(config.ValidatorPrefix), key...))
		})
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
	return buildListItemsByParent(db, "proposer", config.ValidatorHeightMapKey, config.BlockByValidatorPrefix, config.BlockHashPrefix, false)
}

// ListEnvelopesByProcessHandler writes a list of envelopes which share the given process
func ListEnvelopesByProcessHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildListItemsByParent(db, "process", config.ProcessEnvelopeHeightMapKey, config.EnvPIDPrefix, config.EnvPackagePrefix, true)
}

// ListEnvelopesHandler writes a list of envelopes
func ListEnvelopesHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildListItemsHandler(db, config.EnvPackagePrefix, nil)
}

// EnvelopeHeightByProcessHandler writes the number of envelopes which share the given processID
func EnvelopeHeightByProcessHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildHeightByParentHandler(db, "process", config.ProcessEnvelopeHeightMapKey)
}

// NumBlocksByValidatorHandler writes the number of blocks which share the given proposer
func NumBlocksByValidatorHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildHeightByParentHandler(db, "proposer", config.ValidatorHeightMapKey)
}

// ProcessHeightByEntityHandler writes the number of processes which share the given entity
func ProcessHeightByEntityHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildHeightByParentHandler(db, "entity", config.EntityProcessHeightMapKey)
}

// ListTxsHandler writes a list of txs starting with the given height key
func ListTxsHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		froms, ok := r.URL.Query()["from"]
		if !ok || len(froms[0]) < 1 {
			log.Errorf("Url Param 'from' is missing")
			http.Error(w, "Url Param 'from' missing", http.StatusNotFound)
			return
		}
		from, err := strconv.Atoi(froms[0])
		if err != nil {
			log.Error(err)
		}
		hashes := listItemsByHeight(db, config.ListSize, from, []byte(config.TxHeightPrefix))
		if len(hashes) == 0 {
			log.Errorf("No txs available at height %d", from)
			http.Error(w, "No txs available", http.StatusNotFound)
			return
		}
		var rawTxs types.ItemList
		for _, hash := range hashes {
			rawTx, err := db.Get(append([]byte(config.TxHashPrefix), hash...))
			if err != nil {
				log.Error(err)
			}
			var tx types.StoreTx
			err = proto.Unmarshal(rawTx, &tx)
			if err != nil {
				log.Error(err)
			}
			send := types.SendTx{
				Hash:  hash,
				Store: &tx,
			}
			rawSend, err := proto.Marshal(&send)
			if err != nil {
				log.Error(err)
			}
			rawTxs.Items = append(rawTxs.GetItems(), rawSend)
			if err != nil {
				log.Error(err)
			}
		}
		msg, err := proto.Marshal(&rawTxs)
		if err != nil {
			log.Error(err)
		}
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
			http.Error(w, "Url Param 'id' missing", http.StatusNotFound)
			return
		}
		height, err := strconv.Atoi(ids[0])
		if err != nil {
			log.Error(err)
		}
		id := util.EncodeInt(height)
		key := append([]byte(config.TxHeightPrefix), id...)
		hash, err := db.Get(key)
		if err != nil {
			log.Error(err)
			http.Error(w, "Unable to get tx hash", http.StatusNotFound)
			return
		}
		raw, err := db.Get(append([]byte(config.TxHashPrefix), hash...))
		if err != nil {
			log.Error(err)
			http.Error(w, "Unable to get raw tx", http.StatusNotFound)
			return
		}

		var tx types.StoreTx
		err = proto.Unmarshal(raw, &tx)
		if err != nil {
			log.Error(err)
		}

		send := types.SendTx{
			Hash:  hash,
			Store: &tx,
		}

		rawTx, err := proto.Marshal(&send)
		if err != nil {
			log.Error(err)
		}
		w.Write(rawTx)
		log.Debugf("Sent tx %d", height)
	}
}

// TxHeightFromHashHandler indirects the given tx hash
func TxHeightFromHashHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ids, ok := r.URL.Query()["hash"]
		if !ok || len(ids[0]) < 1 {
			log.Errorf("Url Param 'id' is missing")
			http.Error(w, "Url Param 'id' missing", http.StatusNotFound)
			return
		}
		id := ids[0]
		hash, err := hex.DecodeString(id)
		if err != nil {
			log.Error(err)
			log.Errorf("Cannot decode tx hash")
			http.Error(w, "Tx hash invalid", http.StatusNotFound)
			return
		}
		key := append([]byte(config.TxHashPrefix), hash...)
		has, err := db.Has(key)
		if err != nil {
			log.Error(err)
			log.Errorf("Tx Height not found")
			http.Error(w, "Tx hash invalid", http.StatusNotFound)
			return
		}
		if !has {
			log.Errorf("Tx hash key not found")
			http.Error(w, "Tx hash key not found", http.StatusNotFound)
			return
		}
		raw, err := db.Get(key)
		if err != nil {
			log.Error(err)
			http.Error(w, "Tx hash not found", http.StatusNotFound)
			return
		}

		var tx types.StoreTx
		err = proto.Unmarshal(raw, &tx)
		if err != nil {
			log.Error(err)
			http.Error(w, "Unable to get tx", http.StatusNotFound)
			return
		}
		height := &types.Height{Height: tx.GetHeight()}
		rawHeight, err := proto.Marshal(height)
		if err != nil {
			log.Error(err)
			http.Error(w, "Tx height invalid", http.StatusNotFound)
			return
		}
		w.Write(rawHeight)
	}
}

// EnvelopeHeightFromNullifierHandler returns the height of the corresponding envelope nullifier
func EnvelopeHeightFromNullifierHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ids, ok := r.URL.Query()["nullifier"]
		if !ok || len(ids[0]) < 1 {
			log.Errorf("Url Param 'nullifier' is missing")
			http.Error(w, "Url Param 'nullifier' missing", http.StatusNotFound)
			return
		}
		id := ids[0]
		hash, err := hex.DecodeString(id)
		if err != nil {
			log.Error(err)
			http.Error(w, "Unable to decode nullifier", http.StatusNotFound)
			return
		}
		key := append([]byte(config.EnvNullifierPrefix), hash...)
		has, err := db.Has(key)
		if err != nil {
			log.Error(err)
			http.Error(w, "Unable to get envelope height", http.StatusNotFound)
			return
		}
		if !has {
			log.Errorf("Envelope nullifier does not exist")
			http.Error(w, "Envelope nullifier does not exist", http.StatusNotFound)
			return
		}
		raw, err := db.Get(key)
		if err != nil {
			log.Error(err)
			http.Error(w, "Unable to get envelope key", http.StatusNotFound)
			return
		}

		var height types.Height
		err = proto.Unmarshal(raw, &height)
		if err != nil {
			log.Error(err)
			http.Error(w, "Unable to unmarshal height", http.StatusNotFound)
			return
		}
		w.Write(raw)

	}
}
