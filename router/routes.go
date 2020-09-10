package router

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	dvotedb "gitlab.com/vocdoni/go-dvote/db"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/api"
	"gitlab.com/vocdoni/vocexplorer/config"
	vocdb "gitlab.com/vocdoni/vocexplorer/db"
	ptypes "gitlab.com/vocdoni/vocexplorer/proto"
	"gitlab.com/vocdoni/vocexplorer/util"
	"google.golang.org/protobuf/proto"
)

//PingHandler responds to a ping
func PingHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	}
}

// StatsHandler is the public api for all blockchain statistics & information
func StatsHandler(db *dvotedb.BadgerDB, cfg *config.Cfg) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debugf("Serving statistics ")
		stats := new(api.VochainStats)

		if !cfg.Detached {
			// If unable to get api information, don't return error so db information can still serve
			t := api.StartTendermintClient(cfg.TendermintHost)
			status := api.GetHealth(t)
			if status == nil {
				log.Warnf("Unable to get vochain status")
			} else {
				stats.Network = status.NodeInfo.Network
				stats.Version = status.NodeInfo.Version
				stats.SyncInfo = status.SyncInfo
			}

			genesis := api.GetGenesis(t)
			if status == nil {
				log.Warnf("Unable to get genesis block")
			} else {
				stats.GenesisTimeStamp = genesis.GenesisTime
				stats.ChainID = genesis.ChainID
			}

			gw, cancel := api.InitGateway(cfg.GatewayHost)
			defer cancel()
			blockTime, blockTimeStamp, height, err := gw.GetBlockStatus()
			if err != nil {
				log.Warn(err)
			} else {
				stats.BlockTime = blockTime
				stats.BlockTimeStamp = blockTimeStamp
				stats.Height = height
			}
		}

		blockHeight := vocdb.GetHeight(db, config.LatestBlockHeightKey, 1)
		entityCount := vocdb.GetHeight(db, config.LatestEntityCountKey, 0)
		envelopeCount := vocdb.GetHeight(db, config.LatestEnvelopeCountKey, 0)
		processCount := vocdb.GetHeight(db, config.LatestProcessCountKey, 0)
		transactionHeight := vocdb.GetHeight(db, config.LatestTxHeightKey, 1)
		validatorCount := vocdb.GetHeight(db, config.LatestValidatorCountKey, 0)

		stats.BlockHeight = blockHeight.GetHeight()
		stats.EntityCount = entityCount.GetHeight()
		stats.EnvelopeCount = envelopeCount.GetHeight()
		stats.ProcessCount = processCount.GetHeight()
		stats.TransactionHeight = transactionHeight.GetHeight()
		stats.ValidatorCount = validatorCount.GetHeight()

		msg, err := json.Marshal(stats)
		if err != nil {
			log.Warn(err)
			http.Error(w, "Unable to marshal stats", http.StatusInternalServerError)
			return
		}

		w.Write(msg)
	}
}

// HeightHandler writes the int64 height value corresponding to given key
func HeightHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		keys, ok := r.URL.Query()["key"]
		if !ok || len(keys[0]) < 1 {
			log.Warnf("Url Param 'key' is missing")
			http.Error(w, "Url Param 'key' missing", http.StatusBadRequest)
			return
		}
		val, err := db.Get([]byte(keys[0]))
		if err != nil {
			log.Warn(err)
			http.Error(w, "Key not found", http.StatusInternalServerError)
			return
		}

		var height ptypes.Height
		err = proto.Unmarshal(val, &height)
		if err != nil {
			log.Warn(err)
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
			log.Warnf("Url Param 'key' is missing")
			http.Error(w, "Url Param 'key' missing", http.StatusBadRequest)
			return
		}
		val, err := db.Get([]byte(keys[0]))
		if err != nil {
			log.Warnf("Key %s not found: %s", keys[0], err.Error())
			http.Error(w, "Key not found", http.StatusInternalServerError)
			return
		}

		var heightMap ptypes.HeightMap
		err = proto.Unmarshal(val, &heightMap)
		if err != nil {
			log.Warn(err)
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
			log.Warnf("Url Param '" + IDName + "' is missing")
			http.Error(w, "Url Param '"+IDName+"' missing", http.StatusBadRequest)
			return
		}
		id := ids[0]
		addressBytes, err := hex.DecodeString(id)
		if err != nil {
			log.Warn(err)
		}
		key := append([]byte(itemPrefix), addressBytes...)
		raw, err := db.Get(key)
		if err != nil {
			log.Warn(err)
		}
		w.Write(raw)
		log.Debugf("Sent Item")
	}
}

