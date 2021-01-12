package router

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/db"
	ptypes "gitlab.com/vocdoni/vocexplorer/proto"
	"gitlab.com/vocdoni/vocexplorer/util"
	"go.vocdoni.io/dvote/log"
	"google.golang.org/protobuf/proto"
)

// GetBlockHeaderHandler writes a StoreBlock from the DB
func GetBlockHeaderHandler(d *db.ExplorerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildItemByHeightHandler(d,
		config.LatestBlockHeightKey,
		config.BlockHeightPrefix,
		func(key []byte) ([]byte, error) {
			return d.Db.Get(append([]byte(config.BlockHashPrefix), key...))
		},
		packBlock,
	)
}

// GetBlockHandler writes a full block from the vochain
func GetBlockHandler(d *db.ExplorerDB) func(w http.ResponseWriter, r *http.Request) {
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
		block, err := d.Vs.GetBlock(int64(height))
		if err != nil {
			log.Warn(err)
			http.Error(w, "Cannot get block at height "+util.IntToString(height), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(block.Block)
	}
}

// ListBlocksHandler writes a list of blocks by height
func ListBlocksHandler(d *db.ExplorerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildListItemsHandler(d,
		config.BlockHeightPrefix,
		func(key []byte) ([]byte, error) {
			return d.Db.Get(append([]byte(config.BlockHashPrefix), key...))
		},
		packBlock,
	)
}

// ListBlocksByValidatorHandler writes a list of blocks which share the given proposer
func ListBlocksByValidatorHandler(d *db.ExplorerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildListItemsByParent(d, "proposer", config.ValidatorHeightMapKey, config.BlockByValidatorPrefix, config.BlockHashPrefix, false, packBlock)
}

// NumBlocksByValidatorHandler writes the number of blocks which share the given proposer
func NumBlocksByValidatorHandler(d *db.ExplorerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildHeightByParentHandler(d, "proposer", config.ValidatorHeightMapKey)
}

// SearchBlocksHandler writes a list of blocks by search term
func SearchBlocksHandler(d *db.ExplorerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildSearchHandler(d,
		config.BlockHashPrefix,
		false,
		nil,
		packBlock,
	)
}

// SearchBlocksByValidatorHandler writes a list of blocks by search term belonging to given validator
func SearchBlocksByValidatorHandler(d *db.ExplorerDB) func(w http.ResponseWriter, r *http.Request) {
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
		items := db.SearchBlocksByValidator(d.Db, config.ListSize, searchTerm, validator)
		if len(items) == 0 {
			log.Warn("Retrieved no items")
			http.Error(w, "No items available", http.StatusInternalServerError)
			return
		}
		var itemList ptypes.ItemList
		for _, rawItem := range items {
			itemList.Items = append(itemList.GetItems(), rawItem)
		}

		msg, err := json.Marshal(&itemList)
		if err != nil {
			log.Warn(err)
			http.Error(w, "Unable to encode data", http.StatusInternalServerError)
			return
		}
		w.Write(msg)
		log.Debugf("Sent %d blocks for search term %s, validator %s", len(itemList.GetItems()), searchTerm, validator)
	}
}

func packBlock(raw []byte) []byte {
	var item ptypes.StoreBlock
	err := proto.Unmarshal(raw, &item)
	if err != nil {
		log.Error(err)
	}
	new, err := json.Marshal(item.Mirror())
	if err != nil {
		log.Error(err)
	}
	return new
}
