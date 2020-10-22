package router

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/api/dbtypes"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/db"
	ptypes "gitlab.com/vocdoni/vocexplorer/proto"
	"gitlab.com/vocdoni/vocexplorer/util"
	"google.golang.org/protobuf/proto"
)

func buildItemByIDHandler(d *db.ExplorerDB, IDName, itemPrefix string, getItem func(key []byte) ([]byte, error), pack func([]byte) []byte) func(w http.ResponseWriter, r *http.Request) {
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
		raw, err := d.Db.Get(key)
		if err != nil {
			log.Warn(err)
		}
		if getItem != nil {
			raw, err = getItem(raw)
			if err != nil {
				log.Warn(err)
				http.Error(w, "item not found", http.StatusInternalServerError)
				return
			}
		}
		if pack != nil {
			raw = pack(raw)
		}
		w.Write(raw)
		log.Debugf("Sent Item")
	}
}

func buildItemByHeightHandler(d *db.ExplorerDB, heightKey, key string, getItem func(key []byte) ([]byte, error), pack func(key []byte) []byte) func(w http.ResponseWriter, r *http.Request) {
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

		envHeight := db.GetHeight(d.Db, heightKey, 0)
		if height > int(envHeight.GetHeight()) {
			log.Warnf("Requested item does not exist")
			http.Error(w, "Requested item does not exist", http.StatusInternalServerError)
			return
		}
		itemKey := append([]byte(key), util.EncodeInt(height)...)
		rawItem, err := d.Db.Get(itemKey)
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
		if pack != nil {
			rawItem = pack(rawItem)
		}
		w.Write(rawItem)
		log.Debugf("sent item with key %s, height %d", key, height)
	}
}

func buildListItemsHandler(d *db.ExplorerDB, key string, getItem func(key []byte) ([]byte, error), pack func(key []byte) []byte) func(w http.ResponseWriter, r *http.Request) {
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
		items := db.ListItemsByHeight(d.Db, config.ListSize, from, []byte(key))
		if len(items) == 0 {
			log.Warnf("Retrieved no items from key %s, index %d", key, from)
			http.Error(w, "No items available", http.StatusInternalServerError)
			return
		}
		var itemList dbtypes.ItemList
		for _, rawItem := range items {
			if getItem != nil {
				rawItem, err = getItem(rawItem)
				if err != nil {
					log.Warn(err)
					http.Error(w, "item not found", http.StatusInternalServerError)
					return
				}
			}
			if pack != nil {
				rawItem = pack(rawItem)
			}
			itemList.Items = append(itemList.Items, rawItem)
		}

		msg, err := json.Marshal(&itemList)
		if err != nil {
			log.Warn(err)
			http.Error(w, "Unable to encode data", http.StatusInternalServerError)
			return
		}
		w.Write(msg)
		log.Debugf("Sent %d items", len(itemList.Items))
	}
}

func buildSearchHandler(d *db.ExplorerDB, key string, getKey bool, getItem func(key []byte) ([]byte, error), pack func(key []byte) []byte) func(w http.ResponseWriter, r *http.Request) {
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
			items = db.SearchKeys(d.Db, config.ListSize, searchTerm, []byte(key))
			// items = db.SearchKeys(d.Db, config.ListSize, term, []byte(key))
		} else {
			items = db.SearchItems(d.Db, config.ListSize, searchTerm, []byte(key))
			// items = db.SearchItems(d.Db, config.ListSize, term, []byte(key))
		}
		if len(items) == 0 {
			log.Warn("Retrieved no items")
			http.Error(w, "No items available", http.StatusInternalServerError)
			return
		}
		var itemList dbtypes.ItemList
		for _, rawItem := range items {
			if getItem != nil {
				rawItem, err = getItem(rawItem)
				if err != nil {
					log.Warnf("No item found for key %s, search term %s: %s", key, searchTerm, err.Error())
					http.Error(w, "item not found", http.StatusInternalServerError)
					return
				}
			}
			if pack != nil {
				rawItem = pack(rawItem)
			}
			itemList.Items = append(itemList.Items, rawItem)
		}

		msg, err := json.Marshal(&itemList)
		if err != nil {
			log.Warn(err)
			http.Error(w, "Unable to encode data", http.StatusInternalServerError)
			return
		}
		w.Write(msg)
		log.Debugf("Sent %d items for search term %s, key %s", len(itemList.Items), searchTerm, key)
	}
}

func buildListItemsByParent(d *db.ExplorerDB, parentName, heightMapKey, getHeightPrefix, itemPrefix string, marshalHeight bool, pack func(key []byte) []byte) func(w http.ResponseWriter, r *http.Request) {
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

		heightMap := db.GetHeightMap(d.Db, heightMapKey)
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
		keys := db.ListItemsByHeight(d.Db, config.ListSize, from, append([]byte(getHeightPrefix), parentBytes...))
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
				rawPackage, err = d.Db.Get(append([]byte(itemPrefix), util.EncodeInt(height.GetHeight())...))
			} else {
				rawPackage, err = d.Db.Get(append([]byte(itemPrefix), rawKey...))
			}
			if err != nil {
				log.Warn(err)
			} else {
				if pack != nil {
					rawPackage = pack(rawPackage)
				}
				rawItems.Items = append(rawItems.GetItems(), rawPackage)
			}
		}
		if len(rawItems.Items) <= 0 {
			http.Error(w, "No items found", http.StatusInternalServerError)
			return
		}
		msg, err := json.Marshal(&rawItems)
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

func buildHeightByParentHandler(d *db.ExplorerDB, parentName, heightMapKey string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		parents, ok := r.URL.Query()[parentName]
		if !ok || len(parents[0]) < 1 {
			log.Warnf("Url Param '" + parentName + "' is missing")
			http.Error(w, "Url Param 'process' missing", http.StatusBadRequest)
			return
		}
		var heightMap ptypes.HeightMap
		has, err := d.Db.Has([]byte(heightMapKey))
		if err != nil || !has {
			log.Warn("No item height not found")
			http.Error(w, "No items available", http.StatusInternalServerError)
			return
		}
		if has {
			rawHeightMap, err := d.Db.Get([]byte(heightMapKey))
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
		msg, err := json.Marshal(sendHeight)
		if err != nil {
			log.Warn(err)
			http.Error(w, "Unable to encode data", http.StatusInternalServerError)
			return
		}
		w.Write(msg)
		log.Debugf("Found %d items by %s %s", sendHeight.GetHeight(), parentName, parents[0])
	}
}