func buildItemByHeightHandler(db *dvotedb.BadgerDB, heightKey, key string, getItem func(key []byte) ([]byte, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		heights, ok := r.URL.Query()["height"]
		if !ok || len(heights[0]) < 1 {
			log.Warnf("Url Param 'height' is missing")
			http.Error(w, "Url Param 'height' missing", http.StatusBadRequest)
			return
		}
		height, err := strconv.Atoi(heights[0])
		if err != nil {
			log.Warn(err)
			http.Error(w, "Cannot decode height", http.StatusInternalServerError)
			return
		}

		envHeight := vocdb.GetHeight(db, heightKey, 0)
		if height > int(envHeight.GetHeight()) {
			log.Warnf("Requested item does not exist")
			http.Error(w, "Requested item does not exist", http.StatusInternalServerError)
			return
		}
		itemKey := append([]byte(key), util.EncodeInt(height)...)
		rawItem, err := db.Get(itemKey)
		if err != nil {
			log.Warn(err)
			http.Error(w, "item not found", http.StatusInternalServerError)
			return
		}
		if getItem != nil {
			rawItem, err = getItem(rawItem)
			if err != nil {
				log.Warn(err)
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
			log.Warnf("Url Param 'from' is missing")
			http.Error(w, "Url Param 'from' missing", http.StatusBadRequest)
			return
		}
		from, err := strconv.Atoi(froms[0])
		if err != nil {
			log.Warn(err)
		}
		items := vocdb.ListItemsByHeight(db, config.ListSize, from, []byte(key))
		if len(items) == 0 {
			log.Warnf("Retrieved no items from key %s, index %d", key, from)
			http.Error(w, "No items available", http.StatusInternalServerError)
			return
		}
		var itemList ptypes.ItemList
		for _, rawItem := range items {
			if getItem != nil {
				rawItem, err = getItem(rawItem)
				if err != nil {
					log.Warn(err)
					http.Error(w, "item not found", http.StatusInternalServerError)
					return
				}
			}
			itemList.Items = append(itemList.GetItems(), rawItem)
		}

		msg, err := proto.Marshal(&itemList)
		if err != nil {
			log.Warn(err)
			http.Error(w, "Unable to encode data", http.StatusInternalServerError)
			return
		}
		w.Write(msg)
		log.Debugf("Sent %d items", len(itemList.GetItems()))
	}
}

func buildSearchHandler(db *dvotedb.BadgerDB, key string, getKey bool, getItem func(key []byte) ([]byte, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		terms, ok := r.URL.Query()["term"]
		if !ok || len(terms[0]) < 1 {
			log.Warnf("Url Param 'term' is missing")
			http.Error(w, "Url Param 'term' missing", http.StatusBadRequest)
			return
		}
		searchTerm := strings.ToLower(terms[0])

		var err error
		var items [][]byte
		if getKey {
			items = vocdb.SearchKeys(db, config.ListSize, searchTerm, []byte(key))
			// items = vocdb.SearchKeys(db, config.ListSize, term, []byte(key))
		} else {
			items = vocdb.SearchItems(db, config.ListSize, searchTerm, []byte(key))
			// items = vocdb.SearchItems(db, config.ListSize, term, []byte(key))
		}
		if len(items) == 0 {
			log.Warn("Retrieved no items")
			http.Error(w, "No items available", http.StatusInternalServerError)
			return
		}
		var itemList ptypes.ItemList
		for _, rawItem := range items {
			if getItem != nil {
				rawItem, err = getItem(rawItem)
				if err != nil {
					log.Warnf("No item found for key %s, search term %s: %s", key, searchTerm, err.Error())
					http.Error(w, "item not found", http.StatusInternalServerError)
					return
				}
			}
			itemList.Items = append(itemList.GetItems(), rawItem)
		}

		msg, err := proto.Marshal(&itemList)
		if err != nil {
			log.Warn(err)
			http.Error(w, "Unable to encode data", http.StatusInternalServerError)
			return
		}
		w.Write(msg)
		log.Debugf("Sent %d items for search term %s, key %s", len(itemList.GetItems()), searchTerm, key)
	}
}

func buildListItemsByParent(db *dvotedb.BadgerDB, parentName, heightMapKey, getHeightPrefix, itemPrefix string, marshalHeight bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		froms, ok := r.URL.Query()["from"]
		if !ok || len(froms[0]) < 1 {
			log.Warnf("Url Param 'from' is missing")
			http.Error(w, "Url Param 'from' missing", http.StatusBadRequest)
			return
		}
		parents, ok := r.URL.Query()[parentName]
		if !ok || len(parents[0]) < 1 {
			log.Warnf("Url Param '" + parentName + "' is missing")
			http.Error(w, "Url Param '"+parentName+"' is missing", http.StatusBadRequest)
			return
		}
		from, err := strconv.Atoi(froms[0])
		if err != nil {
			log.Warn(err)
		}

		heightMap := vocdb.GetHeightMap(db, heightMapKey)
		itemHeight, ok := heightMap.Heights[parents[0]]
		if !ok {
			log.Warn("Parent does not exist")
			http.Error(w, "No items available", http.StatusInternalServerError)
			return
		}
		from = util.Min(from, int(itemHeight))

		// Get keys
		parentBytes, err := hex.DecodeString(parents[0])
		if err != nil {
			log.Warn(err)
		}
		keys := vocdb.ListItemsByHeight(db, config.ListSize, from, append([]byte(getHeightPrefix), parentBytes...))
		if len(keys) == 0 {
			log.Warn("No keys retrieved")
			http.Error(w, "No items available", http.StatusInternalServerError)
			return
		}
		var rawItems ptypes.ItemList
		// Get packages by height:package
		for _, rawKey := range keys {
			var rawPackage []byte
			if marshalHeight {
				height := new(ptypes.Height)
				if err := proto.Unmarshal(rawKey, height); err != nil {
					log.Warn(err)
				}
				log.Debugf("Getting item from height %d", height.GetHeight())
				rawPackage, err = db.Get(append([]byte(itemPrefix), util.EncodeInt(height.GetHeight())...))
			} else {
				rawPackage, err = db.Get(append([]byte(itemPrefix), rawKey...))
			}
			if err != nil {
				log.Warn(err)
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
			log.Warn(err)
			http.Error(w, "Unable to encode data", http.StatusInternalServerError)
			return
		}
		if err != nil {
			log.Warn(err)
		}
		w.Write(msg)
		log.Debugf("Sent %d items by %s %s", len(rawItems.GetItems()), parentName, parents[0])
	}
}

func buildHeightByParentHandler(db *dvotedb.BadgerDB, parentName, heightMapKey string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		parents, ok := r.URL.Query()[parentName]
		if !ok || len(parents[0]) < 1 {
			log.Warnf("Url Param '" + parentName + "' is missing")
			http.Error(w, "Url Param 'process' missing", http.StatusBadRequest)
			return
		}
		var heightMap ptypes.HeightMap
		has, err := db.Has([]byte(heightMapKey))
		if err != nil || !has {
			log.Warn("No item height not found")
			http.Error(w, "No items available", http.StatusInternalServerError)
			return
		}
		if has {
			rawHeightMap, err := db.Get([]byte(heightMapKey))
			if err != nil {
				log.Warn(err)
			}
			err = proto.Unmarshal(rawHeightMap, &heightMap)
			if err != nil {
				log.Warn(err)
			}
		}
		height, ok := heightMap.Heights[parents[0]]
		if !ok {
			height = 0
		}
		sendHeight := &ptypes.Height{Height: int64(height)}
		msg, err := proto.Marshal(sendHeight)
		if err != nil {
			log.Warn(err)
			http.Error(w, "Unable to encode data", http.StatusInternalServerError)
			return
		}
		w.Write(msg)
		log.Debugf("Found %d items by %s %s", sendHeight.GetHeight(), parentName, parents[0])
	}
}

