package router

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/vocdoni/vocexplorer/config"
	"github.com/vocdoni/vocexplorer/db"
	ptypes "github.com/vocdoni/vocexplorer/proto"
	"github.com/vocdoni/vocexplorer/util"
	"go.vocdoni.io/dvote/log"
	"google.golang.org/protobuf/proto"
)

// ListTxsHandler writes a list of transactions starting with the given height key
func ListTxsHandler(d *db.ExplorerDB) func(w http.ResponseWriter, r *http.Request) {
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
		hashes := db.ListItemsByHeight(d.Db, config.ListSize, from, []byte(config.TxHeightPrefix))
		if len(hashes) == 0 {
			log.Warnf("No transactions available at height %d", from)
			http.Error(w, "No transactions available", http.StatusInternalServerError)
			return
		}
		var rawTxs ptypes.ItemList
		for _, hash := range hashes {
			rawTx, err := d.Db.Get(append([]byte(config.TxHashPrefix), hash...))
			if err != nil {
				log.Warn(err)
				http.Error(w, "Unable to get raw tx", http.StatusInternalServerError)
				return
			}
			rawTxs.Items = append(rawTxs.GetItems(), packTransaction(rawTx))
			if err != nil {
				log.Warn(err)
			}
		}
		msg, err := json.Marshal(&rawTxs)
		if err != nil {
			log.Warn(err)
		}
		w.Write(msg)
		log.Debugf("Sent %d transactions", len(rawTxs.GetItems()))
	}
}

// GetTxByHeightHandler writes the tx corresponding to given height key
func GetTxByHeightHandler(d *db.ExplorerDB) func(w http.ResponseWriter, r *http.Request) {
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
		hash, err := d.Db.Get(key)
		if err != nil {
			log.Warn(err)
			http.Error(w, "Unable to get tx hash", http.StatusInternalServerError)
			return
		}
		raw, err := d.Db.Get(append([]byte(config.TxHashPrefix), hash...))
		if err != nil {
			log.Warn(err)
			http.Error(w, "Unable to get raw tx", http.StatusInternalServerError)
			return
		}
		w.Write(packTransaction(raw))
		log.Debugf("Sent tx %d", height)
	}
}

// GetTxByHashHandler writes the tx corresponding to given hash key
func GetTxByHashHandler(d *db.ExplorerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildItemByIDHandler(d, "id", config.TxHashPrefix, nil, packTransaction)
}

// TxHeightFromHashHandler indirects the given tx hash
func TxHeightFromHashHandler(d *db.ExplorerDB) func(w http.ResponseWriter, r *http.Request) {
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
		has, err := d.Db.Has(key)
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
		raw, err := d.Db.Get(key)
		if err != nil {
			log.Warn(err)
			http.Error(w, "Tx hash not found", http.StatusInternalServerError)
			return
		}

		var tx ptypes.Transaction
		err = proto.Unmarshal(raw, &tx)
		if err != nil {
			log.Warn(err)
			http.Error(w, "Unable to get tx", http.StatusInternalServerError)
			return
		}
		height := &ptypes.Height{Height: tx.GetTxHeight()}
		rawHeight, err := json.Marshal(height.Mirror())
		if err != nil {
			log.Warn(err)
			http.Error(w, "Tx height invalid", http.StatusInternalServerError)
			return
		}
		w.Write(rawHeight)
	}
}

// SearchTransactionsHandler writes a list of transactions by search term
func SearchTransactionsHandler(d *db.ExplorerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildSearchHandler(d,
		config.TxHashPrefix,
		true,
		func(hash []byte) ([]byte, error) {
			log.Debugf("Hash found: %X", hash)
			rawTx, err := d.Db.Get(append([]byte(config.TxHashPrefix), hash...))
			if err != nil {
				return nil, err
			}
			return rawTx, nil
		},
		packTransaction,
	)
}

func packTransaction(raw []byte) []byte {
	var item ptypes.Transaction
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