// GetEnvelopeHandler writes a single envelope
func GetEnvelopeHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildItemByHeightHandler(db, config.LatestEnvelopeCountKey, config.EnvPackagePrefix, nil)
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
	return buildItemByHeightHandler(db, config.LatestEntityCountKey, config.EntityHeightPrefix, nil)
}

// GetProcessHandler writes a single process
func GetProcessHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildItemByHeightHandler(db, config.LatestProcessCountKey, config.ProcessHeightPrefix, nil)
}

// ListEntitiesHandler writes a list of entities from 'from'
func ListEntitiesHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildListItemsHandler(db, config.EntityHeightPrefix, nil)
}

// ListProcessesHandler writes a list of processes from 'from'
func ListProcessesHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildListItemsHandler(db, config.ProcessHeightPrefix, nil)
}

// ListProcessesByEntityHandler writes a list of processes belonging to 'entity'
func ListProcessesByEntityHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildListItemsByParent(db, "entity", config.EntityProcessCountMapKey, config.ProcessByEntityPrefix, config.ProcessHeightPrefix, true)
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
	return buildListItemsByParent(db, "process", config.ProcessEnvelopeCountMapKey, config.EnvPIDPrefix, config.EnvPackagePrefix, true)
}

// ListEnvelopesHandler writes a list of envelopes
func ListEnvelopesHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildListItemsHandler(db, config.EnvPackagePrefix, nil)
}

// EnvelopeHeightByProcessHandler writes the number of envelopes which share the given processID
func EnvelopeHeightByProcessHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildHeightByParentHandler(db, "process", config.ProcessEnvelopeCountMapKey)
}

// NumBlocksByValidatorHandler writes the number of blocks which share the given proposer
func NumBlocksByValidatorHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildHeightByParentHandler(db, "proposer", config.ValidatorHeightMapKey)
}

// ProcessHeightByEntityHandler writes the number of processes which share the given entity
func ProcessHeightByEntityHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildHeightByParentHandler(db, "entity", config.EntityProcessCountMapKey)
}

// ListTxsHandler writes a list of transactions starting with the given height key
func ListTxsHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		froms, ok := r.URL.Query()["from"]
		if !ok || len(froms[0]) < 1 {
			log.Warnf("Url Param 'from' is missing")
			http.Error(w, "Url Param 'from' missing", http.StatusNotFound)
			return
		}
		from, err := strconv.Atoi(froms[0])
		if err != nil {
			log.Warn(err)
		}
		hashes := vocdb.ListItemsByHeight(db, config.ListSize, from, []byte(config.TxHeightPrefix))
		if len(hashes) == 0 {
			log.Warnf("No transactions available at height %d", from)
			http.Error(w, "No transactions available", http.StatusInternalServerError)
			return
		}
		var rawTxs ptypes.ItemList
		for _, hash := range hashes {
			rawTx, err := db.Get(append([]byte(config.TxHashPrefix), hash...))
			if err != nil {
				log.Warn(err)
			}
			var tx ptypes.StoreTx
			err = proto.Unmarshal(rawTx, &tx)
			if err != nil {
				log.Warn(err)
			}
			send := ptypes.SendTx{
				Hash:  hash,
				Store: &tx,
			}
			rawSend, err := proto.Marshal(&send)
			if err != nil {
				log.Warn(err)
			}
			rawTxs.Items = append(rawTxs.GetItems(), rawSend)
			if err != nil {
				log.Warn(err)
			}
		}
		msg, err := proto.Marshal(&rawTxs)
		if err != nil {
			log.Warn(err)
		}
		w.Write(msg)
		log.Debugf("Sent %d transactions", len(rawTxs.GetItems()))
	}
}

// GetTxHandler writes the tx corresponding to given height key
func GetTxHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ids, ok := r.URL.Query()["id"]
		if !ok || len(ids[0]) < 1 {
			log.Warnf("Url Param 'id' is missing")
			http.Error(w, "Url Param 'id' missing", http.StatusNotFound)
			return
		}
		height, err := strconv.Atoi(ids[0])
		if err != nil {
			log.Warn(err)
		}
		id := util.EncodeInt(height)
		key := append([]byte(config.TxHeightPrefix), id...)
		hash, err := db.Get(key)
		if err != nil {
			log.Warn(err)
			http.Error(w, "Unable to get tx hash", http.StatusInternalServerError)
			return
		}
		raw, err := db.Get(append([]byte(config.TxHashPrefix), hash...))
		if err != nil {
			log.Warn(err)
			http.Error(w, "Unable to get raw tx", http.StatusInternalServerError)
			return
		}

		var tx ptypes.StoreTx
		err = proto.Unmarshal(raw, &tx)
		if err != nil {
			log.Warn(err)
		}

		send := ptypes.SendTx{
			Hash:  hash,
			Store: &tx,
		}

		rawTx, err := proto.Marshal(&send)
		if err != nil {
			log.Warn(err)
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
			log.Warnf("Url Param 'id' is missing")
			http.Error(w, "Url Param 'id' missing", http.StatusNotFound)
			return
		}
		id := ids[0]
		hash, err := hex.DecodeString(id)
		if err != nil {
			log.Warn(err)
			log.Warnf("Cannot decode tx hash")
			http.Error(w, "Tx hash invalid", http.StatusInternalServerError)
			return
		}
		key := append([]byte(config.TxHashPrefix), hash...)
		has, err := db.Has(key)
		if err != nil {
			log.Warn(err)
			log.Warnf("Tx Height not found")
			http.Error(w, "Tx hash invalid", http.StatusInternalServerError)
			return
		}
		if !has {
			log.Warnf("Tx hash key not found")
			http.Error(w, "Tx hash key not found", http.StatusInternalServerError)
			return
		}
		raw, err := db.Get(key)
		if err != nil {
			log.Warn(err)
			http.Error(w, "Tx hash not found", http.StatusInternalServerError)
			return
		}

		var tx ptypes.StoreTx
		err = proto.Unmarshal(raw, &tx)
		if err != nil {
			log.Warn(err)
			http.Error(w, "Unable to get tx", http.StatusInternalServerError)
			return
		}
		height := &ptypes.Height{Height: tx.GetHeight()}
		rawHeight, err := proto.Marshal(height)
		if err != nil {
			log.Warn(err)
			http.Error(w, "Tx height invalid", http.StatusInternalServerError)
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
			log.Warnf("Url Param 'nullifier' is missing")
			http.Error(w, "Url Param 'nullifier' missing", http.StatusNotFound)
			return
		}
		id := ids[0]
		hash, err := hex.DecodeString(id)
		if err != nil {
			log.Warn(err)
			http.Error(w, "Unable to decode nullifier", http.StatusInternalServerError)
			return
		}
		key := append([]byte(config.EnvNullifierPrefix), hash...)
		has, err := db.Has(key)
		if err != nil {
			log.Warn(err)
			http.Error(w, "Unable to get envelope height", http.StatusInternalServerError)
			return
		}
		if !has {
			log.Warnf("Envelope nullifier does not exist")
			http.Error(w, "Envelope nullifier does not exist", http.StatusInternalServerError)
			return
		}
		raw, err := db.Get(key)
		if err != nil {
			log.Warn(err)
			http.Error(w, "Unable to get envelope key", http.StatusInternalServerError)
			return
		}

		var height ptypes.Height
		err = proto.Unmarshal(raw, &height)
		if err != nil {
			log.Warn(err)
			http.Error(w, "Unable to unmarshal height", http.StatusInternalServerError)
			return
		}
		w.Write(raw)

	}
}

// SearchBlocksHandler writes a list of blocks by search term
func SearchBlocksHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildSearchHandler(db,
		config.BlockHashPrefix,
		false,
		nil,
	)
}

// SearchBlocksByValidatorHandler writes a list of blocks by search term belonging to given validator
func SearchBlocksByValidatorHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		terms, ok := r.URL.Query()["term"]
		if !ok || len(terms[0]) < 1 {
			log.Warnf("Url Param 'term' is missing")
			http.Error(w, "Url Param 'term' missing", http.StatusBadRequest)
			return
		}
		searchTerm := strings.ToLower(terms[0])

		validators, ok := r.URL.Query()["validator"]
		if !ok || len(validators[0]) < 1 {
			log.Warnf("Url Param 'validator' is missing")
			http.Error(w, "Url Param 'validator' missing", http.StatusBadRequest)
			return
		}
		validator := strings.ToLower(validators[0])

		var err error
		var items [][]byte
		items = vocdb.SearchBlocksByValidator(db, config.ListSize, searchTerm, validator)
		if len(items) == 0 {
			log.Warn("Retrieved no items")
			http.Error(w, "No items available", http.StatusInternalServerError)
			return
		}
		var itemList ptypes.ItemList
		for _, rawItem := range items {
			itemList.Items = append(itemList.GetItems(), rawItem)
		}

		msg, err := proto.Marshal(&itemList)
		if err != nil {
			log.Warn(err)
			http.Error(w, "Unable to encode data", http.StatusInternalServerError)
			return
		}
		w.Write(msg)
		log.Debugf("Sent %d blocks for search term %s, validator %s", len(itemList.GetItems()), searchTerm, validator)
	}
}

// SearchTransactionsHandler writes a list of transactions by search term
func SearchTransactionsHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildSearchHandler(db,
		config.TxHashPrefix,
		true,
		func(hash []byte) ([]byte, error) {
			log.Debugf("Hash found: %X", hash)
			rawTx, err := db.Get(append([]byte(config.TxHashPrefix), hash...))
			if err != nil {
				return []byte{}, err
			}
			var tx ptypes.StoreTx
			err = proto.Unmarshal(rawTx, &tx)
			if err != nil {
				return []byte{}, err
			}
			send := ptypes.SendTx{
				Hash:  hash,
				Store: &tx,
			}
			rawSend, err := proto.Marshal(&send)
			if err != nil {
				return []byte{}, err
			}
			return rawSend, nil
		},
	)
}

// SearchEnvelopesHandler writes a list of envelopes by search term
func SearchEnvelopesHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildSearchHandler(db,
		config.EnvNullifierPrefix,
		false,
		func(key []byte) ([]byte, error) {
			height := &ptypes.Height{}
			err := proto.Unmarshal(key, height)
			if err != nil {
				log.Warn("Unable to unmarshal envelope height")
			}
			return db.Get(append([]byte(config.EnvPackagePrefix), util.EncodeInt(height.GetHeight())...))
		},
	)
}

// SearchProcessesHandler writes a list of processes by search term
func SearchProcessesHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildSearchHandler(db,
		config.ProcessIDPrefix,
		true,
		nil,
	)
}

// SearchEntitiesHandler writes a list of entities by search term
func SearchEntitiesHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildSearchHandler(db,
		config.EntityIDPrefix,
		true,
		nil,
	)
}

// SearchValidatorsHandler writes a list of validators by search term
func SearchValidatorsHandler(db *dvotedb.BadgerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildSearchHandler(db,
		config.ValidatorPrefix,
		false,
		nil,
	)
}
